package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type Mailer interface {
	SendLoginCode(ctx context.Context, email, code string, expiresIn time.Duration) error
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SMTPMailer struct {
	config SMTPConfig
}

func NewSMTPMailer(config SMTPConfig) *SMTPMailer {
	return &SMTPMailer{config: config}
}

func (m *SMTPMailer) SendLoginCode(ctx context.Context, email, code string, expiresIn time.Duration) error {
	if m == nil || m.config.Host == "" || m.config.Port == 0 || m.config.Username == "" || m.config.Password == "" {
		return ErrMailUnavailable
	}
	from := strings.TrimSpace(m.config.From)
	if from == "" {
		from = m.config.Username
	}

	subject := "E-Director login verification code"
	body := fmt.Sprintf("您的 E-Director 登录验证码是：%s\n\n验证码 %d 分钟内有效。若非本人操作，请忽略本邮件。", code, int(expiresIn.Minutes()))
	message := []byte(strings.Join([]string{
		"From: " + from,
		"To: " + email,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n"))

	addr := net.JoinHostPort(m.config.Host, fmt.Sprintf("%d", m.config.Port))
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)

	result := make(chan error, 1)
	go func() {
		if m.config.Port == 465 {
			result <- m.sendTLS(addr, auth, from, email, message)
			return
		}
		result <- smtp.SendMail(addr, auth, from, []string{email}, message)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-result:
		return err
	}
}

func (m *SMTPMailer) sendTLS(addr string, auth smtp.Auth, from, to string, message []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: m.config.Host, MinVersion: tls.VersionTLS12})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.config.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}
	return writer.Close()
}
