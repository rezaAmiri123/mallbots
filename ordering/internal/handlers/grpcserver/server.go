package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezaAmiri123/mallbots/internal/errorsotel"
	"github.com/rezaAmiri123/mallbots/ordering/internal/application"
	"github.com/rezaAmiri123/mallbots/ordering/internal/application/commands"
	"github.com/rezaAmiri123/mallbots/ordering/internal/domain"
	"github.com/rezaAmiri123/mallbots/ordering/orderingpb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	orderingpb.UnimplementedOrderingServiceServer
}

var _ orderingpb.OrderingServiceServer = (*server)(nil)

func RegisterServer(app application.App, register grpc.ServiceRegistrar) error {
	orderingpb.RegisterOrderingServiceServer(register, server{app: app})
	return nil
}

func (s server) CreateOrder(ctx context.Context, request *orderingpb.CreateOrderRequest) (*orderingpb.CreateOrderResponse, error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("OrderID", id),
		attribute.String("CustomerID", request.GetCustomerId()),
		attribute.String("PaymentID", request.GetPaymentId()),
	)

	items := make([]domain.Item, len(request.Items))
	for i, item := range request.Items {
		items[i] = s.itemToDomain(item)
	}

	err := s.app.CreateOrder(ctx, commands.CreateOrder{
		ID:         id,
		CustomerID: request.GetCustomerId(),
		PaymentID:  request.GetPaymentId(),
		Items:      items,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}
	
	return &orderingpb.CreateOrderResponse{Id: id}, err
}

func (s server) itemToDomain(item *orderingpb.Item) domain.Item {
	return domain.Item{
		ProductID:   item.GetProductId(),
		ProductName: item.GetProductId(),
		StoreID:     item.GetStoreId(),
		StoreName:   item.GetStoreName(),
		Price:       item.GetPrice(),
		Quantity:    int(item.GetQuantity()),
	}
}
