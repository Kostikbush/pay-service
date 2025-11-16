package postgres

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	migrations "pay-service/db/migrations"
)

func RunMigrations(dbURL string) error {
	src, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(openStdDB(dbURL), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
