package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/depot/depotpb"
	"github.com/rezaAmiri123/mallbots/depot/internal/application"
	"github.com/rezaAmiri123/mallbots/depot/internal/application/commands"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	depotpb.UnimplementedDepotServiceServer
}

var _ depotpb.DepotServiceServer = (*server)(nil)

func Register(app application.App, register grpc.ServiceRegistrar) error {
	depotpb.RegisterDepotServiceServer(register, server{
		app: app,
	})
	return nil
}

func (s server) CreateShoppingList(ctx context.Context, request *depotpb.CreateShoppingListRequest) (*depotpb.CreateShoppingListResponse, error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("ShoppingListID", id),
		attribute.String("OrderID", request.GetOrderId()),
	)

	items := make([]commands.OrderItem, 0, len(request.GetItems()))
	for _, item := range request.GetItems() {
		items = append(items, s.itemToDomain(item))
	}

	err := s.app.CreateShoppingList(ctx, commands.CreateShoppingList{
		ID:      id,
		OrderID: request.GetOrderId(),
		Items:   items,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &depotpb.CreateShoppingListResponse{Id: id}, err
}

func (s server) itemToDomain(item *depotpb.OrderItem) commands.OrderItem {
	return commands.OrderItem{
		StoreID:   item.GetStoreId(),
		ProductID: item.GetProductId(),
		Quantity:  int(item.GetQuantity()),
	}
}
