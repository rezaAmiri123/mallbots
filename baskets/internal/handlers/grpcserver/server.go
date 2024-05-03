package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/baskets/basketspb"
	"github.com/rezaAmiri123/mallbots/baskets/internal/application"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	basketspb.UnimplementedBasketServiceServer
}

var _ basketspb.BasketServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	basketspb.RegisterBasketServiceServer(registrar, server{
		app: app,
	})
	return nil
}

func (s server) StartBasket(ctx context.Context, request *basketspb.StartBasketRequest) (resp *basketspb.StartBasketResponse, err error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("BasketID", id),
		attribute.String("CustomerID", request.GetCustomerId()),
	)

	err = s.app.StartBasket(ctx, application.StartBasket{
		ID:         id,
		CustomerID: request.CustomerId,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &basketspb.StartBasketResponse{Id: id}, err
}

func (s server) CheckoutBasket(ctx context.Context, request *basketspb.CheckoutBasketRequest) (resp *basketspb.CheckoutBasketResponse, err error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("BasketID", request.GetId()),
		attribute.String("PaymentID", request.GetPaymentId()),
	)

	err = s.app.CheckoutBasket(ctx, application.CheckoutBasket{
		ID:         request.GetId(),
		PaymentID: request.GetPaymentId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &basketspb.CheckoutBasketResponse{}, err

}

func (s server) AddItem(ctx context.Context, request *basketspb.AddItemRequest) (resp *basketspb.AddItemResponse, err error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("BasketID", request.GetId()),
		attribute.String("ProductID", request.GetProductId()),
	)

	err = s.app.AddItem(ctx, application.AddItem{
		ID:        request.GetId(),
		ProductID: request.GetProductId(),
		Quantity:  int(request.GetQuantity()),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &basketspb.AddItemResponse{}, nil
}
