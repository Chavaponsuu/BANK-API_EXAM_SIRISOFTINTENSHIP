package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"

	"github.com/krizad/go-gin-api/config"
	"github.com/krizad/go-gin-api/migrations"
)

var DB *sql.DB

func Connect(cfg *config.Config) {
	var err error
	DB, err = sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("connected to PostgreSQL")
}

func RunMigrations() error {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Println("migrations: no migrations applied")
	} else {
		fmt.Printf("migrations: at version %d (dirty=%v)\n", version, dirty)
	}

	return nil
}

func RollbackMigrations() error {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration down failed: %w", err)
	}

	fmt.Println("migrations: rolled back all")
	return nil
}

func ForceMigration(version int) error {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Force(version); err != nil {
		return fmt.Errorf("migration force failed: %w", err)
	}

	fmt.Printf("migrations: forced to version %d\n", version)
	return nil
}

func Close() {
	if DB != nil {
		_ = DB.Close()
	}
}
