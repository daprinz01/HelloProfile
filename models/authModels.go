package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// LoginRequest is used to contruct the login request
type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Claims is  a struct that will be encoded to a JWT.
// jwt.StandardClaims is an embedded type to provide expiry time
type Claims struct {
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
	jwt.StandardClaims
}

// VerifiedClaims is the result from verifying claims
type VerifiedClaims struct {
	Email string
	Role  string
}

// UserDetail is used to create a new User
type UserDetail struct {
	Firstname                 string    `json:"firstname,omitempty"`
	Lastname                  string    `json:"lastname,omitempty"`
	Username                  string    `json:"username,omitempty"`
	Email                     string    `json:"email,omitempty"`
	Phone                     string    `json:"phone,omitempty"`
	IsEmailConfirmed          bool      `json:"isEmailConfirmed,omitempty"`
	Password                  string    `json:"password,omitempty"`
	IsPasswordSystemGenerated bool      `json:"isPasswordSystemGenerated,omitempty"`
	CreatedAt                 time.Time `json:"createdAt,omitempty"`
	IsLockedOut               bool      `json:"isLockedOut,omitempty"`
	ProfilePicture            string    `json:"profilePicture,omitempty"`
	IsActive                  bool      `json:"isActive,omitempty"`
	LanguageName              string    `json:"languageName,omitempty"`
	TimezoneName              string    `json:"timezoneName,omitempty"`
	Zone                      string    `json:"zone,omitempty"`
	Country                   string    `json:"country,omitempty"`
	City                      string    `json:"city,omitempty"`
	Profiles                  []Profile `json:"profiles,omitempty"`
}

// RefreshResponse is used to send success message for a successful refresh of auth token
type RefreshResponse struct {
	ResponseCode    string `json:"responseCode,omitempty"`
	ResponseMessage string `json:"responseMessage,omitempty"`
}

// SendOtpRequest is used to recieve otp requests
type SendOtpRequest struct {
	Email           string `json:"email,omitempty"`
	Phone           string `json:"phonenumber,omitempty"`
	IsEmailPrefered bool   `json:"isEmailPreferred"`
	Purpose         string `json:"purpose,omitempty"`
	Application     string `json:"application,omitempty"`
	VerifyPath      string `json:"verifyPath,omitempty"`
}

// VerifyOtpRequest is used to verify an OTP against a user
type VerifyOtpRequest struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phonenumber,omitempty"`
	OTP   string `json:"otp,omitempty"`
}

// ResetPasswordRequest is used to reset user password after
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword,omitempty"`
}

// ChangePasswordRequest is used to change the user's password with new password
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// GoogleClaims -
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

type GoogleJWT struct {
	GoogleJWT string `json:"token"`
}
type FeedbackRequest struct {
	Sender        string   `json:"sender,omitempty"`
	Message       string   `json:"message,omitempty"`
	AttachmentUrl []string `json:"attachments,omitempty"`
}
