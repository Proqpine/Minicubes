package main

import (
	"caldave/internal/config"
	"caldave/internal/server"
	"context"
	"fmt"
	"os"
)

func main() {
	cfg := config.NewConfig()
	ctx := context.Background()
	if err := server.Run(cfg, ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %s\n", err)
		os.Exit(1)
	}
}
