package application

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/ordering/internal/application/commands"
	"github.com/rezaAmiri123/mallbots/ordering/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateOrder(ctx context.Context, cmd commands.CreateOrder) error
	}
	Queries interface {
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateOrderHandler
	}
	appQueries struct {
	}
)

var _ App = (*Application)(nil)

func New(
	orders domain.OrderRepository,
	puplisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateOrderHandler: commands.NewCreateOrderHandler(orders, puplisher),
		},
		appQueries: appQueries{},
	}
}
