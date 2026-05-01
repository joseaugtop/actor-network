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

// AllShortestPaths returns every shortest path between from and to,
// but only when its length (in edges) is <= maxLen.
// Built on top of a BFS that records every predecessor at the same
// minimum distance, then walks the predecessor DAG to enumerate paths.
func (s *Service) AllShortestPaths(from, to string, maxLen int) ([][]string, error) {
	if !s.HasVertex(from) || !s.HasVertex(to) {
		return nil, ErrVertexNotFound
	}
	if from == to {
		return [][]string{{from}}, nil
	}

	adj, _ := s.g.AdjacencyMap()

	dist := map[string]int{from: 0}
	parents := map[string][]string{}
	queue := []string{from}
	target := -1

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		d := dist[v]
		if target != -1 && d >= target {
			continue
		}
		if d >= maxLen {
			continue
		}
		for n := range adj[v] {
			nd, seen := dist[n]
			if !seen {
				dist[n] = d + 1
				parents[n] = []string{v}
				if n == to {
					target = d + 1
				}
				queue = append(queue, n)
			} else if nd == d+1 {
				parents[n] = append(parents[n], v)
			}
		}
	}

	if target == -1 {
		return nil, nil
	}

	var paths [][]string
	var build func(node string, suffix []string)
	build = func(node string, suffix []string) {
		path := append([]string{node}, suffix...)
		if node == from {
			cp := make([]string, len(path))
			copy(cp, path)
			paths = append(paths, cp)
			return
		}
		for _, p := range parents[node] {
			build(p, path)
		}
	}
	build(to, nil)

	sort.Slice(paths, func(i, j int) bool {
		for k := 0; k < len(paths[i]) && k < len(paths[j]); k++ {
			if paths[i][k] != paths[j][k] {
				return paths[i][k] < paths[j][k]
			}
		}
		return false
	})
	return paths, nil
}
