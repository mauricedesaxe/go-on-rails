package jobs

import (
	"fmt"
	"os"

	"github.com/mailjet/mailjet-apiv3-go"
)

type EmailJob struct {
	From    string
	To      string
	Subject string
	Body    string
}

func (e EmailJob) Execute() error {
	// Validate the email job
	from := e.From
	to := e.To
	subject := e.Subject
	body := e.Body
	if from == "" || to == "" || subject == "" || body == "" {
		return fmt.Errorf("from, to, subject, and body are required")
	}

	// Initialize the Mailjet client
	mjPublicKey, ok := os.LookupEnv("MJ_APIKEY_PUBLIC")
	if !ok {
		return fmt.Errorf("mj public key not found")
	}
	mjPrivateKey, ok := os.LookupEnv("MJ_APIKEY_PRIVATE")
	if !ok {
		return fmt.Errorf("mj private key not found")
	}
	mailjetClient := mailjet.NewMailjetClient(mjPublicKey, mjPrivateKey)

	// Send the email
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: from,
				Name:  from,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
					Name:  to,
				},
			},
			Subject:  subject,
			TextPart: body,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	return err
}
