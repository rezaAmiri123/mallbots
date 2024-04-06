package domain

const (
	StoreCreatedEvent = "stores.StoreCreated"
)
type StoreCreated struct {
	Name     string
	Location string
}

// Key implements registry.Registerable
func (StoreCreated) Key() string { return StoreCreatedEvent }
