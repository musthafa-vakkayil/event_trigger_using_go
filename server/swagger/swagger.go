package swagger

// User represents a user in the system
// @Description UserDto object used for input and output
type UserDto struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}
