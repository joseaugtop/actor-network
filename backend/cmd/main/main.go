package main

import (
	"encoding/json"
	"fmt"
	"grafotb1/model"
	"os"

	"github.com/dominikbraun/graph"
)

func showMovies(movies []model.Movie){
	for _, m := range movies {
        fmt.Printf("id:", m.Id, "title:",m.Title, "cast:",m.Cast,)
    }
}

//seed
func populateGraphFromArray(movies []model.Movie) graph.Graph[string, string] {
    g := graph.New(graph.StringHash)

    for _, movie := range movies {
        g.AddVertex(movie.Title)
        for _, actor := range movie.Cast {
            g.AddVertex(actor)
            g.AddEdge(movie.Title, actor)
        }
    }

    return g
}

//show
func showGraph(g graph.Graph[string, string]) {
    adjacencyMap, _ := g.AdjacencyMap()

    for vertex, adjacent  := range adjacencyMap {
		fmt.Println("{")
        fmt.Printf("\tVértice: %s,\n\tAssociados: {\n", vertex)
        for neighbor := range adjacent {
            fmt.Printf("\t\t%s\n", neighbor)

        }
        fmt.Printf("\t},\n")
        fmt.Println("},")
        fmt.Println()
    }
}

func main() {
	data, err := os.ReadFile("../api/latest_movies.json")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }

	var movies []model.Movie
    if err := json.Unmarshal(data, &movies); err != nil {
        fmt.Printf("Error: %v\n", err)
    }

	g := populateGraphFromArray(movies)	

	showGraph(g)

}