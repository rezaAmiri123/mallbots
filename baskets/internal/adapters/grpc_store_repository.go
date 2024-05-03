package adapters

import (
	"context"
	"fmt"

	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
	"github.com/stackus/errors"
	"google.golang.org/grpc/codes"
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
		fmt.Println("************************")
		fmt.Println(err, storeID)
		fmt.Println(errors.GRPCCode(err))
		if errors.GRPCCode(err) == codes.NotFound {
			return nil, errors.ErrNotFound.Msg("store was not located")
		}
		return nil, errors.Wrap(err, "requesting store")

	}

	return r.storeToDomain(resp.Store), nil
}

func (r GrpcStoreRepository) storeToDomain(store *storespb.Store) *domain.Store {
	return &domain.Store{
		ID:   store.GetId(),
		Name: store.GetName(),
	}
}
