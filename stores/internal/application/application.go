package application

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/stores/internal/application/commands"
	"github.com/rezaAmiri123/mallbots/stores/internal/application/queries"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CreateStore(ctx context.Context, cmd commands.CreateStore) error
		AddProduct(ctx context.Context, cmd commands.AddProduct) error
		IncreaseProductPrice(ctx context.Context, cmd commands.IncreaseProductPrice) error
	}

	Queries interface {
		GetStore(ctx context.Context, query queries.GetStore) (*domain.MallStore, error)
		GetProduct(ctx context.Context, query queries.GetProduct) (*domain.CatalogProduct, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateStoreHandler
		commands.AddProductHandler
		commands.IncreaseProductPriceHandler
	}
	appQueries struct {
		queries.GetStoreHandler
		queries.GetProductHandler
	}
)

var _ App = (*Application)(nil)

func New(
	stores domain.StoreRepository,
	mall domain.MallRepository,
	catalog domain.CatalogRepository,
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateStoreHandler:          commands.NewCreateStoreHandler(stores, publisher),
			AddProductHandler:           commands.NewAddProductHandler(products, publisher),
			IncreaseProductPriceHandler: commands.NewIncreaseProductPriceHandler(products, publisher),
		},
		appQueries: appQueries{
			GetStoreHandler:   queries.NewGetStoreHandler(mall),
			GetProductHandler: queries.NewGetProductHandler(catalog),
		},
	}
}
