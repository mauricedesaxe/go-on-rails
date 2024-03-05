package jobs

import (
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go"
)

type EmailJob struct {
	From    string
	To      string
	Subject string
	Body    string
	Client  *mailjet.Client
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
	_, err := e.Client.SendMailV31(&messages)
	return err
}
