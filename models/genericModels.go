package models

// Errormessage is used to construct error messages
type Errormessage struct {
	Errorcode    string `json:"responseCode"`
	ErrorMessage string `json:"responseDescription"`
}

// SuccessResponse is used to form the success message
type SuccessResponse struct {
	ResponseCode    string      `json:"responseCode"`
	ResponseMessage string      `json:"responseDescription"`
	ResponseDetails interface{} `json:"responseDetails"`
}
