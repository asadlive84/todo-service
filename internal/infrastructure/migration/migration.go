package migration

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/mysql" 

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/go-sql-driver/mysql" 
)

func RunMigrations(dbURL string, migrationsPath string) error {
	if migrationsPath == "" {
		migrationsPath = "./migrations"
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("migration initialization error: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No pending migrations")
			return nil
		}
		return fmt.Errorf("migration error: %w", err)
	}

	fmt.Println("Migrations completed successfully")
	return nil
}
