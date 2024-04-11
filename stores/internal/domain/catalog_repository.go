package domain

import "context"

type CatalogProduct struct {
	ID          string
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

type CatalogRepository interface {
	AddProduct(ctx context.Context, productID, storeID, name, description, sku string, price float64) error
	Find(ctx context.Context, productID string) (*CatalogProduct, error)
}
