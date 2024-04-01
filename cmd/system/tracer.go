package system

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)
func (s *System) initOpenTelemetry() error {
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		return err
	}

	s.tp = sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(s.tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	s.waiter.Cleanup(func() {
		if err := s.tp.Shutdown(context.Background()); err != nil {
			s.logger.Error().Err(err).Msg("ran into an issue shutting down the tracer provider")
		}
	})

	return nil
}
