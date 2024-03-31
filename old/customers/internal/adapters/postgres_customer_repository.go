package adapters

import (
	"context"
	"fmt"

	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
)

type PostgresCustomerRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CustomerRepository = (*PostgresCustomerRepository)(nil)

func NewPostgresCustomerRepository(tableName string, db postgres.DB) PostgresCustomerRepository {
	return PostgresCustomerRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PostgresCustomerRepository) Save(ctx context.Context, customer *domain.Customer) error {
	const query = `INSERT INTO %s
	                (id, name, sms_number, enabled)
					VALUES ($1,$2,$3,$4)`

	_, err := r.db.ExecContext(ctx, r.table(query),
		customer.ID(),
		customer.Name,
		customer.SmsNumber,
		customer.Enabled,
	)

	return err
}

func (r PostgresCustomerRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	const query = `SELECT name, sms_number, enabled
				   from %s WHERE id = $1 LIMIT 1`

	customer := domain.NewCustomer(customerID)

	err := r.db.QueryRowContext(ctx, r.table(query), customerID).Scan(
		&customer.Name,
		&customer.SmsNumber,
		&customer.Enabled,
	)

	return customer, err
}

func (r PostgresCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	const query = `UPDATE %s SET
	               name = $2, sms_number = $3, enabled = $4
				   WHERE id = $1`
	_, err := r.db.ExecContext(ctx, r.table(query),
		customer.ID(),
		customer.Name,
		customer.SmsNumber,
		customer.Enabled,
	)

	return err
}

func (r PostgresCustomerRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
