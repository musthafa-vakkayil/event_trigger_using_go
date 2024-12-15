package model

// User represents a user in the system
// @Description User object used for input and output
type User struct {
	UserId       string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

// @Description UserDto object used for input and output
type UserDto struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}
