package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"grafotb1/service"
)

type Server struct {
	svc *service.Service
}

func New(svc *service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/actors", s.handleActors)
	mux.HandleFunc("/bfs", s.handleBFS)
	mux.HandleFunc("/bfs8", s.handleBFS8)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "mensagem": "servidor ativo"})
	})
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

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
	Paths   []pathItem `json:"paths"`
	Count   int        `json:"count"`
	Length  int        `json:"length,omitempty"`
	Found   bool       `json:"found"`
	Message string     `json:"message,omitempty"`
}

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

	paths, err := s.svc.AllShortestPaths(from, to, maxLen)
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
	writeJSON(w, http.StatusOK, pathsResponse{
		Paths:  items,
		Count:  len(items),
		Length: len(paths[0]) - 1,
		Found:  true,
	})
}

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
