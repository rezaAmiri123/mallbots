package system

import (
	"context"

	"github.com/rezaAmiri123/mallbots/internal/config"
)


type Service interface {
	config.AppConfig
	DB
	JS
	Mux
	RPC
	Logger
	Waiter
}

type Module interface {
	Startup(context.Context, Service) error
}
