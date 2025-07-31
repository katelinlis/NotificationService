package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/avenir/notification-service/internal/config"
	"github.com/avenir/notification-service/internal/domain/repository"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type store struct {
	notification repository.NotificationRepository
}

func newStore(db *sql.DB) repository.Store {
	return &store{
		notification: NewNotificationRepository(db),
	}
}

func (s *store) Notification() repository.NotificationRepository {
	return s.notification
}

func InitPostgres(cfg config.PostgresConfig) repository.Store {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("driver error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://scripts/migrations", // путь должен быть валидный внутри Docker
		"postgres", driver)
	if err != nil {
		log.Fatalf("migration init error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	// Initialize the store

	return newStore(db)
}
