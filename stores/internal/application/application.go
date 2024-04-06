package application

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/stores/internal/application/commands"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateStore(ctx context.Context, cmd commands.CreateStore) error
	}
	Queries interface {
	}
	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateStoreHandler
	}
	appQueries struct{}
)

var _ App = (*Application)(nil)

func New(
	stores domain.StoreRepository,
	// products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateStoreHandler: commands.NewCreateStoreHandler(stores, publisher),
		},
		appQueries: appQueries{},
	}
}
