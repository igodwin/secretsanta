package notifier

import (
	"fmt"
	"github.com/igodwin/secretsanta/pkg/participant"
)

const stdoutAssignmentTemplate = `%s has %s
`

type Stdout struct {
}

func (s *Stdout) SendNotification(participant *participant.Participant) error {
	fmt.Printf(stdoutAssignmentTemplate, participant.Name, participant.Recipient.Name)
	return nil
}

func (s *Stdout) IsConfigured() error {
	return nil
}
