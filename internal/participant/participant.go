package participant

import "fmt"

type Participant struct {
	Name       string   `json:"name"`
	Email      []string `json:"email"`
	Exclusions []string `json:"exclusions"`
	Match      *Participant
}

func (i *Participant) UpdateMatch(participant *Participant) error {
	for _, a := range i.Exclusions {
		if participant.Name == a {
			return fmt.Errorf("participant %s is excluded", participant.Name)
		}
	}
	i.Match = participant

	return nil
}
