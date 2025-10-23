package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs SQL migrations from the given folder against the DATABASE_URL.
// folder should be a relative path like "./migrations" resolved from the working directory.
func RunMigrations(databaseURL, folder string) {
	m, err := migrate.New("file://"+folder, databaseURL)
	if err != nil {
		log.Printf("migrate: init error: %v", err)
		return
	}
	defer func() { _, _ = m.Close() }()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("migrate: up error: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Printf("migrate: no change")
	} else {
		log.Printf("migrate: applied")
	}
}
