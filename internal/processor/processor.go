package processor

import (
	"encoding/json"
	"log"

	"github.com/avenir/notification-service/internal/delivery/push"
	"github.com/avenir/notification-service/internal/delivery/ws"
	"github.com/avenir/notification-service/internal/domain/model"
	"github.com/avenir/notification-service/internal/infra/redis"
)

type NotificationProcessor struct {
	redis redis.RedisService
	ws    *ws.WebSocketSender
	push  push.PushSender
}

func New(redis redis.RedisService, ws *ws.WebSocketSender, push push.PushSender) *NotificationProcessor {
	return &NotificationProcessor{redis, ws, push}
}

func (p *NotificationProcessor) HandleMessageCreated(data []byte) error {
	var msg model.MessageCreatedEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Invalid payload: %v", err)
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
