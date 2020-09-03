package models

// SendEmailRequest is used to get email request
type SendEmailRequest struct {
	From           EmailAddress      `json:"from"`
	To             []EmailAddress    `json:"to"`
	CC             []EmailAddress    `json:"cc"`
	BCC            []EmailAddress    `json:"bcc"`
	Subject        string            `json:"subject"`
	Message        string            `json:"message"`
	AttachmentName []EmailAttachment `json:"attachments"`
}

// EmailAddress is used to collect name and email of recipients
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// EmailAttachment is used to send attachments
type EmailAttachment struct {
	FileName string `json:"fileName"`
	Base64   string `json:"base64"`
}

