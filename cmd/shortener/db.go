package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func createConnect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to DB: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("cannot get driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to get migrate instance: %w", err)
	}

	return m.Up()
}
