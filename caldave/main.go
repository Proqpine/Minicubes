package main

import (
	"caldave/internal/config"
	"caldave/internal/server"
	"log"
)

func main() {
	cfg := config.NewConfig()
	if err := server.Run(cfg); err != nil {
		log.Fatalf("could not run the server: %v", err)
	}
}
