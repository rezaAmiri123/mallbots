package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/depot/internal/domain"
	"github.com/stackus/errors"
)

type PostgresShoppingListRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.ShoppingListRepository = (*PostgresShoppingListRepository)(nil)

func NewPostgresShoppingListRepository(tableName string, db postgres.DB) PostgresShoppingListRepository {
	return PostgresShoppingListRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PostgresShoppingListRepository) Find(ctx context.Context, shoppingListID string) (*domain.ShoppingList, error) {
	const query = `SELECT order_id, stops,assigned_bot_id, status
	               FROM %s WHERE id = $1 LIMIT 1`
	shoppingList := domain.NewShoppingList(shoppingListID)

	var stops []byte
	var status string

	err := r.db.QueryRowContext(ctx, r.table(query), shoppingListID).Scan(
		&shoppingList.OrderID,
		&stops,
		&shoppingList.AssignedBotID,
		&status,
	)
	if err != nil {
		return nil, errors.ErrInternalServerError.Err(err)
	}

	shoppingList.Status = domain.ToShoppingListStatus(status)

	err = json.Unmarshal(stops, &shoppingList.Stops)
	if err != nil {
		return nil, errors.ErrInternalServerError.Err(err)
	}

	return shoppingList, nil
}

func (r PostgresShoppingListRepository) Save(ctx context.Context, list *domain.ShoppingList) error {
	const query = `INSERT INTO %s (id, order_id, stops, assigned_bot_id, status)
	               VALUES ($1, $2, $3, $4, $5)`
	stops, err := json.Marshal(list.Stops)
	if err != nil {
		return errors.ErrInternalServerError.Err(err)
	}

	_, err = r.db.ExecContext(ctx, r.table(query), list.ID(), list.OrderID, stops, list.AssignedBotID, list.Status.String())
	if err != nil {
		return errors.ErrInternalServerError.Err(err)
	}

	return nil
}

func (r PostgresShoppingListRepository) Update(ctx context.Context, list *domain.ShoppingList) error {
	const query = `UPDATE %s SET stops = $2, assigned_bot_id = $3, status = $4
				   WHERE id = $1`
	stops, err := json.Marshal(list.Stops)
	if err != nil {
		return errors.ErrInternalServerError.Err(err)
	}

	_, err = r.db.ExecContext(ctx, r.table(query), list.ID(), stops, list.AssignedBotID, list.Status.String())
	if err != nil {
		return errors.ErrInternalServerError.Err(err)
	}

	return nil
}

func (r PostgresShoppingListRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
