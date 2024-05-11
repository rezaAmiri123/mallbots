package domain

import "context"

type StoreCacheRepository interface{
	Add(ctx context.Context, storeID, name,location string)error
	StoreRepository
}
