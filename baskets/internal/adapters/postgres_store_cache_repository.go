package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
	"github.com/stackus/errors"
)

type PostgresStoreCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.StoreRepository
}

var _ domain.StoreCacheRepository = (*PostgresStoreCacheRepository)(nil)

func NewPostgresStoreCacheRepository(
	tableName string,
	db postgres.DB,
	fallback domain.StoreRepository,
) PostgresStoreCacheRepository {
	return PostgresStoreCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r PostgresStoreCacheRepository) Add(ctx context.Context, storeID, name string) error {
	const query = `INSERT INTO %s (id, name) VALUES ($1,$2)
				  ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)

	return err
}

func (r PostgresStoreCacheRepository) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	const query = `SELECT name from %s
				  WHERE id = $1 LIMIT 1`
	store := &domain.Store{
		ID: storeID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning store")
		}
		store, err = r.fallback.Find(ctx, storeID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.Wrap(err, "store fallback failed")
		}
		// attempt to add it to the cache
		return store, r.Add(ctx, store.ID, store.Name)
	}
	
	return store, nil
}

func (r PostgresStoreCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
