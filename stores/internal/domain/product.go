package domain

import "github.com/rezaAmiri123/edatV2/es"

const ProductAggregate = "stores.Product"

type Product struct {
	es.Aggregate
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

var _ interface{
	es.EventApplier
	es.Snapshotter
} = (*Product)(nil)

