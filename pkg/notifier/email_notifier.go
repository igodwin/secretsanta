package notifier

import (
	"fmt"
	"github.com/igodwin/secretsanta/pkg/participant"
	"net/smtp"
	"strings"
)

const (
	subjectSuffix     = "'s Secret Santa Assignment"
	emailBodyTemplate = `Hello %s,

You have been given the important task of finding the perfect gift for %s this year! You are the only person who knows this, so you should try to keep it a surprise. Also don't delete this email too soon, because I'm not going to remember who you have.

Think about some things your unknown gifter should know you would like for Christmas, and send an email to the rest of the group so that your Secret Santa will see it.

Merry Christmas!

Papa Elf`
)

type SendMailFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

type EmailNotifier struct {
	Host         string
	Port         string
	Identity     string
	Username     string
	Password     string
	FromAddress  string
	FromName     string
	SendMailFunc SendMailFunc
}

func (e *EmailNotifier) SendNotification(participant *participant.Participant) error {
	allAddresses := strings.Join(participant.ContactInfo, ",")
	auth := smtp.PlainAuth(e.Identity, e.Username, e.Password, e.Host)
	from := fmt.Sprintf("<%s>", e.FromAddress)
	if e.FromName != "" {
		from = fmt.Sprintf(`"%s" <%s>`, e.FromName, e.FromAddress)
	}
	formattedMessage := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s%s\r\n\r\n%s",
		from,
		allAddresses,
		participant.Name,
		subjectSuffix,
		fmt.Sprintf(emailBodyTemplate, participant.Name, participant.Recipient.Name)))

	if e.SendMailFunc == nil {
		e.SendMailFunc = smtp.SendMail
	}

	err := e.SendMailFunc(fmt.Sprintf("%s:%s", e.Host, e.Port), auth, e.FromAddress, append(participant.ContactInfo, e.FromAddress), formattedMessage)
	if err != nil {
		return err
	}
	return nil
}

func (e *EmailNotifier) IsConfigured() error {
	if e.Host == "" && e.Port == "" && e.Username == "" && e.Password == "" && e.FromAddress == "" {
		return fmt.Errorf("smtp is not configured")
	}
	return nil
}
