package agent

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	edatlog "github.com/rezaAmiri123/edatV2/log"
)

func (a *Agent) setupMonitoring() error {
	mux := chi.NewMux()
	mux.Use(middleware.Heartbeat("/liveness"))
	mux.Method("GET", "/metrics", promhttp.Handler())
	//a.setupSwagger(mux)
	webServer := &http.Server{
		Addr:    fmt.Sprintf("%s", a.config.Monitoring.Address),
		Handler: mux,
	}
	// a.container.AddSingleton(constants.HttpServerKey, func(c di.Container) (any, error) {
	// 	return webServer, nil
	// })
	//a.httpServer = webServer
	
	go func() {
		logger := edatlog.DefaultLogger
		logger.Info(fmt.Sprintf("run http at %s", a.config.Monitoring.Address))
		err := webServer.ListenAndServe()
		if err != nil {
			_ = a.Shutdown()
		}
	}()
	go func() {
		http.ListenAndServe(a.config.Monitoring.PprofAddress, nil)
	}()
	return nil

}

// http://localhost:6060/debug/pprof/
// go tool pprof -http localhost:8085 http://localhost:6060/debug/pprof/heap?debug=1