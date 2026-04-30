package services

import (
	"os"

	"github.com/resend/resend-go/v2"
)

func SendEmail(to string, subject string, body string) error {
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Flowtap <onboarding@resend.dev>",
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	_, err := client.Emails.Send(params)

	return err
}