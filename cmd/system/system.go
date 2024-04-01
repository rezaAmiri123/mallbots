package system

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/rezaAmiri123/edatV2/waiter"
	"github.com/rezaAmiri123/mallbots/internal/config"
	"github.com/rs/zerolog"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type System struct {
	cfg    config.Config
	db     *sql.DB
	nc     *nats.Conn
	js     nats.JetStreamContext
	mux    *chi.Mux
	rpc    *grpc.Server
	waiter waiter.Waiter
	logger zerolog.Logger
	tp     *sdktrace.TracerProvider
}

func NewSystem(cfg config.Config) (*System, error) {
	s := &System{cfg: cfg}

	s.initWaiter()

	if err := s.initLogger();err!=nil{
		return nil, err
	}
	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initJS(); err != nil {
		return nil, err
	}

	if err := s.initOpenTelemetry(); err != nil {
		return nil, err
	}

	s.initMux()
	s.initRpc()
	

	return s, nil
}


func (s *System) Config() config.Config {
	return s.cfg
}
