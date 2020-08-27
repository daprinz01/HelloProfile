package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// LoginRequest is used to contruct the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims is  a struct that will be encoded to a JWT.
// jwt.StandardClaims is an embedded type to provide expiry time
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// VerifiedClaims is the result from verifying claims
type VerifiedClaims struct {
	Email string
	Role  string
}

// UserDetail is used to create a new User
type UserDetail struct {
	Firstname                 string    `json:"firstname"`
	Lastname                  string    `json:"lastname"`
	Username                  string    `json:"username"`
	Email                     string    `json:"email"`
	IsEmailConfirmed          bool      `json:"is_email_confirmed"`
	Password                  string    `json:"password"`
	IsPasswordSystemGenerated bool      `json:"is_password_system_generated"`
	Address                   string    `json:"address"`
	City                      string    `json:"city"`
	State                     string    `json:"state"`
	Country                   string    `json:"country"`
	CreatedAt                 time.Time `json:"created_at"`
	IsLockedOut               bool      `json:"is_locked_out"`
	ProfilePicture            string    `json:"profile_picture"`
	IsActive                  bool      `json:"is_active"`
	LanguageName              string    `json:"language_name"`
	RoleName                  string    `json:"role_name"`
	TimezoneName              string    `json:"timezone_name"`
	Zone                      string    `json:"zone"`
	ProviderName              string    `json:"provider_name"`
	ClientID                  string    `json:"client_id"`
	ClientSecret              string    `json:"client_secret"`
	ProviderLogo              string    `json:"provider_logo"`
}

// RefreshResponse is used to send success message for a successful refresh of auth token
type RefreshResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}
