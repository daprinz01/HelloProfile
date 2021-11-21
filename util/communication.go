package util

import (
	"fmt"
	"net/smtp"
	"os"

	log "github.com/sirupsen/logrus"
)

// SendEmail is used to send email to users
func SendEmail(from string, to []string, message string) error {
	fields := log.Fields{"microservice": "helloprofile.service", "function": "SendEmail"}
	//Read variables from environment
	var smtpHost, smtpPort, smtpUser, smtpPassword string
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	smtpUser = os.Getenv("SMTP_USER")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	// Set up authentication information.
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	// to := []string{"recipient@example.net"}
	msg := []byte(message)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, from, to, msg)
	if err != nil {
		log.WithFields(fields).Error("Error occured sending Email")
		return err
	}
	return nil
}
