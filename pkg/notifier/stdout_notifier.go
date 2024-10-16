package notifier

import (
	"fmt"
	"github.com/igodwin/secretsanta/pkg/participant"
)

const stdoutAssignmentTemplate = `%s has %s
`

type StdOut struct {
}

func (s *StdOut) SendNotification(participant *participant.Participant) error {
	fmt.Printf(stdoutAssignmentTemplate, participant.Name, participant.Recipient.Name)
	return nil
}

func (s *StdOut) IsConfigured() error {
	return nil
}
