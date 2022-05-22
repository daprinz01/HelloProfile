// Code generated by sqlc. DO NOT EDIT.

package helloprofiledb

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type BasicBlock struct {
	ID              uuid.UUID      `json:"id"`
	ProfilePhotoUrl sql.NullString `json:"profile_photo_url"`
	CoverPhotoUrl   sql.NullString `json:"cover_photo_url"`
	CoverColour     sql.NullString `json:"cover_colour"`
	Fullname        string         `json:"fullname"`
	Title           string         `json:"title"`
	Bio             string         `json:"bio"`
}

type CallToAction struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	DisplayName string    `json:"display_name"`
}

type Contact struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProfileID uuid.UUID `json:"profile_id"`
}

type ContactBlock struct {
	ID      uuid.UUID `json:"id"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
	Address string    `json:"address"`
	Website string    `json:"website"`
}

type ContactCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Content struct {
	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	ImageUrl string    `json:"image_url"`
}

type EmailVerification struct {
	ID        uuid.UUID      `json:"id"`
	Email     sql.NullString `json:"email"`
	Otp       string         `json:"otp"`
	CreatedAt time.Time      `json:"created_at"`
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
	BasicBlockID   uuid.NullUUID  `json:"basic_block_id"`
	ContactBlockID uuid.NullUUID  `json:"contact_block_id"`
	Status         bool           `json:"status"`
	ProfileName    string         `json:"profile_name"`
	PageColor      string         `json:"page_color"`
	Font           string         `json:"font"`
	Url            sql.NullString `json:"url"`
	IsDefault      bool           `json:"is_default"`
}

type ProfileContent struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	DisplayTitle   string    `json:"display_title"`
	Description    string    `json:"description"`
	Url            string    `json:"url"`
	ProfileID      uuid.UUID `json:"profile_id"`
	CallToActionID uuid.UUID `json:"call_to_action_id"`
	ContentID      uuid.UUID `json:"content_id"`
	Order          int32     `json:"order"`
}

type ProfileSocial struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	SocialsID uuid.UUID `json:"socials_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	Order     int32     `json:"order"`
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

type SavedProfile struct {
	ID        uuid.UUID `json:"id"`
	ProfileID uuid.UUID `json:"profile_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	IsAdded   bool      `json:"is_added"`
}

type Social struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Placeholder string    `json:"placeholder"`
	ImageUrl    string    `json:"image_url"`
}

type User struct {
	ID                        uuid.UUID      `json:"id"`
	Firstname                 sql.NullString `json:"firstname"`
	Lastname                  sql.NullString `json:"lastname"`
	Username                  sql.NullString `json:"username"`
	Email                     string         `json:"email"`
	Phone                     sql.NullString `json:"phone"`
	Country                   sql.NullString `json:"country"`
	City                      sql.NullString `json:"city"`
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
	Email                     string         `json:"email"`
	Phone                     sql.NullString `json:"phone"`
	Username                  sql.NullString `json:"username"`
	Password                  sql.NullString `json:"password"`
	Country                   sql.NullString `json:"country"`
	City                      sql.NullString `json:"city"`
	ProfilePicture            sql.NullString `json:"profile_picture"`
	IsEmailConfirmed          bool           `json:"is_email_confirmed"`
	IsLockedOut               bool           `json:"is_locked_out"`
	IsPasswordSystemGenerated bool           `json:"is_password_system_generated"`
	CreatedAt                 time.Time      `json:"created_at"`
	IsActive                  bool           `json:"is_active"`
	RoleName                  sql.NullString `json:"role_name"`
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

type UserRole struct {
	ID     uuid.UUID     `json:"id"`
	UserID uuid.NullUUID `json:"user_id"`
	RoleID uuid.NullUUID `json:"role_id"`
}
