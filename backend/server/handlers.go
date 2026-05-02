package server

import (
	"errors"
	"net/http"
	"strconv"

	"grafotb1/service"
)

// --- Tipos de resposta ------------------------------------------------------

type pathResponse struct {
	Path    []string `json:"path"`
	Length  int      `json:"length"`
	Found   bool     `json:"found"`
	Message string   `json:"message,omitempty"`
}

type pathItem struct {
	Path   []string `json:"path"`
	Length int      `json:"length"`
}

type pathsResponse struct {
	Paths     []pathItem `json:"paths"`
	Count     int        `json:"count"`
	MinLength int        `json:"minLength,omitempty"`
	MaxLength int        `json:"maxLength,omitempty"`
	Truncated bool       `json:"truncated,omitempty"`
	Cap       int        `json:"cap,omitempty"`
	Found     bool       `json:"found"`
	Message   string     `json:"message,omitempty"`
}

// --- Handlers ---------------------------------------------------------------

// GET /show — adjacências completas do grafo.
func (s *Server) handleShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}
	adj := s.svc.AdjacencyMap()
	writeJSON(w, http.StatusOK, map[string]any{
		"count":     len(adj),
		"adjacency": adj,
	})
}

// GET /actors — lista ordenada de atores.
func (s *Server) handleActors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}
	actors := s.svc.Actors()
	writeJSON(w, http.StatusOK, map[string]any{
		"count":  len(actors),
		"actors": actors,
	})
}

// GET /bfs?from=&to= — caminho mínimo entre dois atores.
func (s *Server) handleBFS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}
	from, to, ok := readEndpoints(w, r)
	if !ok {
		return
	}

	path, err := s.svc.ShortestPath(from, to)
	if errors.Is(err, service.ErrVertexNotFound) {
		writeJSON(w, http.StatusNotFound, pathResponse{
			Path: []string{}, Length: -1, Found: false,
			Message: "ator não encontrado",
		})
		return
	}
	if path == nil {
		writeJSON(w, http.StatusOK, pathResponse{
			Path: []string{}, Length: -1, Found: false,
			Message: "nenhum relacionamento encontrado",
		})
		return
	}
	writeJSON(w, http.StatusOK, pathResponse{
		Path:   path,
		Length: len(path) - 1,
		Found:  true,
	})
}

// GET /bfs8?from=&to=&max= — todos os caminhos até max arestas (padrão 8).
func (s *Server) handleBFS8(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}
	from, to, ok := readEndpoints(w, r)
	if !ok {
		return
	}

	maxLen := 8
	if v := r.URL.Query().Get("max"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxLen = n
		}
	}

	paths, truncated, err := s.svc.AllPathsUpTo(from, to, maxLen)
	if errors.Is(err, service.ErrVertexNotFound) {
		writeJSON(w, http.StatusNotFound, pathsResponse{
			Paths: []pathItem{}, Found: false, Message: "ator não encontrado",
		})
		return
	}
	if len(paths) == 0 {
		writeJSON(w, http.StatusOK, pathsResponse{
			Paths: []pathItem{}, Count: 0, Found: false,
			Message: "nenhum relacionamento encontrado dentro do comprimento máximo",
		})
		return
	}

	items := make([]pathItem, len(paths))
	for i, p := range paths {
		items[i] = pathItem{Path: p, Length: len(p) - 1}
	}
	resp := pathsResponse{
		Paths:     items,
		Count:     len(items),
		MinLength: len(paths[0]) - 1,
		MaxLength: len(paths[len(paths)-1]) - 1,
		Truncated: truncated,
		Found:     true,
	}
	if truncated {
		resp.Cap = service.MaxPathsCap
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Helpers ----------------------------------------------------------------

func readEndpoints(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "parâmetros 'from' e 'to' são obrigatórios",
		})
		return "", "", false
	}
	return from, to, true
}
