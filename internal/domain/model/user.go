package model

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}
