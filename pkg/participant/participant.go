package participant

import "fmt"

type Participant struct {
	Name             string   `json:"name"`
	NotificationType string   `json:"notification_type"`
	ContactInfo      []string `json:"contact_info"`
	Exclusions       []string `json:"exclusions"`
	Recipient        *Participant
}

func (p *Participant) UpdateRecipient(participant *Participant) error {
	if p.Name == participant.Name {
		return fmt.Errorf("cannot update match with self")
	}
	for _, a := range p.Exclusions {
		if participant.Name == a {
			return fmt.Errorf("participant %s is excluded", participant.Name)
		}
	}
	p.Recipient = participant

	return nil
}
