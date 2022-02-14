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
	PageColor    string       `json:"page_color,omitempty"`
	Font         string       `json:"font,omitempty"`
	Basic        Basic        `json:"basic,omitempty"`
	ContactBlock ContactBlock `json:"contact,omitempty"`
	Socials      []Socials    `json:"socials,omitempty"`
	Links        []Content    `json:"links,omitempty"`
	Articles     []Content    `json:"articles,omitempty"`
	Videos       []Content    `json:"videos,omitempty"`
	Audios       []Content    `json:"audios,omitempty"`
	Forms        []Content    `json:"forms,omitempty"`
	Meetings     []Content    `json:"meetings,omitempty"`
	Events       []Content    `json:"events,omitempty"`
}

// Socials keeps the social accounts a user has added to their profile
type Socials struct {
	Username    string    `json:"username,omitempty"`
	Name        string    `json:"name,omitempty"`
	Placeholder string    `json:"placeholder,omitempty"`
	ImageUrl    string    `json:"image_url,omitempty"`
	Order       int32     `json:"order,omitempty"`
	SocialsID   uuid.UUID `json:"socialsID,omitempty"`
	ProfileID   uuid.UUID `json:"profileID,omitempty"`
}

type CallToAction struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	DisplayName string    `json:"displayName"`
}

type Content struct {
	ID             uuid.UUID    `json:"id"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Url            string       `json:"url"`
	Order          int32        `json:"order"`
	CallToAction   CallToAction `json:"callToAction"`
	CallToActionID uuid.UUID    `json:"callToActionId"`
	ContentID      uuid.UUID    `json:"contentId"`
}

type ContentType struct {
	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	ImageUrl string    `json:"imageUrl"`
}

type Basic struct {
	ID              uuid.UUID `json:"id"`
	ProfilePhotoUrl string    `json:"profile_photo_url"`
	CoverPhotoUrl   string    `json:"cover_photo_url"`
	CoverColour     string    `json:"cover_colour"`
	Fullname        string    `json:"fullname"`
	Title           string    `json:"title"`
	Bio             string    `json:"bio"`
}
type ContactBlock struct {
	ID      uuid.UUID `json:"id"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
	Address string    `json:"address"`
	Website string    `json:"website"`
}
