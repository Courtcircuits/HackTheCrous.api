package main

import (

	// http.HandleFunc("/user", s.handleGetUserByID)

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Courtcircuits/HackTheCrous.api/api"
	"github.com/Courtcircuits/HackTheCrous.api/graph"
	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

func main() {
	port := util.Get("PORT")

	store := storage.NewPostgresDatabase()
	cache := storage.NewRedisCache()
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))

	server := api.NewServer(":"+port, *store, cache, h)
	server.Start()
}
