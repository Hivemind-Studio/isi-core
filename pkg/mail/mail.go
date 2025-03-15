package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type EmailConfig struct {
	APIKey string
	Sender string
	Email  string
	Url    string
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

	emailPayload := map[string]interface{}{
		"sender": map[string]string{
			"name":  c.config.Sender,
			"email": c.config.Email,
		},
		"to": []map[string]string{
			{
				"email": to[0],
				"name":  to[0],
			},
		},
		"subject":     subject,
		"htmlContent": body.String(),
	}

	jsonPayload, err := json.Marshal(emailPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", c.config.Url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", c.config.APIKey)

	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("failed to send email to %s cause %s", to[0], err)
	}

	log.Printf("Mail sent! to %s", to[0])
	return nil
}
