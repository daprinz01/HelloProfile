package models

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
