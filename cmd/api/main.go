package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"basic-gin/internal/server"
)

func main() {
	log.Println("Starting application...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
