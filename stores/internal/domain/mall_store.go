package domain

import "context"

type MallStore struct {
	ID            string
	Name          string
	Location      string
	Participating bool
}

type MallRepository interface {
	AddStore(ctx context.Context, storeID, name, location string) error
	Find(ctx context.Context, storeID string)(*MallStore,error)
}
