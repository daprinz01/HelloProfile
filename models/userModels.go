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

// Socials keeps the social accounts a user has added to their profile
type Socials struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Placeholder string `json:"placeholder"`
	ImageUrl    string `json:"image_url"`
	Order       int32  `json:"order"`
}

type CallToAction struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	DisplayName string    `json:"displayName"`
}

type Content struct {
	ID           uuid.UUID    `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Url          string       `json:"url"`
	Order        int32        `json:"order"`
	CallToAction CallToAction `json:"callToAction"`
}
