package adapters

import (
	"context"

	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
	"github.com/stackus/errors"
	"google.golang.org/grpc/codes"
)

type GrpcProductRepository struct {
	client storespb.StoresServiceClient
}

var _ domain.ProductRepository = (*GrpcProductRepository)(nil)

func NewGrpcProductRepository(client storespb.StoresServiceClient) GrpcProductRepository {
	return GrpcProductRepository{
		client: client,
	}
}

func (r GrpcProductRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	resp, err := r.client.GetProduct(ctx, &storespb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		if errors.GRPCCode(err) == codes.NotFound {
			return nil, errors.ErrNotFound.Msg("product was not located")
		}
		return nil, errors.Wrap(err, "requesting product")
	}

	return r.productToDomain(resp.GetProduct()), nil
}

func (r GrpcProductRepository) productToDomain(store *storespb.Product) *domain.Product {
	return &domain.Product{
		ID:      store.GetId(),
		StoreID: store.GetStoreId(),
		Name:    store.GetName(),
		Price:   store.GetPrice(),
	}
}
