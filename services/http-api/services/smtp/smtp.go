package smtp

import (
	"fmt"
	"log"
	"net/smtp"
)

type IoteaSmtpClient interface {
	SendMail(from string, to []string, subject string, message []byte) error
}

type ioteaSmtpClient struct {
	smtpAddress string
	smtpAuth    smtp.Auth
}

func (c *ioteaSmtpClient) SendMail(from string, to []string, subject string, message []byte) error {
	return smtp.SendMail(c.smtpAddress, c.smtpAuth, from, to, message)
}

type mockIoteaSmtpClient struct{}

func (m *mockIoteaSmtpClient) SendMail(from string, to []string, subject string, message []byte) error {
	log.Printf("Sending mock email - from: %s - to: %s - subject: %s", from, to, subject)
	log.Println(string(message))
	return nil
}

func NewSmtpClient(smtpHost string, smtpPort int, smtpUser *string, smtpPassword *string) IoteaSmtpClient {
	smtpAddress := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	smtpAuth := (func() smtp.Auth {
		if smtpUser == nil || smtpPassword == nil {
			return nil
		}
		return smtp.PlainAuth("", *smtpUser, *smtpPassword, smtpHost)
	})()

	return &ioteaSmtpClient{
		smtpAddress: smtpAddress,
		smtpAuth:    smtpAuth,
	}
}

func NewMockSmtpClient() IoteaSmtpClient {
	return &mockIoteaSmtpClient{}
}

var SmtpClient IoteaSmtpClient
