package config

import (
	"github.com/srcndev/message-service/pkg/errs"
)

var (
	errAppPortEmpty = errs.NewWithDefaults("APP_PORT_EMPTY", "APP_PORT cannot be empty")
	errAppURLEmpty  = errs.NewWithDefaults("APP_URL_EMPTY", "APP_URL cannot be empty")
)
