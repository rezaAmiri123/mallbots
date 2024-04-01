package system

import "context"


type Service interface {
	AppConfig
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
