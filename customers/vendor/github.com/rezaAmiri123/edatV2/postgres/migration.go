package postgres

import (
	"database/sql"
	"io/fs"

	"github.com/pressly/goose/v3"
)

func MigrateUp(db *sql.DB,fs fs.FS) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	
	if err := goose.Up(db, "."); err != nil {
		return err
	}
	
	return nil
}

func MigrateDown(db *sql.DB,fs fs.FS) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Down(db, "."); err != nil {
		return err
	}
	return nil
}
