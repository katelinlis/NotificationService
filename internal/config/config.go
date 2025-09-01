package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type Config struct {
	Redis          redis.Options
	Broker         []string
	PushFCMKey     string
	UserServiceURL string
	Postgres       PostgresConfig
}

func Load() *Config {
	return &Config{
		Redis: redis.Options{
			Addr: "redis-master.default.svc.cluster.local:6379",
		},
		Broker:     []string{"100.64.0.2:9092"},
		PushFCMKey: os.Getenv("FCM_KEY"),
		Postgres: PostgresConfig{
			Host:     "100.64.0.2",
			Port:     5432,
			User:     "avenir",
			Password: os.Getenv("DATABASE_PASSWORD"),
			DBName:   "notification_service",
		},
		UserServiceURL: "http://avenirbackend-service.avenir.svc.cluster.local:3000",
	}
}
