package main

import (
	"flag"
	"log"

	"github.com/igodwin/secretsanta/internal/api"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	server := api.NewServer(*addr)

	log.Printf("Secret Santa Web Server")
	log.Printf("Visit http://localhost%s to get started", *addr)

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
