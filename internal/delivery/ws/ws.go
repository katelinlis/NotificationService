package ws

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/katelinlis/NotificationService/internal/domain/model"
	"github.com/katelinlis/NotificationService/internal/infra/redis"
	"github.com/katelinlis/NotificationService/utils"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

type WebSocketSender struct {
	redis     redis.RedisService
	bufferMap chan model.MessageCreatedEvent
	clientMap sync.Map // clientMap stores the client ID and its corresponding socket connection
}

func (w *WebSocketSender) Send(instanceID string, msg model.MessageCreatedEvent) error {
	return w.redis.SendPublish(instanceID, msg)
}

func (w *WebSocketSender) Read() {
	instanceID := os.Getenv("HOSTNAME")
	if instanceID == "" {
		instanceID = "default-instance"
	}

	sub := w.redis.Subscribe(context.Background(), instanceID)

	//msg, err := sub.ReceiveMessage(context.Background())

	for {
		msg, err := sub.ReceiveMessage(context.Background())
		if err != nil {
			continue
		}

		var events model.MessageCreatedEvent
		if err := json.Unmarshal([]byte(msg.Payload), &events); err != nil {
			continue
		}

		if events.ReceiverID == 0 {
			continue
		}

		w.bufferMap <- events

	}

}

func (w *WebSocketSender) EmitMessages() {
	for msg := range w.bufferMap {
		client, ok := w.clientMap.Load(msg.ReceiverID)

		if ok {
			client.(*socket.Socket).Emit("message", msg)
		}
	}
}

func InitWS(redis redis.RedisService) *WebSocketSender {

	instanceID := os.Getenv("HOSTNAME")
	if instanceID == "" {
		instanceID = "default-instance"
	}
	httpServer := types.NewWebServer(nil)
	WSS := &WebSocketSender{
		redis:     redis,
		bufferMap: make(chan model.MessageCreatedEvent, 1000),
		clientMap: sync.Map{},
	}

	s := socket.DefaultServerOptions()
	s.SetServeClient(true)
	s.SetPingInterval(20 * time.Second) // Интервал отправки Ping
	s.SetPingTimeout(60 * time.Second)  // Тайм-аут до разрыва соединения
	s.SetCors(&types.Cors{              // Кросс-доменные заголовки
		Origin:      "*",
		Credentials: true,
	})

	// Создаем сервер Socket.IO
	io := socket.NewServer(httpServer, s)

	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		println("connect")

		done := make(chan struct{})
		clientID := 0

		client.On("auth", func(data ...any) {
			// Обработка авторизации клиента
			if len(data) > 0 {

				token, _ := data[0].(string)
				//println(token)
				jwt, err := utils.JWTParse(token)
				println(jwt)
				if err != nil {
					println(err.Error())
					done <- struct{}{}
					return
				}

				clientID = int(jwt.ClientID)
				println(clientID)
				WSS.clientMap.Store(clientID, client)
				WSS.redis.SetInstanceID(clientID, instanceID)
			}
		})
		client.On("event", func(datas ...any) {
		})
		client.On("disconnect", func(...any) {
			WSS.clientMap.Delete(clientID)
			WSS.redis.DelInstanceID(clientID)

		})

	})

	println("rest")
	httpServer.Listen("0.0.0.0:4000", nil)

	return WSS
}
