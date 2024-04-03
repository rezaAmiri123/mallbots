package system

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

func (s *System) WaitForMonitoring(ctx context.Context) error {
	mux := chi.NewMux()
	mux.Use(middleware.Heartbeat("/liveness"))
	mux.Method("GET", "/metrics", promhttp.Handler())

	webServer := &http.Server{
		Addr:    fmt.Sprintf("%s", s.cfg.Monitoring.Address),
		Handler: mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("monitoring server started; listening at")
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("error: ", err.Error())
			return err
		}
		return nil
	})
	group.Go(func() error {
		fmt.Printf("pprop server started; listening at")
		defer fmt.Println("pprop server shutdown")
		if err := http.ListenAndServe(s.cfg.Monitoring.PprofAddress, nil); err != nil && err != http.ErrServerClosed {
			fmt.Println("error: ", err.Error())
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("monitoring server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

//  http://localhost:6060/debug/pprof/
//  go tool pprof -http localhost:8085 http://localhost:6060/debug/pprof/heap?debug=1
