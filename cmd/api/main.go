package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("app init failed: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if runErr := application.Run(); runErr != nil {
			log.Printf("server run error: %v", runErr)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
}
