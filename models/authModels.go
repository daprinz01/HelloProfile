package models

// LoginRequest is used to contruct the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
