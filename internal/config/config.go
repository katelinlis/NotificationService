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
			Addr: "redis-master.default.svc.cluster.local:6379",
		},
		Broker:     []string{"10.8.1.1:9092"},
		PushFCMKey: "your-fcm-key",
	}
}
