package application

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
	"github.com/stackus/errors"
)

type (
	StartBasket struct {
		ID         string
		CustomerID string
	}

	CheckoutBasket struct {
		ID        string
		PaymentID string
	}

	GetBasket struct {
		ID string
	}

	AddItem struct {
		ID        string
		ProductID string
		Quantity  int
	}

	App interface {
		StartBasket(ctx context.Context, start StartBasket) error
		CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error
		GetBasket(ctx context.Context, getbasket GetBasket) (*domain.Basket, error)
		AddItem(ctx context.Context, add AddItem) error
	}

	Application struct {
		baskets   domain.BasketRepository
		stores    domain.StoreRepository
		products  domain.ProductRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

var _ App = (*Application)(nil)

func New(
	baskets domain.BasketRepository,
	stores domain.StoreRepository,
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		baskets:   baskets,
		stores:    stores,
		products:  products,
		publisher: publisher,
	}
}

func (a Application) StartBasket(ctx context.Context, start StartBasket) error {
	basket, err := a.baskets.Load(ctx, start.ID)
	if err != nil {
		return err
	}

	event, err := basket.Start(start.CustomerID)
	if err != nil {
		return err
	}

	err = a.baskets.Save(ctx, basket)
	if err != nil {
		return err
	}

	return a.publisher.Publish(ctx, event)
}

func (a Application) GetBasket(ctx context.Context, getbasket GetBasket) (*domain.Basket, error) {
	return a.baskets.Load(ctx, getbasket.ID)
}

func (a Application) AddItem(ctx context.Context, add AddItem) error {
	basket, err := a.baskets.Load(ctx, add.ID)
	if err != nil {
		return err
	}

	product, err := a.products.Find(ctx, add.ProductID)
	if err != nil {
		return err
	}

	store, err := a.stores.Find(ctx, product.StoreID)
	if err != nil {
		return err
	}

	err = basket.AddItem(store, product, add.Quantity)
	if err != nil {
		return err
	}

	return a.baskets.Save(ctx, basket)
}

func (a Application) CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error {
	basket, err := a.baskets.Load(ctx, checkout.ID)
	if err != nil {
		return err
	}

	event, err := basket.Checkout(checkout.PaymentID)
	if err != nil {
		return errors.Wrap(err, "baskets checkout")
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		return errors.Wrap(err, "baskets checkout")
	}

	return a.publisher.Publish(ctx, event)
}
