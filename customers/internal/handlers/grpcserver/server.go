package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
	"github.com/rezaAmiri123/mallbots/customers/internal/application"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	customerspb.UnimplementedCustomersServiceServer
	app application.App
}

var _ customerspb.CustomersServiceServer = (*server)(nil)

func RegisterServer(app application.App, register grpc.ServiceRegistrar) {
	customerspb.RegisterCustomersServiceServer(register, server{
		app: app,
	})
}

func (s server) RegisterCustomer(ctx context.Context, request *customerspb.RegisterCustomerRequest) (resp *customerspb.RegisterCustomerResponse, err error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("CustomerID", id),
	)

	err = s.app.RegisterCustomer(ctx, application.RegisterCustomer{
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

func (s server) GetCustomer(ctx context.Context, request *customerspb.GetCustomerRequest) (resp *customerspb.GetCustomerResponse, err error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CustomerID", request.GetId()),
	)

	customer, err := s.app.GetCustomer(ctx, application.GetCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	resp = &customerspb.GetCustomerResponse{
		Customer: s.customerToProto(customer),
	}

	return resp, nil
}

func (s server) customerToProto(customer *domain.Customer) *customerspb.Customer {
	return &customerspb.Customer{
		Id:        customer.ID(),
		Name:      customer.Name,
		SmsNumber: customer.SmsNumber,
		Enabled:   customer.Enabled,
	}
}
