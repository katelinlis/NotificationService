package push

import "github.com/katelinlis/NotificationService/internal/domain/model"

type PushSender struct{}

func (p PushSender) Send(msg model.MessageCreatedEvent) error {
	return nil
}

func InitPush() PushSender {
	return PushSender{}
}
