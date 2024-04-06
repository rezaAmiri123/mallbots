package domain

const(
	ProductAddedEvent          = "stores.ProductAdded"
)

type ProductAdded struct {
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}
func(ProductAdded)Key()string{return ProductAddedEvent}