package domain

import (
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/registry"
)

const ShoppingListAggregate = "depot.ShoppingList"

type ShoppingList struct {
	ddd.Aggregate
	OrderID       string
	Stops         Stops
	AssignedBotID string
	Status        ShoppingListStatus
}

var _ interface {
	registry.Registrable
} = (*ShoppingList)(nil)

func NewShoppingList(id string) *ShoppingList {
	return &ShoppingList{
		Aggregate: ddd.NewAggregate(id, ShoppingListAggregate),
	}
}

func (ShoppingList) Key() string { return ShoppingListAggregate }

func CreateShoppingList(id, orderID string) *ShoppingList {
	shoppingList := NewShoppingList(id)
	shoppingList.OrderID = orderID
	shoppingList.Status = ShoppingListIsPending
	shoppingList.Stops = make(Stops)

	shoppingList.AddEvent(ShppingListCreatedEvent, &ShppingListCreated{
		ShoppingList: shoppingList,
	})

	return shoppingList
}

func (sl *ShoppingList) AddItem(store *Store, product *Product, quantity int) error {
	if _, exists := sl.Stops[store.ID]; !exists {
		sl.Stops[store.ID] = &Stop{
			StoreName:     store.Name,
			StoreLocation: store.Location,
			Items:         make(Items),
		}
	}

	return sl.Stops[store.ID].AddItem(product, quantity)
}
