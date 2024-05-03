package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
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

func (r PostgresProductCacheRepository) Add(ctx context.Context, productID, storeID, name string, price float64) error {
	const query = `INSERT INTO %s (id, store_id, name, price) VALUES ($1,$2,$3,$4)
				  ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, r.table(query), productID, storeID, name, price)

	return err
}

func (r PostgresProductCacheRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	const query = `SELECT store_id, name, price from %s
				  WHERE id = $1 LIMIT 1`
	product := &domain.Product{
		ID: productID,
	}

	fmt.Println("err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(")
	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(
		&product.StoreID, &product.Name, &product.Price,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning product")
		}
		fmt.Println("sql.ErrNoRows.Error(): ", sql.ErrNoRows.Error())
		fmt.Println("err.Error(): ", err.Error())
		// if strings.Contains(sql.ErrNoRows.Error(), err.Error()){
		// 	return nil, errors.Wrap(err, "scanning product")
		// }

		fmt.Println("product, err = r.fallback.Find(ctx, productID)")
		product, err = r.fallback.Find(ctx, productID)
		if err != nil {
			return nil, errors.Wrap(err, "product fallback failed")
		}
		// attempt to add it to the cache
		return product, r.Add(ctx, product.ID, product.StoreID,product.Name,product.Price)
	}

	return product, nil
}

func (r PostgresProductCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
