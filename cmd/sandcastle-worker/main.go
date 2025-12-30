package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mbobrovskyi/kube-sandcastle/internal/executor"
	"github.com/mbobrovskyi/kube-sandcastle/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	exec := executor.NewExecutor()

	workerHandler := server.NewWorkerHandler(exec)
	router := server.NewRouter(workerHandler)

	server := server.NewServer(router)
	if err := server.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
