// Package service contém a lógica do grafo de capitais:
// montagem do grafo (Seed), exibição das adjacências (Show) e a busca do
// Caminho Mais Barato usando o algoritmo de Dijkstra com um heap (fila de
// prioridade).
package service

import (
	"container/heap"
	"errors"
	"math"
	"sort"

	"grafotb1/model"
)

// ErrCityNotFound é retornado quando a origem ou o destino não existem no grafo.
var ErrCityNotFound = errors.New("capital não encontrada")

// Service guarda o grafo de capitais como uma Lista de Adjacências.
//
//   - adj[A][B] = distância em km entre as capitais A e B.
//   - tolls[A]  = pedágio cobrado ao passar pela capital A.
//
// O grafo é NÃO direcionado: se existe a aresta A→B, também existe B→A com a
// mesma distância.
type Service struct {
	adj   map[string]map[string]int
	tolls map[string]int
}

// Seed monta o grafo a partir dos dados lidos do capitais.json.
//
// O JSON é uma lista de objetos e o nome da capital é a CHAVE de cada objeto
// (por isso o tipo []map[string]model.City).
func Seed(entries []map[string]model.City) *Service {
	s := &Service{
		adj:   make(map[string]map[string]int),
		tolls: make(map[string]int),
	}

	for _, entry := range entries {
		for name, city := range entry {
			s.tolls[name] = city.Toll
			// Garante que a capital exista no grafo mesmo que não tenha
			// vizinhos (ex.: Macapá não tem rota terrestre).
			if s.adj[name] == nil {
				s.adj[name] = make(map[string]int)
			}
			for neighbor, distance := range city.Neighbors {
				s.addEdge(name, neighbor, distance)
			}
		}
	}
	return s
}

// addEdge cria a aresta nos dois sentidos (grafo não direcionado).
func (s *Service) addEdge(from, to string, distance int) {
	if s.adj[from] == nil {
		s.adj[from] = make(map[string]int)
	}
	if s.adj[to] == nil {
		s.adj[to] = make(map[string]int)
	}
	s.adj[from][to] = distance
	s.adj[to][from] = distance
}

