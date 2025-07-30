package main

import (
	"context"

	"github.com/avenir/notification-service/internal/config"
	"github.com/avenir/notification-service/internal/delivery/kafka"
	"github.com/avenir/notification-service/internal/delivery/push"
	"github.com/avenir/notification-service/internal/delivery/ws"
	"github.com/avenir/notification-service/internal/infra/redis"
	"github.com/avenir/notification-service/internal/processor"
)

func main() {
	// Initialize the notification service
	cfg := config.Load()
	//log := logger.New()

	redis := redis.InitRedis(cfg.Redis.Addr)
	ws := ws.InitWS(redis)
	push := push.InitPush()

	consumer := kafka.NewConsumer(cfg.Broker)

	processor := processor.New(redis, ws, push)

	consumer.Subscribe(context.Background(), "message.created", "notification-group", processor.HandleMessageCreated)

	select {}
}
