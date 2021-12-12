// Code generated by sqlc. DO NOT EDIT.

package helloprofiledb

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.NullUUID  `json:"user_id"`
	Street      string         `json:"street"`
	City        string         `json:"city"`
	State       sql.NullString `json:"state"`
	CountryID   uuid.NullUUID  `json:"country_id"`
	AddressType uuid.NullUUID  `json:"address_type"`
}

type AddressType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Contact struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	ProfileID         uuid.UUID `json:"profile_id"`
	ContactCategoryID uuid.UUID `json:"contact_category_id"`
}

type ContactCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Country struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	FlagImageUrl sql.NullString `json:"flag_image_url"`
	CountryCode  sql.NullString `json:"country_code"`
}

type EmailVerification struct {
	ID        uuid.UUID      `json:"id"`
	Email     sql.NullString `json:"email"`
	Otp       string         `json:"otp"`
	CreatedAt time.Time      `json:"created_at"`
}

type IdentityProvider struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	ImageUrl     string    `json:"image_url"`
}

type Language struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type LanguageProficiency struct {
	ID          uuid.UUID      `json:"id"`
	Proficiency sql.NullString `json:"proficiency"`
}

type Otp struct {
	ID               uuid.UUID      `json:"id"`
	UserID           uuid.UUID      `json:"user_id"`
	OtpToken         sql.NullString `json:"otp_token"`
	CreatedAt        time.Time      `json:"created_at"`
	IsSmsPreferred   bool           `json:"is_sms_preferred"`
	IsEmailPreferred bool           `json:"is_email_preferred"`
	Purpose          sql.NullString `json:"purpose"`
}

type Profile struct {
	ID             uuid.UUID      `json:"id"`
	UserID         uuid.UUID      `json:"user_id"`
	Status         bool           `json:"status"`
	ProfileName    string         `json:"profile_name"`
	Fullname       string         `json:"fullname"`
	Title          string         `json:"title"`
	Bio            string         `json:"bio"`
	Company        string         `json:"company"`
	CompanyAddress string         `json:"company_address"`
	ImageUrl       sql.NullString `json:"image_url"`
	Phone          string         `json:"phone"`
	Email          string         `json:"email"`
	AddressID      uuid.NullUUID  `json:"address_id"`
	Website        sql.NullString `json:"website"`
	IsDefault      bool           `json:"is_default"`
	Color          sql.NullInt32  `json:"color"`
}

type Recent struct {
	ID         uuid.UUID `json:"id"`
	ProfileID  uuid.UUID `json:"profile_id"`
	Title      string    `json:"title"`
	Highlights string    `json:"highlights"`
	Year       int32     `json:"year"`
	Link       string    `json:"link"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type State struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CountryID uuid.UUID `json:"country_id"`
}

type Timezone struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Zone string    `json:"zone"`
}

type User struct {
	ID                        uuid.UUID      `json:"id"`
	Firstname                 sql.NullString `json:"firstname"`
	Lastname                  sql.NullString `json:"lastname"`
	Username                  sql.NullString `json:"username"`
	Email                     string         `json:"email"`
	Phone                     sql.NullString `json:"phone"`
	IsEmailConfirmed          bool           `json:"is_email_confirmed"`
	Password                  sql.NullString `json:"password"`
	IsPasswordSystemGenerated bool           `json:"is_password_system_generated"`
	CreatedAt                 time.Time      `json:"created_at"`
	IsLockedOut               bool           `json:"is_locked_out"`
	ImageUrl                  sql.NullString `json:"image_url"`
	IsActive                  bool           `json:"is_active"`
}

type UserDetail struct {
	ID                        uuid.UUID      `json:"id"`
	Firstname                 sql.NullString `json:"firstname"`
	Lastname                  sql.NullString `json:"lastname"`
	Username                  sql.NullString `json:"username"`
	Email                     string         `json:"email"`
	Phone                     sql.NullString `json:"phone"`
	IsEmailConfirmed          bool           `json:"is_email_confirmed"`
	Password                  sql.NullString `json:"password"`
	IsPasswordSystemGenerated bool           `json:"is_password_system_generated"`
	CreatedAt                 time.Time      `json:"created_at"`
	IsLockedOut               bool           `json:"is_locked_out"`
	ProfilePicture            sql.NullString `json:"profile_picture"`
	IsActive                  bool           `json:"is_active"`
	TimezoneName              sql.NullString `json:"timezone_name"`
	Zone                      sql.NullString `json:"zone"`
}

type UserLanguage struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	LanguageID  uuid.UUID      `json:"language_id"`
	Proficiency sql.NullString `json:"proficiency"`
}

type UserLogin struct {
	ID                  uuid.UUID      `json:"id"`
	UserID              uuid.UUID      `json:"user_id"`
	LoginTime           time.Time      `json:"login_time"`
	LoginStatus         bool           `json:"login_status"`
	ResponseCode        sql.NullString `json:"response_code"`
	ResponseDescription sql.NullString `json:"response_description"`
	Device              sql.NullString `json:"device"`
	IpAddress           sql.NullString `json:"ip_address"`
	Longitude           sql.NullString `json:"longitude"`
	Latitude            sql.NullString `json:"latitude"`
	Resolved            bool           `json:"resolved"`
}

type UserProvider struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	IdentityProviderID uuid.UUID `json:"identity_provider_id"`
}

type UserRole struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}

type UserTimezone struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	TimezoneID uuid.UUID `json:"timezone_id"`
}
