package smtp

import (
	"context"
	"fmt"
	"strings"

	gomail "gopkg.in/mail.v2"
)

type Config struct {
	Host     string
	Username string
	Password string
	Port     int
	From     string
	Name     string
}

type SMTPEmailNotifier struct {
	dialer *gomail.Dialer
	cfg    Config
}

func New(cfg Config) (SMTPEmailNotifier, error) {
	return SMTPEmailNotifier{
		cfg:    cfg,
		dialer: gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password),
	}, nil
}

type Notification struct {
	Recipient string            `json:"recipient"`
	Subject   string            `json:"subject"`
	Template  string            `json:"template"`
	Data      map[string]string `json:"data"`
}

func (e SMTPEmailNotifier) Send(ctx context.Context, n Notification) error {

	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetAddressHeader("From", e.cfg.From, e.cfg.Name)
	message.SetHeader("To", n.Recipient)
	message.SetHeader("Subject", n.Subject)

	// Set email body

	body := e.buildTemplate(n.Template, n.Data)
	message.SetBody("text/html", body)
	message.SetHeader("charset", "UTF-8")

	// Send the email
	return e.dialer.DialAndSend(message)
}

func (e *SMTPEmailNotifier) buildTemplate(template string, values map[string]string) string {
	message := template
	for k, v := range values {
		message = strings.Replace(message, fmt.Sprintf("{{%s}}", k), v, -1)
	}
	return message
}
