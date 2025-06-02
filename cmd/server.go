package main

import (
	"context"
	"log"

	"github.com/mesameen/proglog/internal/config"
	"github.com/mesameen/proglog/internal/server"
)

func main() {
	config.LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := server.StartHTTPServer(ctx); err != nil {
		log.Fatal(err)
	}
}
