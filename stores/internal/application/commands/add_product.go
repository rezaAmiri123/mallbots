package commands

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
	"github.com/stackus/errors"
)

type (
	AddProduct struct {
		ID          string
		StoreID     string
		Name        string
		Description string
		SKU         string
		Price       float64
	}

	AddProductHandler struct {
		products    domain.ProductRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewAddProductHandler(
	products    domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AddProductHandler {
	return AddProductHandler{
		products: products,
		publisher: publisher,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	event, err := product.InitProduct(cmd.ID, cmd.StoreID,cmd.Name,cmd.Description,cmd.SKU, cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error initializing")
	}

	err = h.products.Save(ctx, product)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
