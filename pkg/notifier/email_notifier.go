package notifier

import (
	"fmt"
	"github.com/igodwin/secretsanta/pkg/participant"
	"net/smtp"
	"strings"
)

const (
	subjectSuffix           = "'s Secret Santa Assignment"
	emailAssignmentTemplate = `Hello %s,

You have been given the important task of finding the perfect gift for %s this year! You are the only person who knows this, so you should try to keep it a surprise. Also don't delete this email too soon, because I'm not going to remember who you have.

Think about some things your unknown gifter should know you would like for Christmas, and send an email to the rest of the group so that your Secret Santa will see it.

Merry Christmas!

Papa Elf`
)

type EmailNotifier struct {
	Host        string
	Port        string
	Identity    string
	Username    string
	Password    string
	FromAddress string
	FromName    string
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
		fmt.Sprintf(emailAssignmentTemplate, participant.Name, participant.Recipient.Name)))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", e.Host, e.Port), auth, e.FromAddress, append(participant.ContactInfo, e.FromAddress), formattedMessage)
	if err != nil {
		return err
	}
	return nil
}
