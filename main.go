package main

import (
	"flag"
	"fmt"

	// http.HandleFunc("/user", s.handleGetUserByID)

	"github.com/Courtcircuits/HackTheCrous.api/api"
	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

func main() {
	defaultPort := util.Get("PORT")
	listenAddr := flag.String("listenaddr", defaultPort, "server listen address")
	flag.Parse()

	fmt.Println("Starting server on port", *listenAddr)

	storage := storage.NewPostgresDatabase()

	server := api.NewServer(*listenAddr, *storage)
	// log.Fatal(server.Start())
	server.Start()
}
