package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/katelinlis/NotificationService/internal/domain/model"
	"github.com/katelinlis/NotificationService/internal/domain/repository"
)

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *model.MessageCreatedEvent) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notifications (user_id, message, read,from_id,type, created_at)
		VALUES ($1, $2, false,$3,$4, $5)
	`, n.ReceiverID, n.Content, n.FromID, n.Type, time.Unix(n.CreatedAt, 0))
	return err
}

func (r *notificationRepository) FindByUserID(context context.Context, id int) ([]model.MessageCreatedEvent, error) {
	rows, err := r.db.QueryContext(context, `
		SELECT id, user_id, message, read,from_id,type, created_at
		FROM notifications
		WHERE user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.MessageCreatedEvent
	for rows.Next() {
		var n model.MessageCreatedEvent
		var createdAt sql.NullTime
		if err := rows.Scan(&n.ID, &n.ReceiverID, &n.Content, &n.IsRead, &n.FromID, &n.Type, &createdAt); err != nil {
			return nil, err
		}

		n.CreatedAt = createdAt.Time.Unix()

		notifications = append(notifications, n)
	}

	return notifications, rows.Err()
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notifications
		SET read = true
		WHERE id = $1
	`, id)
	return err
}

// другие методы реализуются по аналогии
