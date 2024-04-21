package adapters

import (
	"context"

	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
)

type GrpcStoreRepository struct {
	client storespb.StoresServiceClient
}

var _ domain.StoreRepository = (*GrpcStoreRepository)(nil)

func NewGrpcStoreRepository(client storespb.StoresServiceClient) GrpcStoreRepository {
	return GrpcStoreRepository{
		client: client,
	}
}

func (r GrpcStoreRepository) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	resp, err := r.client.GetStore(ctx, &storespb.GetStoreRequest{
		Id: storeID,
	})
	if err != nil {
		return nil, err
	}

	return r.storeToDomain(resp.Store), nil
}

func (r GrpcStoreRepository) storeToDomain(store *storespb.Store) *domain.Store {
	return &domain.Store{
		ID:   store.GetId(),
		Name: store.GetName(),
	}
}
