package domain

const(
	ProductAddedEvent          = "stores.ProductAdded"
	ProductPriceIncreasedEvent = "stores.ProductPriceIncreased"
)

type ProductAdded struct {
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}
func(ProductAdded)Key()string{return ProductAddedEvent}

type ProductPriceChanged struct{
	Delta float64
}

type ProductPriceDelta struct{
	Product *Product
	Delta float64
}
