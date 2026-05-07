// Lógica de grafo: construção da rede ator-filme e buscas (BFS e DFS).
package service

import (
	"errors"
	"fmt"
	"sort"

	"github.com/dominikbraun/graph"

	"grafotb1/model"
)

var ErrVertexNotFound = errors.New("vértice não encontrado")

// vertexLabel gera um rótulo único para o vértice usando o ID do filme.
// Ex: "Dune (841)" vs "Dune (438631)"
func vertexLabel(m model.Movie) string {
	return fmt.Sprintf("%s (%d)", m.Title, m.Id)
}

// cleanLabel remove o sufixo "(ID)" para exibição no frontend.
func cleanLabel(v string) string {
	n := len(v)
	if n < 3 || v[n-1] != ')' {
		return v
	}
	i := n - 2
	for i > 0 && v[i] != '(' {
		i--
	}
	if i <= 0 || v[i] != '(' || v[i-1] != ' ' {
		return v
	}
	return v[:i-1]
}

// --- Grafo ------------------------------------------------------------------

type Service struct {
	g      graph.Graph[string, string]
	actors map[string]bool
	movies map[string]string // key: label interno (com ID), value: título limpo
}

// Seed carrega os filmes no grafo, criando vértices e arestas.
func Seed(movies []model.Movie) *Service {
	g := graph.New(graph.StringHash)
	actors := map[string]bool{}
	titles := map[string]string{} // label interno -> título limpo

	for _, m := range movies {
		label := vertexLabel(m)
		titles[label] = m.Title
		g.AddVertex(label)
		for _, a := range m.Cast {
			actors[a] = true
			g.AddVertex(a)
			g.AddEdge(label, a)
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
// Os nós de filme no caminho são retornados com o título limpo (sem ID).
func (s *Service) ShortestPath(from, to string) ([]string, error) {
	if !s.HasVertex(from) || !s.HasVertex(to) {
		return nil, ErrVertexNotFound
	}
	if from == to {
		return []string{from}, nil
	}

	adj, _ := s.g.AdjacencyMap()

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
				return s.reconstructWithCleanLabels(parent, from, to), nil
			}
			queue = append(queue, n)
		}
	}
	return nil, nil
}

// reconstructWithCleanLabels reconstrói o caminho limpando os rótulos de filme.
func (s *Service) reconstructWithCleanLabels(parent map[string]string, from, to string) []string {
	var path []string
	cur := to
	for {
		path = append([]string{s.cleanNodeLabel(cur)}, path...)
		if cur == from {
			return path
		}
		cur = parent[cur]
	}
}

// cleanNodeLabel retorna o título limpo se for filme, ou o nome original se for ator.
func (s *Service) cleanNodeLabel(v string) string {
	if _, ok := s.movies[v]; ok {
		return s.movies[v]
	}
	return v
}

// --- DFS: todos os caminhos até maxLen --------------------------------------

// Limite para a resposta não estourar o que o navegador renderiza.
const MaxPathsCap = 10000

// Enumera todo caminho simples de from até to com até maxLen arestas.
// Os nós de filme nos caminhos são retornados com o título limpo (sem ID).
func (s *Service) AllPathsUpTo(from, to string, maxLen int) ([][]string, bool, error) {
	if !s.HasVertex(from) || !s.HasVertex(to) {
		return nil, false, ErrVertexNotFound
	}
	if from == to {
		return [][]string{{from}}, false, nil
	}

	adj, _ := s.g.AdjacencyMap()

	sortedNeighbours := make(map[string][]string, len(adj))
	for vertex, neighboursMap := range adj {
		neighbours := make([]string, 0, len(neighboursMap))
		for neighbour := range neighboursMap {
			neighbours = append(neighbours, neighbour)
		}
		sort.Strings(neighbours)
		sortedNeighbours[vertex] = neighbours
	}

	var paths [][]string
	truncated := false
	visited := map[string]bool{from: true}
	path := []string{from}

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

	// Ordena caminhos encontrados
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

	// Limpa os rótulos dos filmes nos resultados
	cleaned := make([][]string, len(paths))
	for i, p := range paths {
		cleaned[i] = make([]string, len(p))
		for j, node := range p {
			cleaned[i][j] = s.cleanNodeLabel(node)
		}
	}

	return cleaned, truncated, nil
}
