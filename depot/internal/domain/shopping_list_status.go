package domain

type ShoppingListStatus string

const (
	ShoppingLisUnknown      ShoppingListStatus = ""
	ShoppingListIsPending   ShoppingListStatus = "pending"
	ShoppingListIsAvailable ShoppingListStatus = "available"
	ShoppingListIsAssigned  ShoppingListStatus = "assigned"
	ShoppingListIsActive    ShoppingListStatus = "active"
	ShoppingListIsCompleted ShoppingListStatus = "completed"
	ShoppingListIsCanceled  ShoppingListStatus = "canceled"
)

func (s ShoppingListStatus) String() string {
	switch s {
	case ShoppingListIsAvailable, ShoppingListIsAssigned, ShoppingListIsActive,
		ShoppingListIsCompleted, ShoppingListIsCanceled:
		return string(s)
	default:
		return ""
	}
}

func ToShoppingListStatus(status string) ShoppingListStatus {
	switch status {
	case ShoppingListIsPending.String():
		return ShoppingListIsPending
	case ShoppingListIsAvailable.String():
		return ShoppingListIsAvailable
	case ShoppingListIsAssigned.String():
		return ShoppingListIsAssigned
	case ShoppingListIsActive.String():
		return ShoppingListIsActive
	case ShoppingListIsCompleted.String():
		return ShoppingListIsCompleted
	case ShoppingListIsCanceled.String():
		return ShoppingListIsCanceled
	default:
		return ShoppingLisUnknown
	}
}
