package adapters

import (
	"context"
	"fmt"

	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
)

type PostgresCatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogRepository = (*PostgresCatalogRepository)(nil)

func NewPostgresCatalogRepository(tableName string, db postgres.DB) PostgresCatalogRepository {
	return PostgresCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PostgresCatalogRepository) AddProduct(ctx context.Context, productID, storeID, name, description, sku string, price float64) error {
	const query = `INSERT INTO %s
					(id, store_id,name, description, sku, price)
					VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, r.table(query), productID, storeID, name, description, sku, price)

	return err
}

func (r PostgresCatalogRepository) Find(ctx context.Context, productID string) (*domain.CatalogProduct, error) {
	const query = `SELECT store_id,name, description, sku, price
				  FROM %s WHERE id = $1 LIMIT 1`
	catalog := &domain.CatalogProduct{
		ID: productID,
	}
	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(
		&catalog.StoreID, &catalog.Name, &catalog.Description, &catalog.SKU, &catalog.Price,
	)
	return catalog, err
}

func (r PostgresCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