// Capitals devolve a lista de capitais em ordem alfabética.
// Serve para alimentar os selects/datalist do frontend.
func (s *Service) Capitals() []string {
	out := make([]string, 0, len(s.adj))
	for name := range s.adj {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// HasCity informa se a capital existe no grafo.
func (s *Service) HasCity(city string) bool {
	_, ok := s.adj[city]
	return ok
}

// Neighbor é uma capital vizinha com a distância correspondente (usado no Show).
type Neighbor struct {
	Name     string `json:"name"`
	Distance int    `json:"distance"`
}

// Show devolve, para cada capital, a lista ordenada de vizinhos com a distância.
// É a "Lista de Adjacências" pedida no trabalho.
func (s *Service) Show() map[string][]Neighbor {
	out := make(map[string][]Neighbor, len(s.adj))
	for city, neighbors := range s.adj {
		list := make([]Neighbor, 0, len(neighbors))
		for name, distance := range neighbors {
			list = append(list, Neighbor{Name: name, Distance: distance})
		}
		sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
		out[city] = list
	}
	return out
}

// --- Caminho Mais Barato (Dijkstra com heap) --------------------------------

// Result é o resultado da busca do caminho mais barato.
type Result struct {
	Path      []string `json:"path"`      // capitais na ordem, da origem ao destino
	Distance  int      `json:"distance"`  // distância total percorrida (km)
	FuelCost  float64  `json:"fuelCost"`  // gasto somente com combustível
	TollCost  float64  `json:"tollCost"`  // gasto somente com pedágios
	TotalCost float64  `json:"totalCost"` // combustível + pedágios
	Found     bool     `json:"found"`     // existe rota?
}

// CheapestPath encontra a rota de MENOR CUSTO entre origem e destino.
//
// O custo de uma viagem tem duas partes:
//
//   - Combustível de um trecho = (distância_km / autonomia) * preçoDoLitro
//   - Pedágio = cobrado ao CHEGAR em cada capital. A origem não paga (estamos
//     saindo dela); todas as capitais seguintes, inclusive o destino, pagam.
//
// O Dijkstra precisa de um único peso por aresta. Então embutimos o pedágio da
// capital de chegada dentro do peso da aresta:
//
//	peso(u → v) = combustível(u, v) + pedágio(v)
//
// O heap (fila de prioridade) garante que sempre expandimos primeiro a capital
// de menor custo acumulado — é isso que torna o Dijkstra correto e eficiente.
func (s *Service) CheapestPath(from, to string, fuelPrice, autonomy float64) (Result, error) {
	if !s.HasCity(from) || !s.HasCity(to) {
		return Result{}, ErrCityNotFound
	}
	if autonomy <= 0 || fuelPrice < 0 {
		return Result{}, errors.New("preço do combustível e autonomia devem ser positivos")
	}

	// cost[c] = menor custo já conhecido para chegar em c a partir da origem.
	cost := make(map[string]float64, len(s.adj))
	for city := range s.adj {
		cost[city] = math.Inf(1) // começa em "infinito": ainda não alcançado.
	}
	cost[from] = 0

	// prev[c] = de qual capital chegamos em c (para remontar o caminho no fim).
	prev := make(map[string]string)

	// Fila de prioridade: sempre devolve a capital de menor custo primeiro.
	pq := &priorityQueue{{city: from, cost: 0}}
	heap.Init(pq)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(pqItem)

		// Chegamos ao destino com o menor custo possível: podemos parar.
		if current.city == to {
			break
		}
		// Entrada desatualizada (já achamos um caminho melhor): ignora.
		if current.cost > cost[current.city] {
			continue
		}

		// Tenta melhorar o custo de cada vizinho (relaxamento das arestas).
		for neighbor, distance := range s.adj[current.city] {
			edgeCost := fuelCost(distance, fuelPrice, autonomy) + float64(s.tolls[neighbor])
			newCost := cost[current.city] + edgeCost
			if newCost < cost[neighbor] {
				cost[neighbor] = newCost
				prev[neighbor] = current.city
				heap.Push(pq, pqItem{city: neighbor, cost: newCost})
			}
		}
	}

	// Destino continuou em "infinito": não há rota.
	if math.IsInf(cost[to], 1) {
		return Result{Path: []string{}, Found: false}, nil
	}

	path := rebuildPath(prev, from, to)
	result := s.summarize(path, fuelPrice, autonomy)
	result.Found = true
	return result, nil
}

// fuelCost calcula o gasto de combustível de um trecho.
func fuelCost(distance int, fuelPrice, autonomy float64) float64 {
	return (float64(distance) / autonomy) * fuelPrice
}

// rebuildPath remonta o caminho seguindo o prev de trás para frente
// (destino → origem) e depois inverte para ficar origem → destino.
func rebuildPath(prev map[string]string, from, to string) []string {
	path := []string{to}
	for current := to; current != from; {
		current = prev[current]
		path = append([]string{current}, path...)
	}
	return path
}

// summarize percorre o caminho final e soma distância, combustível e pedágios.
func (s *Service) summarize(path []string, fuelPrice, autonomy float64) Result {
	result := Result{Path: path}
	for i := 1; i < len(path); i++ {
		from, to := path[i-1], path[i]
		distance := s.adj[from][to]
		result.Distance += distance
		result.FuelCost += fuelCost(distance, fuelPrice, autonomy)
		result.TollCost += float64(s.tolls[to]) // pedágio ao chegar (origem não paga)
	}
	result.TotalCost = result.FuelCost + result.TollCost
	return result
}

// --- Fila de prioridade (heap mínimo) usada pelo Dijkstra -------------------
// Implementa a interface container/heap.Interface, ordenando pelo menor custo.

type pqItem struct {
	city string
	cost float64
}

type priorityQueue []pqItem

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x any) {
	*pq = append(*pq, x.(pqItem))
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
