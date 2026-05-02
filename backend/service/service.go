package service

import (
	"errors"
	"sort"

	"github.com/dominikbraun/graph"

	"grafotb1/model"
)

var ErrVertexNotFound = errors.New("vértice não encontrado")

type Service struct {
	g      graph.Graph[string, string]
	actors map[string]struct{}
	movies map[string]struct{}
}

// New builds the actor-movie graph from the movies dataset.
// Edges are added in both directions so traversals can move freely
// from actor to movie and vice-versa.
func New(movies []model.Movie) *Service {
	g := graph.New(graph.StringHash, graph.Directed())
	actors := map[string]struct{}{}
	titles := map[string]struct{}{}

	for _, m := range movies {
		titles[m.Title] = struct{}{}
		_ = g.AddVertex(m.Title)
		for _, a := range m.Cast {
			actors[a] = struct{}{}
			_ = g.AddVertex(a)
			_ = g.AddEdge(m.Title, a)
			_ = g.AddEdge(a, m.Title)
		}
	}
	return &Service{g: g, actors: actors, movies: titles}
}

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

// AdjacencyMap exposes the underlying adjacency for the show endpoint.
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

// ShortestPath runs a classic BFS and returns the shortest path between
// from and to. Returns nil when no path exists.
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
				return reconstruct(parent, from, to), nil
			}
			queue = append(queue, n)
		}
	}
	return nil, nil
}

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

// MaxPathsCap caps the number of paths returned by AllPathsUpTo so
// the response stays within a size the browser can render.
const MaxPathsCap = 10000

// AllPathsUpTo enumerates every simple path (no repeated vertices) from
// `from` to `to` with edge count <= maxLen. Uses iterative deepening so
// shorter paths are always discovered before longer ones — if the cap
// kicks in at depth k, every path of length < k is guaranteed complete.
// Returns paths sorted by length ascending, then lexicographically.
func (s *Service) AllPathsUpTo(from, to string, maxLen int) ([][]string, bool, error) {
	if !s.HasVertex(from) || !s.HasVertex(to) {
		return nil, false, ErrVertexNotFound
	}
	if from == to {
		return [][]string{{from}}, false, nil
	}

	adj, _ := s.g.AdjacencyMap()

	// Pre-sort neighbours once so DFS expansion is deterministic and
	// the final result is already mostly ordered.
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

	// dfs collects simple paths from `node` to `to` of length exactly
	// `remaining` more edges. `remaining > 0` is the loop invariant.
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
				// `to` only counts when we've used exactly `remaining`
				// edges; landing on it earlier would mean a shorter path,
				// which a previous iteration already collected.
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
