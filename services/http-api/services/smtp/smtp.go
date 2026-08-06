package smtp

import (
	"fmt"
	"log"
	"net/smtp"
)

type GruentSmtpClient interface {
	SendMail(from string, to []string, subject string, message []byte) error
}

type gruentSmtpClient struct {
	smtpAddress string
	smtpAuth    smtp.Auth
}

func (c *gruentSmtpClient) SendMail(from string, to []string, subject string, message []byte) error {
	return smtp.SendMail(c.smtpAddress, c.smtpAuth, from, to, message)
}

type mockGruentSmtpClient struct{}

func (m *mockGruentSmtpClient) SendMail(from string, to []string, subject string, message []byte) error {
	log.Printf("Sending mock email - from: %s - to: %s - subject: %s", from, to, subject)
	log.Println(string(message))
	return nil
}

func NewSmtpClient(smtpHost string, smtpPort int, smtpUser *string, smtpPassword *string) GruentSmtpClient {
	smtpAddress := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	smtpAuth := (func() smtp.Auth {
		if smtpUser == nil || smtpPassword == nil {
			return nil
		}
		return smtp.PlainAuth("", *smtpUser, *smtpPassword, smtpHost)
	})()

	return &gruentSmtpClient{
		smtpAddress: smtpAddress,
		smtpAuth:    smtpAuth,
	}
}

func NewMockSmtpClient() GruentSmtpClient {
	return &mockGruentSmtpClient{}
}

var SmtpClient GruentSmtpClient
