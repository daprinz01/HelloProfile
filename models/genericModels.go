package models

// Errormessage is used to construct error messages
type Errormessage struct {
	Errorcode    string `json:"responseCode,omitempty"`
	ErrorMessage string `json:"responseDescription,omitempty"`
}

// SuccessResponse is used to form the success message
type SuccessResponse struct {
	ResponseCode    string      `json:"responseCode,omitempty"`
	ResponseMessage string      `json:"responseDescription,omitempty"`
	ResponseDetails interface{} `json:"responseDetails,omitempty"`
}
