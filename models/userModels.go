package models

//UserLanguage is used to retrieve the language information of a user
type UserLanguage struct {
	Language    string      `json:"language"`
	Proficiency interface{} `json:"proficiency"`
}

//Timezone is used to retrieve the timezone of the user
type Timezone struct {
	Timezone string `json:"timezone"`
	Zone     string `json:"zone"`
}

// Country is used to retrieve the country information
type Country struct {
	Country     string `json:"country"`
	FlagURL     string `json:"flag_url"`
	CountryCode string `json:"country_code"`
}

// Application is used to retrieve application information
type Application struct {
	Application string `json:"application"`
	Description string `json:"description"`
}

// Role is used to retrieve role information
type Role struct {
	Role        string `json:"role"`
	Description string `json:"description"`
}
