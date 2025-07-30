package model

import "github.com/golang-jwt/jwt/v5"

type MessageCreatedEvent struct {
	ReceiverID int    `json:"reciver_id"`
	MessageID  int    `json:"msg_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Timestamp  int64  `json:"timestamp"`
	FromID     int    `json:"from_id"`
	Username   string
	Avatar     string
}

type MyCustomClaims struct {
	ClientID int64 `json:"id"`
	jwt.RegisteredClaims
}
