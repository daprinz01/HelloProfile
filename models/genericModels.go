package models

// Errormessage is used to construct error messages
type Errormessage struct {
	Errorcode    string `json:"responseCode"`
	ErrorMessage string `json:"responseMessage"`
}

// SuccessResponse is used to form the success message
type SuccessResponse struct {
	ResponseCode    string      `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	ResponseDetails interface{} `json:"responseDetails"`
}
