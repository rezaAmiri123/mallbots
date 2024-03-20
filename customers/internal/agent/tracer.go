package agent

import (
	"context"

	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	edatlog "github.com/rezaAmiri123/edatV2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func (a *Agent) setupTracer() error {
	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure())
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		//sdktrace.WithSampler(sdktrace.AlwaysSample()),
		//sdktrace.WithResource(resource.NewSchemaless(attribute.String("service.name", "myService"))),
		//sdktrace.WithSyncer(exp),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	a.container.AddSingleton(constants.TracerKey, func(c di.Container) (any, error) {
		return tp, nil
	})

	return nil
}
func (a *Agent) cleanupTracer() error {
	tp := a.container.Get(constants.TracerKey).(*sdktrace.TracerProvider)
	logger := edatlog.DefaultLogger
	if err := tp.Shutdown(context.Background()); err != nil {
		logger.Error("ran into an issue shutting down the tracer provider", edatlog.Error(err))
	}
}