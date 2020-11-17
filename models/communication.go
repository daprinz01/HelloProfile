package models

// SendEmailRequest is used to get email request
type SendEmailRequest struct {
	From           EmailAddress      `json:"from,omitempty"`
	To             []EmailAddress    `json:"to,omitempty"`
	CC             []EmailAddress    `json:"cc,omitempty"`
	BCC            []EmailAddress    `json:"bcc,omitempty"`
	Subject        string            `json:"subject,omitempty"`
	Message        string            `json:"message,omitempty"`
	AttachmentName []EmailAttachment `json:"attachments,omitempty"`
}

// EmailAddress is used to collect name and email of recipients
type EmailAddress struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// EmailAttachment is used to send attachments
type EmailAttachment struct {
	FileName string `json:"fileName,omitempty"`
	Base64   string `json:"base64,omitempty"`
}

// SendSmsRequest is used to receive sms request
type SendSmsRequest struct {
	Phone   string `json:"phonenumber,omitempty"`
	Message string `json:"message,omitempty"`
}
