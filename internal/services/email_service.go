package services

import (
	"fmt"
	"starter-kit-restapi-gofiber/internal/config"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	Config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{Config: cfg}
}

func (s *EmailService) SendEmail(to, subject, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.Config.EmailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(s.Config.SMTPHost, s.Config.SMTPPort, s.Config.SMTPUsername, s.Config.SMTPPassword)
	return d.DialAndSend(m)
}

func (s *EmailService) SendResetPasswordEmail(to, token string) error {
	subject := "Reset Password"
	resetURL := fmt.Sprintf("http://localhost:%s/v1/auth/reset-password?token=%s", s.Config.Port, token)
	text := fmt.Sprintf("Dear user,\nTo reset your password, click on this link: %s\nIf you did not request any password resets, then ignore this email.", resetURL)
	return s.SendEmail(to, subject, text)
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	verifyURL := fmt.Sprintf("http://localhost:%s/v1/auth/verify-email?token=%s", s.Config.Port, token)
	text := fmt.Sprintf("Dear user,\nTo verify your email, click on this link: %s\nIf you did not create an account, then ignore this email.", verifyURL)
	return s.SendEmail(to, subject, text)
}