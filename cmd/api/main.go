package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/app"
	"github.com/srcndev/message-service/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Config init failed: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		logger.Fatal("App init failed: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if runErr := application.Run(); runErr != nil {
			logger.Error("Server run error: %v", runErr)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Graceful shutdown failed: %v", err)
	}
}
