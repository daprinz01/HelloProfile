package models

import "github.com/google/uuid"

//UserLanguage is used to retrieve the language information of a user
type UserLanguage struct {
	Language    string `json:"language,omitempty"`
	Proficiency string `json:"proficiency,omitempty"`
}

//Timezone is used to retrieve the timezone of the user
type Timezone struct {
	Timezone string `json:"timezone,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

// Country is used to retrieve the country information
type Country struct {
	Country     string `json:"country,omitempty"`
	FlagURL     string `json:"flagUrl,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

// Application is used to retrieve application information
type Application struct {
	Application string `json:"application,omitempty"`
	Description string `json:"description,omitempty"`
}

// Role is used to retrieve role information
type Role struct {
	Role        string `json:"role,omitempty"`
	Description string `json:"description,omitempty"`
}

// Profile is used to retrieve profile information
type Profile struct {
	ID           uuid.UUID    `json:"id,omitempty"`
	Status       bool         `json:"status,omitempty"`
	ProfileName  string       `json:"profileName,omitempty"`
	IsDefault    bool         `json:"isDefault,omitempty"`
	PageColor    string       `json:"pageColor,omitempty"`
	Font         string       `json:"font,omitempty"`
	Url          string       `json:"url,omitempty"`
	Basic        Basic        `json:"basic,omitempty"`
	ContactBlock ContactBlock `json:"contact,omitempty"`
	Socials      []Socials    `json:"socials,omitempty"`
	Contents     []Content    `json:"contents,omitempty"`
}

type ProfileUrlRequest struct {
	ProfileId   uuid.UUID `json:"profileId,omitempty"`
	ProfileName string    `json:"profileName,omitempty"`
}

// Socials keeps the social accounts a user has added to their profile
type Socials struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Username    string    `json:"username,omitempty"`
	Name        string    `json:"name,omitempty"`
	Placeholder string    `json:"placeholder,omitempty"`
	ImageUrl    string    `json:"imageUrl,omitempty"`
	Order       int32     `json:"order,omitempty"`
	SocialsID   uuid.UUID `json:"socialsID,omitempty"`
	ProfileID   uuid.UUID `json:"profileID,omitempty"`
}

// // Socials keeps the social accounts a user has added to their profile
// type MultipleSocials struct {
// 	Socials []Socials `json:"socials,omitempty"`
// }

type CallToAction struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Type        string    `json:"type,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
}

type Content struct {
	ID             uuid.UUID    `json:"id,omitempty"`
	Title          string       `json:"title,omitempty"`
	Description    string       `json:"description,omitempty"`
	Url            string       `json:"url,omitempty"`
	Order          int32        `json:"order,omitempty"`
	Type           string       `json:"type,omitempty"`
	CallToAction   CallToAction `json:"callToAction,omitempty"`
	CallToActionID uuid.UUID    `json:"callToActionId,omitempty"`
	ContentID      uuid.UUID    `json:"contentId,omitempty"`
}

type ContentType struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Type     string    `json:"type,omitempty"`
	ImageUrl string    `json:"imageUrl,omitempty"`
}

type Basic struct {
	ID              uuid.UUID `json:"id,omitempty"`
	ProfilePhotoUrl string    `json:"profilePhotoUrl,omitempty"`
	CoverPhotoUrl   string    `json:"coverPhotoUrl,omitempty"`
	CoverColour     string    `json:"coverColour,omitempty"`
	Fullname        string    `json:"fullname,omitempty"`
	Title           string    `json:"title,omitempty"`
	Bio             string    `json:"bio,omitempty"`
}
type ContactBlock struct {
	ID      uuid.UUID `json:"id,omitempty"`
	Phone   string    `json:"phone,omitempty"`
	Email   string    `json:"email,omitempty"`
	Address string    `json:"address,omitempty"`
	Website string    `json:"website,omitempty"`
}
type AddContact struct {
	ProfileID  uuid.UUID `json:"profileId,omitempty"`
	CategoryID uuid.UUID `json:"categoryId,omitempty"`
}
type Contact struct {
	Profile Profile `json:"profile,omitempty"`
}

type SaveProfileRequest struct {
	Firstname     string    `json:"firstname,omitempty"`
	Lastname      string    `json:"lastname,omitempty"`
	Email         string    `json:"email,omitempty"`
	ProfileId     uuid.UUID `json:"profileId,omitempty"`
	IsTermsAgreed bool      `json:"isTermsAgreed,omitempty"`
}

type AddProfileFromTemplateRequest struct {
	TemplateID uuid.UUID `json:"templateId,omitempty"`
}
