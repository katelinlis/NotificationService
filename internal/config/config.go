package config

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis      redis.Options
	Broker     []string
	PushFCMKey string
}

func Load() *Config {
	return &Config{
		Redis: redis.Options{
			Addr: "localhost:6379",
		},
		Broker:     []string{"localhost:9092"},
		PushFCMKey: "your-fcm-key",
	}
}
