package push

import "github.com/avenir/notification-service/internal/domain/model"

type PushSender struct{}

func (p PushSender) Send(msg model.MessageCreatedEvent) error {
	return nil
}

func InitPush() PushSender {
	return PushSender{}
}
