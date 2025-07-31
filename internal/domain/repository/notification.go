package repository

import (
	"context"

	"github.com/avenir/notification-service/internal/domain/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *model.MessageCreatedEvent) error
	FindByUserID(ctx context.Context, userID int) ([]model.MessageCreatedEvent, error)
	MarkAsRead(ctx context.Context, id int) error
}
