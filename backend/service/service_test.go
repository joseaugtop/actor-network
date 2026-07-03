package service

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"grafotb1/model"
)

const eps = 1e-6

// --- Exemplo calculado à mão -------------------------------------------------
//
// Grafo pequeno onde dá pra conferir o custo na ponta do lápis:
//
//	tolls: A=10  B=5  C=20
//	arestas (km): A-B=100, B-C=100, A-C=250
//	combustível=R$2/L, autonomia=10 km/L  ->  R$0,20 por km
//
// O destino NÃO paga pedágio (fim da viagem); só as intermediárias pagam.
// Caminho A->B->C: combustível (200 km * 0,20) = 40 ; pedágio  B (só ela) = 5 ; total = 45
// Caminho A->C:    combustível (250 km * 0,20) = 50 ; pedágio  nenhum        = 0 ; total = 50
//
// Logo o mais barato A->C é A->B->C, custo 45.
func TestCheapestPath_ExemploManual(t *testing.T) {
	s := &Service{
		adj:   map[string]map[string]int{},
		tolls: map[string]int{"A": 10, "B": 5, "C": 20},
	}
	s.addEdge("A", "B", 100)
	s.addEdge("B", "C", 100)
	s.addEdge("A", "C", 250)

	res, err := s.CheapestPath("A", "C", 2, 10)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !res.Found {
		t.Fatal("esperava encontrar rota")
	}
	wantPath := []string{"A", "B", "C"}
	if !equal(res.Path, wantPath) {
		t.Errorf("path = %v, quero %v", res.Path, wantPath)
	}
	assertFloat(t, "total", res.TotalCost, 45)
	assertFloat(t, "fuel", res.FuelCost, 40)
	assertFloat(t, "toll", res.TollCost, 5)
	if res.Distance != 200 {
		t.Errorf("distance = %d, quero 200", res.Distance)
	}
}

func TestCheapestPath_SemRota(t *testing.T) {
	s := &Service{
		adj:   map[string]map[string]int{},
		tolls: map[string]int{"A": 0, "B": 0, "Ilha": 0},
	}
	s.addEdge("A", "B", 10)
	s.adj["Ilha"] = map[string]int{} // capital isolada, sem arestas

	res, err := s.CheapestPath("A", "Ilha", 5, 10)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if res.Found {
		t.Errorf("não deveria haver rota até a capital isolada, veio %+v", res)
	}
}

func TestCheapestPath_OrigemIgualDestino(t *testing.T) {
	s := loadRealService(t)
	res, err := s.CheapestPath("Manaus", "Manaus", 5, 10)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !res.Found || res.TotalCost != 0 || len(res.Path) != 1 {
		t.Errorf("origem==destino deveria custar 0 e 1 capital, veio %+v", res)
	}
}

