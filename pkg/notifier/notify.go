package notifier

import "github.com/igodwin/secretsanta/pkg/participant"

type Notifier interface {
	SendNotification(participant *participant.Participant) error
}
