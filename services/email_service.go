package services

import (
	"attendance-system/backend/config"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailService interface {
	Send(to, subject, htmlBody string) error
	IsConfigured() bool
}

type noopEmailService struct{}

func NewEmailService(cfg config.Config) EmailService {
	if cfg.SMTPHost == "" || cfg.SMTPUsername == "" || cfg.SMTPPassword == "" || cfg.SMTPFromEmail == "" {
		return noopEmailService{}
	}
	return &smtpEmailService{cfg: cfg}
}

func (noopEmailService) Send(to, subject, htmlBody string) error {
	return nil
}

func (noopEmailService) IsConfigured() bool {
	return false
}

type smtpEmailService struct {
	cfg config.Config
}

func (s *smtpEmailService) IsConfigured() bool {
	return true
}

func (s *smtpEmailService) Send(to, subject, htmlBody string) error {
	fromName := s.cfg.SMTPFromName
	if fromName == "" {
		fromName = "Attendance System"
	}

	message := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		fromName,
		s.cfg.SMTPFromEmail,
		to,
		subject,
		htmlBody,
	)

	address := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	client, err := smtp.Dial(address)
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: s.cfg.SMTPHost}); err != nil {
			return err
		}
	}
	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(s.cfg.SMTPFromEmail); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(message)); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return client.Quit()
}
