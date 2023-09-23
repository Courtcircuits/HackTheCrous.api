package main

import (
	"flag"
	"fmt"

	// http.HandleFunc("/user", s.handleGetUserByID)

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Courtcircuits/HackTheCrous.api/api"
	"github.com/Courtcircuits/HackTheCrous.api/graph"
	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

func main() {
	defaultPort := ":" + util.Get("PORT")
	listenAddr := flag.String("listenaddr", defaultPort, "server listen address")
	flag.Parse()

	fmt.Println("Starting server on port", *listenAddr)

	storage := storage.NewPostgresDatabase()
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))

	server := api.NewServer(*listenAddr, *storage, h)
	server.Start()
}
