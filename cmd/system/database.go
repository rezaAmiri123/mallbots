package system

import (
	"database/sql"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)
type DB interface{
	DB() *sql.DB
}

func (s *System) initDB() (err error) {
	s.db, err = sql.Open("pgx", s.cfg.Postgres.Conn)
	return err
}

func (s *System) MigrateDB(fs fs.FS) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(s.db, "."); err != nil {
		return err
	}
	return nil
}

func (s *System) DB() *sql.DB {
	return s.db
}

// func (s *System) setupDatabase() error {
// 	db, err := sql.Open("pgx", a.config.Postgres.Conn)
// 	if err != nil {
// 		return fmt.Errorf("cannot load db: %w", err)
// 	}
// 	if err := db.Ping(); err != nil {
// 		return fmt.Errorf("cannot ping db: %w", err)
// 	}

// 	db.SetMaxIdleConns(10)

// 	fmt.Println("connction pool size: ", db.Stats().Idle)
// 	if err = postgres.MigrateUp(db, migrations.FS); err != nil {
// 		return err
// 	}

// 	a.container.AddSingleton(constants.DatabaseKey, func(c di.Container) (any, error) {
// 		return db, nil
// 	})

// 	a.container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
// 		return db.Begin()
// 	})

// 	return nil
// }

// func (a *Agent) cleanupDatabase() error {
// 	db := a.container.Get(constants.DatabaseKey).(*sql.DB)
// 	logger := edatlog.DefaultLogger
// 	if err := db.Close(); err != nil {
// 		logger.Error("ran into an issue shutting down the database connection", edatlog.Error(err))
// 	}
// 	logger.Info("clean up database")
// 	return nil
// }
