package processor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/avenir/notification-service/internal/delivery/push"
	"github.com/avenir/notification-service/internal/delivery/ws"
	"github.com/avenir/notification-service/internal/domain/model"
	"github.com/avenir/notification-service/internal/domain/repository"
	"github.com/avenir/notification-service/internal/infra/httpclient"
	"github.com/avenir/notification-service/internal/infra/redis"
)

type NotificationProcessor struct {
	redis     redis.RedisService
	ws        *ws.WebSocketSender
	push      push.PushSender
	store     repository.Store
	userStore *httpclient.UserClient
}

func New(
	redis redis.RedisService,
	ws *ws.WebSocketSender,
	push push.PushSender,
	store repository.Store,
	userStore *httpclient.UserClient,
) *NotificationProcessor {
	return &NotificationProcessor{redis, ws, push, store, userStore}
}

func (p *NotificationProcessor) HandleMessageCreated(data []byte) error {

	var msg model.MessageCreatedEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Invalid payload: %v", err)
		return err
	}

	// Получаем информацию о пользователе
	user, err := p.userStore.GetUser(context.Background(), msg.ReceiverID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
	}

	if user == nil {
		log.Printf("User not found: %d", msg.ReceiverID)
		return nil
	}

	// // Проверяем, есть ли у пользователя токен для пуш-уведомлений
	// if user.PushToken == "" {
	// 	log.Printf("User %d has no push token", msg.ReceiverID)
	// 	return nil
	// }

	// Сохраняем уведомление в базе данных
	if err := p.store.Notification().Create(context.Background(), &msg); err != nil {
		log.Printf("Failed to save notification: %v", err)
		return err
	}

	// онлайн ли получатель
	instanceID, err := p.redis.GetInstanceID(msg.ReceiverID)
	println("test")
	if instanceID != "" && err == nil {
		println("test2")
		// Онлайн → через WebSocket
		return p.ws.Send(instanceID, msg)

	}

	// Оффлайн → через Push
	return p.push.Send(msg)
}
