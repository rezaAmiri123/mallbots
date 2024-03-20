package agent

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	edatlog "github.com/rezaAmiri123/edatV2/log"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
)

func (a *Agent) setupHttpServer() error {
	mux := chi.NewMux()
	mux.Use(middleware.Heartbeat("/liveness"))
	mux.Method("GET", "/metrics", promhttp.Handler())
	//a.setupSwagger(mux)
	webServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ShutdownTimeout),
		Handler: mux,
	}
	a.container.AddSingleton(constants.HttpServerKey, func(c di.Container) (any, error) {
		return webServer, nil
	})
	//a.httpServer = webServer
	
	go func() {
		logger := edatlog.DefaultLogger
		logger.Info("run http at %s", webServer.Addr)
		err := webServer.ListenAndServe()
		if err != nil {
			_ = a.Shutdown()
		}
	}()
	return nil

}
