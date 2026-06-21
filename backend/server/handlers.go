package server

import (
	"errors"
	"net/http"
	"strconv"

	"grafotb1/service"
)

// GET /capitais — lista de capitais (alimenta os selects do frontend).
func (s *Server) handleCapitais(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}
	capitais := s.svc.Capitals()
	writeJSON(w, http.StatusOK, map[string]any{
		"count":    len(capitais),
		"capitais": capitais,
	})
}

// GET /show — lista de adjacências: cada capital e seus vizinhos com a distância.
func (s *Server) handleShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}
	adj := s.svc.Show()
	writeJSON(w, http.StatusOK, map[string]any{
		"count":     len(adj),
		"adjacency": adj,
	})
}

// GET /caminho?origem=&destino=&combustivel=&autonomia=
// Devolve a rota de menor custo entre duas capitais.
func (s *Server) handleCaminho(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}

	q := r.URL.Query()
	origem := q.Get("origem")
	destino := q.Get("destino")
	if origem == "" || destino == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "parâmetros 'origem' e 'destino' são obrigatórios",
		})
		return
	}

	fuelPrice, errFuel := strconv.ParseFloat(q.Get("combustivel"), 64)
	autonomy, errAuto := strconv.ParseFloat(q.Get("autonomia"), 64)
	if errFuel != nil || errAuto != nil || fuelPrice <= 0 || autonomy <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "'combustivel' e 'autonomia' devem ser números positivos",
		})
		return
	}

	res, err := s.svc.CheapestPath(origem, destino, fuelPrice, autonomy)
	if errors.Is(err, service.ErrCityNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"found":   false,
			"message": "capital de origem ou destino não encontrada",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !res.Found {
		writeJSON(w, http.StatusOK, map[string]any{
			"found":   false,
			"message": "não existe rota entre as capitais selecionadas",
		})
		return
	}
	writeJSON(w, http.StatusOK, res)
}
