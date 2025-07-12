package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/logging"
)

type EmailPayload struct {
	ClientName     string `json:"client_name"`
	ClientEmail    string `json:"client_email"`
	RecipientName  string `json:"recipient_name"`
	RecipientEmail string `json:"recipient_email"`
	Subject        string `json:"subject"`
	HtmlContent    string `json:"html_content"`
}

func SendEmail(clientName string, clientEmail string, recipientName string, recipientEmail string, subject string, htmlContent string) error {
	if _, err := config.LoadConfig(); err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "QugoPy",
			"email": config.AppConfig.BREVO.EMAIL,
		},
		"to": []map[string]string{
			{
				"email": recipientEmail,
				"name":  recipientName,
			},
		},
		"replyTo": map[string]string{
			"name":  clientName,
			"email": clientEmail,
		},
		"subject":     subject,
		"htmlContent": htmlContent,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", config.AppConfig.BREVO.URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("could not create a POST request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", config.AppConfig.BREVO.API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("email send failed: status %s, body: %s", resp.Status, string(body))
	}

	logging.DebugLog(fmt.Sprintln("Email sent! Status: ", resp.Status))
	return nil
}
