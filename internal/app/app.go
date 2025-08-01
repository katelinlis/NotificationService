package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/avenir/notification-service/internal/config"
	"github.com/avenir/notification-service/internal/delivery/kafka"
	"github.com/avenir/notification-service/internal/delivery/push"
	"github.com/avenir/notification-service/internal/delivery/ws"
	"github.com/avenir/notification-service/internal/domain/repository"
	"github.com/avenir/notification-service/internal/infra/httpclient"
	"github.com/avenir/notification-service/internal/infra/postgres"
	"github.com/avenir/notification-service/internal/infra/redis"
	"github.com/avenir/notification-service/internal/processor"
	"github.com/avenir/notification-service/utils"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func StartApp(ctx context.Context) {
	// Initialize the notification service
	cfg := config.Load()
	//log := logger.New()

	redis := redis.InitRedis(cfg.Redis.Addr)
	ws := ws.InitWS(redis)

	push := push.InitPush()

	consumer := kafka.NewConsumer(cfg.Broker)
	postgresDB := postgres.InitPostgres(cfg.Postgres)

	processor := processor.New(
		redis,
		ws,
		push,
		postgresDB,
		httpclient.NewUserClient(cfg.UserServiceURL, http.DefaultClient),
	)

	go consumer.Subscribe(context.Background(), "message.created", "notification-group", processor.HandleMessageCreated)
	go ws.Read()
	go ws.EmitMessages()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	go http.ListenAndServe("localhost:3000", nil)
	go HTTPAPI(postgresDB)

	<-ctx.Done()
}

func HTTPAPI(db repository.Store) {
	mux := http.NewServeMux()

	// Группа /api/internal/
	internalMux := http.NewServeMux()
	internalMux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {})
	internalMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Metrics endpoint"))
	})

	// Вешаем подпрефикс
	mux.Handle("/api/internal/", http.StripPrefix("/api/internal", internalMux))

	// Группа /api/public/
	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	publicMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.AuthCheck(*r, w)
		if err != nil {
			return
		}

		notif, err := db.Notification().FindByUserID(context.Background(), int(claims.ClientID))

		w.WriteHeader(http.StatusOK)
		bytes, err := json.Marshal(notif)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		w.Write(bytes)
	})

	mux.Handle("/api/v1/notifications", http.StripPrefix("/api/v1/notifications/", publicMux))

	http.ListenAndServe(":8080", mux)
}
