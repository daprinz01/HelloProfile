package models

//UserLanguage is used to retrieve the language information of a user
type UserLanguage struct {
	Language    string      `json:"language"`
	Proficiency interface{} `json:"proficiency"`
}
