package server

import (
	"errors"
	"math"
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

	origem, destino, fuelPrice, autonomy, ok := readCaminhoParams(w, r)
	if !ok {
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

// GET /comparar?origem=&destino=&combustivel=&autonomia=
// Roda o NOSSO Dijkstra e o da biblioteca dominikbraun/graph e devolve os dois
// resultados lado a lado. Endpoint de conferência: se "custosBatem" for true,
// nossa implementação concorda com uma biblioteca consagrada.
func (s *Server) handleComparar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método não permitido"})
		return
	}

	origem, destino, fuelPrice, autonomy, ok := readCaminhoParams(w, r)
	if !ok {
		return
	}

	meu, err := s.svc.CheapestPath(origem, destino, fuelPrice, autonomy)
	if errors.Is(err, service.ErrCityNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"error": "capital de origem ou destino não encontrada",
		})
		return
	}
	lib, _ := s.svc.ShortestPathLib(origem, destino, fuelPrice, autonomy)

	// Custos batem se ambos acharam (ou ambos não acharam) e o valor coincide.
	custosBatem := meu.Found == lib.Found && math.Abs(meu.TotalCost-lib.TotalCost) < 1e-6

	writeJSON(w, http.StatusOK, map[string]any{
		"origem":         origem,
		"destino":        destino,
		"combustivel":    fuelPrice,
		"autonomia":      autonomy,
		"meu":            meu,
		"lib":            lib,
		"custosBatem":    custosBatem,
		"diferencaCusto": math.Abs(meu.TotalCost - lib.TotalCost),
		"mesmoCaminho":   iguaisPath(meu.Path, lib.Path),
	})
}

// --- Helpers ----------------------------------------------------------------

// readCaminhoParams lê e valida os 4 parâmetros usados pelos endpoints de rota.
// Em caso de erro, já responde 400 e devolve ok=false.
func readCaminhoParams(w http.ResponseWriter, r *http.Request) (origem, destino string, fuelPrice, autonomy float64, ok bool) {
	q := r.URL.Query()
	origem = q.Get("origem")
	destino = q.Get("destino")
	if origem == "" || destino == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "parâmetros 'origem' e 'destino' são obrigatórios",
		})
		return
	}

	var errFuel, errAuto error
	fuelPrice, errFuel = strconv.ParseFloat(q.Get("combustivel"), 64)
	autonomy, errAuto = strconv.ParseFloat(q.Get("autonomia"), 64)
	if errFuel != nil || errAuto != nil || fuelPrice <= 0 || autonomy <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "'combustivel' e 'autonomia' devem ser números positivos",
		})
		return
	}

	ok = true
	return
}

// iguaisPath compara duas rotas posição a posição.
func iguaisPath(a, b []string) bool {
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
