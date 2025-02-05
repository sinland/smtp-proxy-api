package smtp

import (
	"fmt"
	"net/smtp"
)

type Service struct {
	SMTPServer string
	SMTPPort   int
	username   string
	password   string
}

func NewService(smtpServer string, smtpPort int, username string, password string) *Service {
	return &Service{SMTPServer: smtpServer, SMTPPort: smtpPort, username: username, password: password}
}

func (s *Service) SendEmail(from, to, subject, body string) error {
	// SMTP authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.SMTPServer)

	// Email headers and body
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

	// Sending email
	err := smtp.SendMail(fmt.Sprintf("%s:%d", s.SMTPServer, s.SMTPPort), auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
