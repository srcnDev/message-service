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

// @title           Message Service API
// @version         1.0

// @contact.name   Sercan Yilmaz
// @contact.email  sercanyilmaz.dev@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes http https

// @tag.name messages
// @tag.description Message management operations

// @tag.name sender
// @tag.description Message sender job control operations

// @tag.name health
// @tag.description Health check endpoint

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
