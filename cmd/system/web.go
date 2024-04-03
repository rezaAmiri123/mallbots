package system

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

type Mux interface {
	Mux() *chi.Mux
}

func (s *System) initMux() {
	s.mux = chi.NewMux()
	s.mux.Use(middleware.Heartbeat("/liveness"))
	s.mux.Method("GET", "/metrics", promhttp.Handler())
}

func (s *System) Mux() *chi.Mux {
	return s.mux
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.Web.Http.Port),
		Handler: s.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("web server started; listening at http://localhost%s\n", s.cfg.Web.Http.Port)
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}
