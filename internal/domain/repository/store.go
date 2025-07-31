package repository

type Store interface {
	Notification() NotificationRepository
}
