package domain

const (
	ShppingListCreatedEvent = "depot.ShppingListCreated"
)

type ShppingListCreated struct {
	ShoppingList *ShoppingList
}

func (ShppingListCreated) Key() string { return ShppingListCreatedEvent }
