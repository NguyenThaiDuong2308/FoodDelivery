package service

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"
	"user-service/config"
)

type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPMailer(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	}
}

type Mailer interface {
	SendResetPasswordEmail(ctx context.Context, toEmail, token string) error
}

func (m *SMTPMailer) SendResetPasswordEmail(ctx context.Context, toEmail string, token string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	subject := "Reset Password"
	resetLink := fmt.Sprintf("http://192.168.237.130:5173/reset-forgot-password?token=%s", token)
	body := fmt.Sprintf("Click the link to reset your password: \n\n%s", resetLink)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
		toEmail, subject, body,
	))
	port, err := strconv.Atoi(m.Port)
	if err != nil {
		return errors.New("port must be an integer")
	}
	addr := fmt.Sprintf("%s:%d", m.Host, port)
	err = smtp.SendMail(addr, auth, m.From, []string{toEmail}, msg)
	if err != nil {
		return err
	}
	return nil
}
