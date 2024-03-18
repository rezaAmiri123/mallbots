package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
	"github.com/rezaAmiri123/mallbots/customers/internal/application"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	customerspb.UnimplementedCustomersServiceServer
	app application.App
}

var _ customerspb.CustomersServiceServer = (*GrpcHandler)(nil)

func RegisterServer(app application.App, register grpc.ServiceRegistrar) {
	customerspb.RegisterCustomersServiceServer(register, GrpcHandler{
		app: app,
	})
}

func (h GrpcHandler) RegisterCustomer(ctx context.Context, request *customerspb.RegisterCustomerRequest) (resp *customerspb.RegisterCustomerResponse, err error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("CustomerID", id),
	)

	err = h.app.RegisterCustomer(ctx, application.RegisterCustomer{
		ID:        id,
		Name:      request.GetName(),
		SmsNumber: request.GetSmsNumber(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	resp = &customerspb.RegisterCustomerResponse{
		Id: id,
	}

	return resp, nil
}
