package app

import (
	"github.com/srcndev/message-service/config"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
}

func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		Config: cfg,
	}

	return container, nil
}
