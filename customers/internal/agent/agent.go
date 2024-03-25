package agent

import (
	"io"
	"sync"

	"github.com/rezaAmiri123/edatV2/di"
)


type Agent struct {
	config Config
	container di.Container

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
	closers      []io.Closer
}

func NewAgent(config Config) (*Agent, error) {
	a := &Agent{
		config:    config,
		container: di.New(),
		shutdowns: make(chan struct{}),
	}
	setupsFn := []func() error{
		a.setupLogger,
		a.setupTracer,
		a.setupMonitoring,
		// a.setupRegistry,
		a.setupDatabase,
		// a.setupEventServer,
		a.setupApplication,
		// a.setupEventHandler,
		a.setupGrpcServer,
		//a.setupHttpServer,
	}
	for _, fn := range setupsFn {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	if a.shutdown {
		return nil
	}
	a.shutdown = true
	close(a.shutdowns)
	shutdown := []func() error{

		//func() error {
		//	a.grpcServer.GracefulStop()
		//	return nil
		//},
		//func() error {
		//	return a.httpServer.Shutdown(context.Background())
		//},
		// func() error {
		// 	tp := a.container.Get(constants.TracerKey).(*trace.TracerProvider)
		// 	return tp.Shutdown(context.Background())
		// },
		//func() error {
		//	stream := a.container.Get(constants.StreamKey).(am.MessageStream)
		//	return stream.Unsubscribe()
		//},

		//func() error {
		//	return a.jaegerCloser.Close()
		//},
		a.cleanupGrpcServer,
		a.cleanupDatabase,
		a.cleanupTracer,
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}
