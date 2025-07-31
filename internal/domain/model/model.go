package model

import "github.com/golang-jwt/jwt/v5"

type MessageCreatedEvent struct {
	ID         int    `json:"id"`
	ReceiverID int    `json:"reciver_id"`
	MessageID  int    `json:"msg_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	CreatedAt  int64  `json:"created_at"`
	FromID     int    `json:"from_id"`
	IsRead     bool   `json:"is_read"`
	Username   string
	Avatar     string
}

type MyCustomClaims struct {
	ClientID int64 `json:"id"`
	jwt.RegisteredClaims
}
