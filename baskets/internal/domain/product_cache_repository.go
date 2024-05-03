package domain

import "context"

type ProductCacheRepository interface {
	Add(ctx context.Context, productID, storeID, name string, price float64) error
	ProductRepository
}
