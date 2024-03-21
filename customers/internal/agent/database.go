package agent

import (
	"database/sql"
	"fmt"
	
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
)

func (a *Agent) setupDatabase() error {
	// dbConn, err := postgres.NewDB(postgres.Config{
	// 	PGDriver:     a.PGDriver,
	// 	PGHost:       a.PGHost,
	// 	PGPort:       a.PGPort,
	// 	PGUser:       a.PGUser,
	// 	PGDBName:     a.PGDBName,
	// 	PGPassword:   a.PGPassword,
	// 	PGSearchPath: a.PGSearchPath,
	// })
	// var pgConn *pgxpool.Pool
	// pgConn, err := pgxpool.Connect(context.Background(), a.config.Postgres.Conn)
	db, err := sql.Open("pgx", a.config.Postgres.Conn)
	if err != nil {
		return fmt.Errorf("cannot load db: %w", err)
	}
	if err = postgres.MigrateUp(db, migrations.FS); err != nil {
		return err
	}

	a.container.AddSingleton(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return dbConn, nil
	})

	a.container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return dbConn.Begin()
	})

	return nil
}
