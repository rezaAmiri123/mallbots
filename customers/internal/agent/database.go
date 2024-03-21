package agent

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rezaAmiri123/edatV2/di"
	edatlog "github.com/rezaAmiri123/edatV2/log"
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/customers/internal/adapters/migrations"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
)

func (a *Agent) setupDatabase() error {
	db, err := sql.Open("pgx", a.config.Postgres.Conn)
	if err != nil {
		return fmt.Errorf("cannot load db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("cannot ping db: %w", err)
	}
	if err = postgres.MigrateUp(db, migrations.FS); err != nil {
		return err
	}

	a.container.AddSingleton(constants.DatabaseKey, func(c di.Container) (any, error) {
		return db, nil
	})

	a.container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return db.Begin()
	})

	return nil
}

func (a *Agent) cleanupDatabase() error {
	db := a.container.Get(constants.DatabaseKey).(*sql.DB)
	logger := edatlog.DefaultLogger
	if err := db.Close(); err != nil {
		logger.Error("ran into an issue shutting down the database connection", edatlog.Error(err))
	}
	logger.Info("clean up database")
	return nil
}
