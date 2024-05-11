package application

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/depot/internal/application/commands"
	"github.com/rezaAmiri123/mallbots/depot/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateShoppingList(ctx context.Context, cmd commands.CreateShoppingList) error
	}
	Queries interface {
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateShoppingListHandler
	}
	appQueries struct {
	}
)

var _ App = (*Application)(nil)

func New(
	shoppingLists domain.ShoppingListRepository,
	stores domain.StoreRepository,
	products domain.ProductRepository,
	domainPublisher ddd.EventPublisher[ddd.AggregateEvent],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateShoppingListHandler: commands.NewCreateShoppingListHandler(shoppingLists, stores, products, domainPublisher),
		},
		appQueries: appQueries{},
	}
}
