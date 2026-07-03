// Lê o dataset, monta o grafo e sobe o servidor HTTP.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"grafotb1/model"
	"grafotb1/server"
	"grafotb1/service"
)

func main() {
	data, err := os.ReadFile("api/capitais.json")
	if err != nil {
		log.Fatalf("erro ao ler o arquivo capitais.json: %v", err)
	}

	var cities []map[string]model.City
	if err := json.Unmarshal(data, &cities); err != nil {
		log.Fatalf("erro ao processar o JSON de cidades: %v", err)
	}

	svc := service.Seed(cities)
	log.Printf("grafo construído — %d capitais", len(svc.Capitals()))

	srv := server.New(svc)
	log.Println("servidor iniciado em http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", srv.Routes()))
}
