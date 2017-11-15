package email

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/Sirupsen/logrus"
	mailgun "github.com/mailgun/mailgun-go"
	"github.com/sourcegraph/checkup"
)

// Notifier sends an email notification when something is wrong.
type Notifier struct {
	// MailgunAPIKey stores the API key for Mailgun if configured.
	MailgunAPIKey string
	// MailgunDomain stores the domain for Mailgun if configured.
	MailgunDomain string
	// Recipient is the email address to send the notification to.
	Recipient string
	// Server is the email server.
	Server string
	// Sender is the email address to send the notification from.
	Sender string
	// Auth holds the authentication details for the email server.
	Auth smtp.Auth
}

// Notify checks the health status of the result and sends an email if
// something is not healthy.
func (n Notifier) Notify(results []checkup.Result) error {
	for _, r := range results {
		if !r.Healthy {
			logrus.Debugf("%s is %s: sending email", r.Title, r.Status())
			if err := n.sendEmail(r); err != nil {
				return err
			}
		} else {
			logrus.Debugf("%s is %s", r.Title, r.Status())
		}
	}

	return nil
}

func (n Notifier) sendEmail(result checkup.Result) error {
	if n.MailgunAPIKey != "" && n.MailgunDomain != "" {
		mailgunClient := mailgun.NewMailgun(n.MailgunDomain, n.MailgunAPIKey, "")

		msg, _, err := mailgunClient.Send(mailgunClient.NewMessage(
			/* From */ fmt.Sprintf("%s <%s>", n.Sender, n.Sender),
			/* Subject */ fmt.Sprintf("[UPMAIL]: %s %s", result.Title, result.Status()),
			/* Body */ fmt.Sprintf("Time: %s\n\n%s", time.Now().Format(time.UnixDate), result.String()),
			/* To */ n.Recipient,
		))
		if err != nil {
			return fmt.Errorf("Sending Mailgun message failed: response: %#v error: %v", msg, err)
		}
		logrus.Infof("Mailgun send message succeeded: %#v", msg)
		return nil
	}

	// create the template
	body := fmt.Sprintf(`From: %s
To: %s
Subject: [UPMAIL]: %s %s
%s

Time: %s
`, n.Sender, n.Recipient, result.Title, result.Status(), result.String(), time.Now().Format(time.UnixDate))

	// send the email
	if err := smtp.SendMail(n.Server, n.Auth, n.Sender, []string{n.Recipient}, []byte(body)); err != nil {
		return fmt.Errorf("Send mail failed: %v", err)
	}

	return nil
}
