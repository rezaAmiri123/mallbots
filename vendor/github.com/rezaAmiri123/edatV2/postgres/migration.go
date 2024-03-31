package postgres

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"io/fs"
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
