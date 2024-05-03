package domain

const (
	OrderCreatedEvent = "ordering.OrderCreated"
)

type OrderCreated struct {
	CustomerID string
	PaymentID  string
	ShoppingID string
	Items      []Item
}

func (OrderCreated) Key() string { return OrderCreatedEvent }
