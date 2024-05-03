package adapters

import (
	"context"
	"fmt"

	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
	"github.com/stackus/errors"
)

type PostgresMallRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MallRepository = (*PostgresMallRepository)(nil)

func NewPostgresMallRepository(tableName string, db postgres.DB) PostgresMallRepository {
	return PostgresMallRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PostgresMallRepository) AddStore(ctx context.Context, storeID, name, location string) error {
	const query = `INSERT INTO %s
					(id, name, location, participating)
					VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name, location, false)

	return err
}

func (r PostgresMallRepository) Find(ctx context.Context, storeID string) (*domain.MallStore, error) {
	const query = `SELECT name, location, participating
				  FROM %s WHERE id = $1 LIMIT 1`
	store := &domain.MallStore{
		ID: storeID,
	}
	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(
		&store.Name, &store.Location, &store.Participating,
	)
	if err != nil{
		return nil, errors.Wrap(err, "scanning store")
	}
	return store, nil
}

// func(r PostgresMallRepository){}
// func(r PostgresMallRepository){}
// func(r PostgresMallRepository){}
// func(r PostgresMallRepository){}
// func(r PostgresMallRepository){}
// func(r PostgresMallRepository){}
func (r PostgresMallRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
