package app

import (
	"context"
	"net/http"

	"github.com/avenir/notification-service/internal/config"
	"github.com/avenir/notification-service/internal/delivery/kafka"
	"github.com/avenir/notification-service/internal/delivery/push"
	"github.com/avenir/notification-service/internal/delivery/ws"
	"github.com/avenir/notification-service/internal/infra/httpclient"
	"github.com/avenir/notification-service/internal/infra/postgres"
	"github.com/avenir/notification-service/internal/infra/redis"
	"github.com/avenir/notification-service/internal/processor"
)

func StartApp(ctx context.Context) {
	// Initialize the notification service
	cfg := config.Load()
	//log := logger.New()

	redis := redis.InitRedis(cfg.Redis.Addr)
	ws := ws.InitWS(redis)

	push := push.InitPush()

	consumer := kafka.NewConsumer(cfg.Broker)

	processor := processor.New(
		redis,
		ws,
		push,
		postgres.InitPostgres(cfg.Postgres),
		httpclient.NewUserClient(cfg.UserServiceURL, http.DefaultClient),
	)

	go consumer.Subscribe(context.Background(), "message.created", "notification-group", processor.HandleMessageCreated)
	go ws.Read()
	go ws.EmitMessages()

	<-ctx.Done()
}
