// Code generated by sqlc. DO NOT EDIT.

package authdb

import (
	"database/sql"
	"time"
)

type Application struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ApplicationsRole struct {
	ID             int64         `json:"id"`
	ApplicationsID sql.NullInt64 `json:"applications_id"`
	RolesID        sql.NullInt64 `json:"roles_id"`
}

type Country struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	FlagImageUrl sql.NullString `json:"flag_image_url"`
}

type IdentityProvider struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ImageUrl     string `json:"image_url"`
}

type Language struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type State struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	CountryID sql.NullInt64 `json:"country_id"`
}

type Timezone struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Zone string `json:"zone"`
}

type User struct {
	ID                        int64          `json:"id"`
	Firstname                 sql.NullString `json:"firstname"`
	Lastname                  sql.NullString `json:"lastname"`
	Username                  sql.NullString `json:"username"`
	Email                     string         `json:"email"`
	IsEmailConfirmed          bool           `json:"is_email_confirmed"`
	Password                  sql.NullString `json:"password"`
	IsPasswordSystemGenerated bool           `json:"is_password_system_generated"`
	Address                   sql.NullString `json:"address"`
	City                      sql.NullString `json:"city"`
	State                     sql.NullString `json:"state"`
	Country                   sql.NullString `json:"country"`
	CreatedAt                 time.Time      `json:"created_at"`
	IsLockedOut               bool           `json:"is_locked_out"`
	ImageUrl                  sql.NullString `json:"image_url"`
	IsActive                  bool           `json:"is_active"`
}

type UserDetail struct {
	ID                        int64          `json:"id"`
	Firstname                 sql.NullString `json:"firstname"`
	Lastname                  sql.NullString `json:"lastname"`
	Username                  sql.NullString `json:"username"`
	Email                     string         `json:"email"`
	IsEmailConfirmed          bool           `json:"is_email_confirmed"`
	Password                  sql.NullString `json:"password"`
	IsPasswordSystemGenerated bool           `json:"is_password_system_generated"`
	Address                   sql.NullString `json:"address"`
	City                      sql.NullString `json:"city"`
	State                     sql.NullString `json:"state"`
	Country                   sql.NullString `json:"country"`
	CreatedAt                 time.Time      `json:"created_at"`
	IsLockedOut               bool           `json:"is_locked_out"`
	ProfilePicture            sql.NullString `json:"profile_picture"`
	IsActive                  bool           `json:"is_active"`
	LanguageName              sql.NullString `json:"language_name"`
	RoleName                  sql.NullString `json:"role_name"`
	TimezoneName              sql.NullString `json:"timezone_name"`
	Zone                      sql.NullString `json:"zone"`
	ProviderName              sql.NullString `json:"provider_name"`
	ClientID                  sql.NullString `json:"client_id"`
	ClientSecret              sql.NullString `json:"client_secret"`
	ProviderLogo              sql.NullString `json:"provider_logo"`
}

type UserLanguage struct {
	ID         int64         `json:"id"`
	UserID     sql.NullInt64 `json:"user_id"`
	LanguageID sql.NullInt64 `json:"language_id"`
}

type UserProvider struct {
	ID                 int64         `json:"id"`
	UserID             sql.NullInt64 `json:"user_id"`
	IdentityProviderID sql.NullInt64 `json:"identity_provider_id"`
}

type UserRole struct {
	ID     int64         `json:"id"`
	UserID sql.NullInt64 `json:"user_id"`
	RoleID sql.NullInt64 `json:"role_id"`
}

type UserTimezone struct {
	ID         int64         `json:"id"`
	UserID     sql.NullInt64 `json:"user_id"`
	TimezoneID sql.NullInt64 `json:"timezone_id"`
}
