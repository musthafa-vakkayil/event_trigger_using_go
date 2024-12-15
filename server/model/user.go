package model

type User struct {
	UserId       int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}
