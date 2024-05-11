package domain

import "context"

type ProductCacheRepository interface {
	Add(ctx context.Context, ProductID, name, storeID string) error
	ProductRepository
}
