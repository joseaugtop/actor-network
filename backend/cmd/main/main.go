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
	data, err := os.ReadFile("api/latest_movies.json")
	if err != nil {
		log.Fatalf("erro ao ler o arquivo latest_movies.json: %v", err)
	}

	var movies []model.Movie
	if err := json.Unmarshal(data, &movies); err != nil {
		log.Fatalf("erro ao processar o JSON de filmes: %v", err)
	}

	svc := service.New(movies)
	log.Printf("grafo construído — %d filmes, %d atores únicos", len(movies), len(svc.Actors()))

	srv := server.New(svc)
	log.Println("servidor iniciado em http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", srv.Routes()))
}
