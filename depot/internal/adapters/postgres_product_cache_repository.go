package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/depot/internal/domain"
	"github.com/stackus/errors"
)

type PostgresProductCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.ProductRepository
}

var _ domain.ProductCacheRepository = (*PostgresProductCacheRepository)(nil)

func NewPostgresProductCacheRepository(
	tableName string,
	db postgres.DB,
	fallback domain.ProductRepository,
) PostgresProductCacheRepository {
	return PostgresProductCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r PostgresProductCacheRepository) Add(ctx context.Context, productID, name, storeID string) error {
	const query = `INSERT INTO %s (id, name, store_id)
				  VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, r.table(query), productID, name, storeID)

	return err
}

func (r PostgresProductCacheRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	const query = `SELECT name, store_id FROM %s
				   WHERE id = $1  LIMIT 1`
	product := &domain.Product{
		ID: productID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(
		&product.Name, &product.StoreID,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning product")
		}
		product, err = r.fallback.Find(ctx, productID)
		if err != nil {
			return nil, errors.Wrap(err, "product fallback failed")
		}
		// attempt to add it to the cache
		return product, r.Add(ctx, product.ID, product.Name, product.StoreID)
	}

	return product, nil
}

func (r PostgresProductCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
