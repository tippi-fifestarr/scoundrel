package main

import (
	"log"
	"os"

	"github.com/tippi-fifestarr/scoundrel/api"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := api.NewServer()
	log.Fatal(server.Start(":" + port))
}
