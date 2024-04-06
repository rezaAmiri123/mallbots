///go:build integration || database

package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
	"github.com/rezaAmiri123/mallbots/migrations"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type customerSuite struct {
	tableName string
	container testcontainers.Container
	db        *sql.DB
	repo      PostgresCustomerRepository
	suite.Suite
}

func TestPostgresCustomerRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode: skipping")
	}
	suite.Run(t, &customerSuite{tableName: constants.CustomersTableName})
}

func (s *customerSuite) SetupSuite() {
	var err error

	ctx := context.Background()
	initDir, err := filepath.Abs("./../../../docker/database")
	if err != nil {
		s.T().Fatal(err)
	}

	const dbUrl = "postgres://mallbots_user:mallbots_pass@localhost:%s/mallbots?sslmode=disable"
	s.container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			Hostname:     "postgres",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_PASSWORD": "itsasecret",
			},
			Mounts: []testcontainers.ContainerMount{
				testcontainers.BindMount(initDir, "/docker-entrypoint-initdb.d"),
			},
			WaitingFor: wait.ForSQL("5432/tcp", "pgx", func(host string, port nat.Port) string {
				return fmt.Sprintf(dbUrl, port.Port())
			}).WithStartupTimeout(5 * time.Second),
		},
		Started: true,
	})

	if err != nil {
		s.T().Fatal(err)
	}

	endpoint, err := s.container.Endpoint(ctx, "")
	if err != nil {
		s.T().Fatal(err)
	}

	s.db, err = sql.Open("pgx", fmt.Sprintf("postgres://mallbots_user:mallbots_pass@%s/mallbots?sslmode=disable", endpoint))
	if err != nil {
		s.T().Fatal(err)
	}

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		s.T().Fatal(err)
	}
	if err := goose.Up(s.db, "."); err != nil {
		s.T().Fatal(err)
	}
}

func (s *customerSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		s.T().Fatal(err)
	}
	if err := s.container.Terminate(context.Background()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *customerSuite) SetupTest() {
	s.repo = NewPostgresCustomerRepository(s.tableName, s.db)
}

func (s *customerSuite) TearDownTest() {
	_, err := s.db.ExecContext(context.Background(), fmt.Sprintf("TRUNCATE %s", s.tableName))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *customerSuite) TestPostgresCustomerRepository_Save() {
	customer := domain.NewCustomer("customer-id")
	customer.Name = "name"
	customer.SmsNumber = "sms-number"
	
	err := s.repo.Save(context.Background(), customer)
	s.NoError(err)

	query := fmt.Sprintf("SELECT name from %s WHERE id = $1 LIMIT 1", s.tableName)
	row:= s.db.QueryRow(query, "customer-id")
	s.NoError(row.Err())

	var name string
	s.NoError(row.Scan(&name))
	s.Equal("name", name)
}

func (s *customerSuite) TestPostgresCustomerRepository_Find() {
	query := `INSERT INTO %s (id, name, sms_number, enabled)
					VALUES ('customer-id', 'customer-name', 'customer-sms-number', true)`
				
	_, err := s.db.Exec(fmt.Sprintf(query, s.tableName))
	s.NoError(err)

	customer, err := s.repo.Find(context.Background(),"customer-id")
	s.NoError(err)
	s.Equal(customer.Name, "customer-name")
}

func (s *customerSuite) TestPostgresCustomerRepository_Update() {
	query := `INSERT INTO %s (id, name, sms_number, enabled)
					VALUES ('customer-id', 'customer-name', 'customer-sms-number', true)`
				
	_, err := s.db.Exec(fmt.Sprintf(query, s.tableName))
	s.NoError(err)

	customer := domain.NewCustomer("customer-id")
	customer.Name = "name-changed"
	err = s.repo.Update(context.Background(),customer)
	s.NoError(err)
	
	query = fmt.Sprintf("SELECT name from %s WHERE id = $1 LIMIT 1", s.tableName)
	row:= s.db.QueryRow(query, "customer-id")
	s.NoError(row.Err())

	var name string
	s.NoError(row.Scan(&name))
	s.Equal(name, customer.Name)

}