// --- Conferência contra Bellman-Ford no grafo REAL ---------------------------
//
// Bellman-Ford é um algoritmo DIFERENTE do Dijkstra mas resolve o mesmo
// problema (menor custo) usando o mesmo peso de aresta. Se os dois concordam
// para TODOS os pares de capitais e em vários cenários de preço/autonomia,
// temos confiança forte de que o caminho do Dijkstra é realmente o mais barato.
func TestCheapestPath_BateComBellmanFord(t *testing.T) {
	s := loadRealService(t)
	capitais := s.Capitals()

	cenarios := []struct {
		fuel, autonomy float64
	}{
		{5.89, 12},  // carro econômico
		{7.20, 8},   // combustível caro, carro "beberrão"
		{1.00, 1},   // peso do pedágio fica pequeno perto do combustível
		{0.01, 100}, // combustível quase irrelevante: domina o pedágio
	}

	for _, c := range cenarios {
		for _, origem := range capitais {
			// Bellman-Ford resolve a origem para TODOS os destinos de uma vez.
			ref := bellmanFord(s, origem, c.fuel, c.autonomy)

			for _, destino := range capitais {
				res, err := s.CheapestPath(origem, destino, c.fuel, c.autonomy)
				if err != nil {
					t.Fatalf("%s->%s: erro %v", origem, destino, err)
				}

				want := ref[destino]
				if math.IsInf(want, 1) {
					if res.Found {
						t.Errorf("%s->%s: Dijkstra achou rota, Bellman-Ford diz que não existe", origem, destino)
					}
					continue
				}
				if !res.Found {
					t.Errorf("%s->%s: Dijkstra não achou rota, Bellman-Ford achou custo %.4f", origem, destino, want)
					continue
				}
				// Bellman-Ford embute o pedágio do destino no custo; o serviço não
				// cobra pedágio no destino. Tira o do destino para comparar (exceto
				// quando origem==destino, que já custa 0).
				if destino != origem {
					want -= float64(s.tolls[destino])
				}
				if math.Abs(res.TotalCost-want) > eps {
					t.Errorf("%s->%s: custo Dijkstra=%.6f, Bellman-Ford=%.6f", origem, destino, res.TotalCost, want)
				}

				// O caminho devolvido tem que ser válido e recalcular o mesmo custo.
				checkPathValido(t, s, res, origem, destino, c.fuel, c.autonomy)
			}
		}
	}
}

// bellmanFord devolve o menor custo da origem até cada capital, com o mesmo
// modelo de custo do serviço: peso(u->v) = combustível(u,v) + pedágio(v).
func bellmanFord(s *Service, origem string, fuel, autonomy float64) map[string]float64 {
	dist := make(map[string]float64, len(s.adj))
	for city := range s.adj {
		dist[city] = math.Inf(1)
	}
	dist[origem] = 0

	// Relaxa todas as arestas |V|-1 vezes.
	for i := 0; i < len(s.adj)-1; i++ {
		changed := false
		for u, vizinhos := range s.adj {
			if math.IsInf(dist[u], 1) {
				continue
			}
			for v, km := range vizinhos {
				w := fuelCost(km, fuel, autonomy) + float64(s.tolls[v])
				if dist[u]+w < dist[v] {
					dist[v] = dist[u] + w
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	return dist
}

// checkPathValido confere que o caminho é real (capitais consecutivas são
// vizinhas), que começa/termina certo e que somar os custos do caminho bate
// com o TotalCost devolvido.
func checkPathValido(t *testing.T, s *Service, res Result, origem, destino string, fuel, autonomy float64) {
	t.Helper()
	if len(res.Path) == 0 {
		t.Errorf("%s->%s: caminho vazio", origem, destino)
		return
	}
	if res.Path[0] != origem || res.Path[len(res.Path)-1] != destino {
		t.Errorf("%s->%s: caminho começa/termina errado: %v", origem, destino, res.Path)
	}

	var total float64
	for i := 1; i < len(res.Path); i++ {
		a, b := res.Path[i-1], res.Path[i]
		km, ok := s.adj[a][b]
		if !ok {
			t.Errorf("%s->%s: %s e %s não são vizinhas no caminho %v", origem, destino, a, b, res.Path)
			return
		}
		total += fuelCost(km, fuel, autonomy)
		if i < len(res.Path)-1 { // destino não paga pedágio
			total += float64(s.tolls[b])
		}
	}
	if math.Abs(total-res.TotalCost) > eps {
		t.Errorf("%s->%s: recalculando o caminho deu %.6f, mas TotalCost=%.6f", origem, destino, total, res.TotalCost)
	}
}

// --- helpers -----------------------------------------------------------------

func loadRealService(t *testing.T) *Service {
	t.Helper()
	data, err := os.ReadFile("../api/capitais.json")
	if err != nil {
		t.Fatalf("não consegui ler capitais.json: %v", err)
	}
	var entries []map[string]model.City
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	return Seed(entries)
}

func assertFloat(t *testing.T, nome string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > eps {
		t.Errorf("%s = %.6f, quero %.6f", nome, got, want)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
