package application

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
)

type (
	StartBasket struct {
		ID         string
		CustomerID string
	}

	App interface {
		StartBasket(ctx context.Context, start StartBasket) error
	}

	Application struct {
		baskets   domain.BasketRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

var _ App = (*Application)(nil)

func New(
	baskets domain.BasketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		baskets:   baskets,
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
