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

/*
Sample Use

	data := EmailData{
		Name:    "Juan",
		Message: "Have a nice day!",
	}

err := mail.SendMail([]string{"jdniels525@gmail.com"}, "Test mail", "pkg/mail/default.html", mail.EmailData(data))

	if err != nil {
		log.Println("Error sending email:", err)
	}
*/
func SendMail(to []string, subject string, templateFile string, data any) error {
	// Check if the template file exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", templateFile)
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fmt.Sprintf("%s <%s>", os.Getenv("MAIL_SENDER_NAME"), os.Getenv("MAIL_AUTH_EMAIL")))
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body.String())

	port, err := strconv.Atoi(os.Getenv("MAIL_SMTP_PORT"))
	if err != nil {
		return fmt.Errorf("failed to parse SMTP port: %v", err)
	}

	fmt.Println("SMTP Username:", os.Getenv("MAIL_SMTP_USERNAME"))
	fmt.Println("SMTP Password:", os.Getenv("MAIL_SMTP_PASSWORD"))

	dialer := gomail.NewDialer(
		os.Getenv("MAIL_SMTP_HOST"),
		port,
		os.Getenv("MAIL_SMTP_USERNAME"),
		os.Getenv("MAIL_SMTP_PASSWORD"),
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	log.Printf("Mail sent! to %s", strings.Join(to, ";"))
	return nil
}
