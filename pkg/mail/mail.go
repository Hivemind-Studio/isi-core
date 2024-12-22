package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	SenderEmail string
	SenderName  string
}

type EmailClient struct {
	config *EmailConfig
}

func NewEmailClient(config *EmailConfig) *EmailClient {
	return &EmailClient{
		config: config,
	}
}

func (c *EmailClient) SendMail(to []string, subject string, templateFile string, data any) error {
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", templateFile)
	}

	emailTemplate, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var body bytes.Buffer
	err = emailTemplate.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fmt.Sprintf("%s <%s>", c.config.SenderName, c.config.SenderEmail))
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body.String())

	port, err := strconv.Atoi(c.config.Port)
	if err != nil {
		return fmt.Errorf("failed to parse SMTP port: %v", err)
	}

	dialer := gomail.NewDialer(
		c.config.Host,
		port,
		c.config.Username,
		c.config.Password,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	log.Printf("Mail sent! to %s", strings.Join(to, ";"))
	return nil
}
