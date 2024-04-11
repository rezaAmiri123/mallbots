package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/stores/internal/application"
	"github.com/rezaAmiri123/mallbots/stores/internal/application/commands"
	"github.com/rezaAmiri123/mallbots/stores/internal/application/queries"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	storespb.UnimplementedStoresServiceServer
}

var _ storespb.StoresServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	storespb.RegisterStoresServiceServer(registrar, server{
		app: app,
	})
	return nil
}

func (s server) CreateStore(ctx context.Context, request *storespb.CreateStoreRequest) (resp *storespb.CreateStoreResponse, err error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("StoreID", id),
	)

	err = s.app.CreateStore(ctx, commands.CreateStore{
		ID:       id,
		Name:     request.GetName(),
		Location: request.GetLocation(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.CreateStoreResponse{Id: id}, err
}

func (s server) GetStore(ctx context.Context, request *storespb.GetStoreRequest) (resp *storespb.GetStoreResponse, err error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("StoreID", request.GetId()),
	)
	store, err := s.app.GetStore(ctx, queries.GetStore{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &storespb.GetStoreResponse{Store: s.storeFromDomain(store)}, nil
}

func (s server) AddProduct(ctx context.Context, request *storespb.AddProductRequest) (resp *storespb.AddProductResponse, err error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("ProductID", id),
	)

	err = s.app.AddProduct(ctx, commands.AddProduct{
		ID:          id,
		StoreID:     request.GetStoreId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SKU:         request.GetSku(),
		Price:       request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.AddProductResponse{Id: id}, err
}

func (s server) IncreaseProductPrice(ctx context.Context, request *storespb.IncreaseProductPriceRequest) (resp *storespb.IncreaseProductPriceResponse, err error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ProductID", request.GetId()),
	)

	err = s.app.IncreaseProductPrice(ctx, commands.IncreaseProductPrice{
		ID:    request.GetId(),
		Price: request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.IncreaseProductPriceResponse{}, err
}

func (s server) storeFromDomain(store *domain.MallStore) *storespb.Store {
	return &storespb.Store{
		Id:            store.ID,
		Name:          store.Name,
		Location:      store.Location,
		Participating: store.Participating,
	}
}
