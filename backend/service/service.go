// Lógica de grafo: construção da rede ator-filme e buscas (BFS e DFS).
package service

import (
	"errors"
	"sort"

	"github.com/dominikbraun/graph"

	"grafotb1/model"
)

var ErrVertexNotFound = errors.New("vértice não encontrado")

// --- Grafo ------------------------------------------------------------------

type Service struct {
	g      graph.Graph[string, string]
	actors map[string]bool
	movies map[string]bool
}

// Seed carrega os filmes no grafo, criando vértices e arestas.
func Seed(movies []model.Movie) *Service {
	g := graph.New(graph.StringHash)
	actors := map[string]bool{}
	titles := map[string]bool{}

	for _, m := range movies {
		titles[m.Title] = true
		_ = g.AddVertex(m.Title)
		for _, a := range m.Cast {
			actors[a] = true
			_ = g.AddVertex(a)
			_ = g.AddEdge(m.Title, a)
		}
	}
	return &Service{g: g, actors: actors, movies: titles}
}

// --- Consultas --------------------------------------------------------------

func (s *Service) Actors() []string {
	out := make([]string, 0, len(s.actors))
	for a := range s.actors {
		out = append(out, a)
	}
	sort.Strings(out)
	return out
}

func (s *Service) HasVertex(v string) bool {
	_, err := s.g.Vertex(v)
	return err == nil
}

func (s *Service) AdjacencyMap() map[string][]string {
	adj, _ := s.g.AdjacencyMap()
	out := make(map[string][]string, len(adj))
	for v, ns := range adj {
		neighbours := make([]string, 0, len(ns))
		for n := range ns {
			neighbours = append(neighbours, n)
		}
		sort.Strings(neighbours)
		out[v] = neighbours
	}
	return out
}

// --- BFS: caminho mínimo ----------------------------------------------------

// Retorna o menor caminho entre from e to, ou nil se não existir.
func (s *Service) ShortestPath(from, to string) ([]string, error) {
	if !s.HasVertex(from) || !s.HasVertex(to) {
		return nil, ErrVertexNotFound
	}
	if from == to {
		return []string{from}, nil
	}

	adj, _ := s.g.AdjacencyMap()

	// parent guarda de onde cada vértice foi alcançado e também
	// faz o papel de "visitados".
	parent := map[string]string{from: ""}
	queue := []string{from}

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		for n := range adj[v] {
			if _, seen := parent[n]; seen {
				continue
			}
			parent[n] = v
			if n == to {
				return reconstruct(parent, from, to), nil
			}
			queue = append(queue, n)
		}
	}
	return nil, nil
}

// reconstrói o caminho do destino até a origem usando o mapa parent, invertendo para a ordem correta.
func reconstruct(parent map[string]string, from, to string) []string {
	var path []string
	cur := to
	for {
		path = append([]string{cur}, path...)
		if cur == from {
			return path
		}
		cur = parent[cur]
	}
}

// --- DFS: todos os caminhos até maxLen --------------------------------------

// Limite para a resposta não estourar o que o navegador renderiza.
const MaxPathsCap = 10000

// Enumera todo caminho simples de from até to com até maxLen arestas.
// Faz aprofundamento iterativo (d = 1..maxLen), então caminhos menores
// vêm primeiro e a truncagem só corta os mais longos.
func (s *Service) AllPathsUpTo(from, to string, maxLen int) ([][]string, bool, error) {
	if !s.HasVertex(from) || !s.HasVertex(to) {
		return nil, false, ErrVertexNotFound
	}
	if from == to {
		return [][]string{{from}}, false, nil
	}

	adj, _ := s.g.AdjacencyMap()

	// Vizinhos pré-ordenados, expansão determinística.
	sortedNeighbours := make(map[string][]string, len(adj))
	for v, ns := range adj {
		ss := make([]string, 0, len(ns))
		for n := range ns {
			ss = append(ss, n)
		}
		sort.Strings(ss)
		sortedNeighbours[v] = ss
	}

	var paths [][]string
	truncated := false
	visited := map[string]bool{from: true}
	path := []string{from}

	// dfs busca caminhos com exatamente `remaining` arestas restantes.
	var dfs func(node string, remaining int)
	dfs = func(node string, remaining int) {
		if truncated {
			return
		}
		for _, n := range sortedNeighbours[node] {
			if visited[n] {
				continue
			}
			if remaining == 1 {
				if n == to {
					cp := make([]string, len(path)+1)
					copy(cp, path)
					cp[len(path)] = n
					paths = append(paths, cp)
					if len(paths) >= MaxPathsCap {
						truncated = true
						return
					}
				}
				continue
			}
			if n == to {
				// chegar em `to` antes da última aresta = caminho
				// mais curto, já pego em iteração anterior.
				continue
			}
			visited[n] = true
			path = append(path, n)
			dfs(n, remaining-1)
			path = path[:len(path)-1]
			delete(visited, n)
			if truncated {
				return
			}
		}
	}

	for d := 1; d <= maxLen && !truncated; d++ {
		dfs(from, d)
	}

	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}
		for k := 0; k < len(paths[i]); k++ {
			if paths[i][k] != paths[j][k] {
				return paths[i][k] < paths[j][k]
			}
		}
		return false
	})
	return paths, truncated, nil
}
