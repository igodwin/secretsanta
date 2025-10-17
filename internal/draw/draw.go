package draw

import (
	"math/rand"

	"github.com/igodwin/secretsanta/pkg/participant"
)

const maxRetries = 1000

func shuffleParticipants(participants []*participant.Participant) []*participant.Participant {
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})
	return participants
}

func Names(participants []*participant.Participant) ([]*participant.Participant, error) {
	participants = shuffleParticipants(participants)

	for i := 0; i < maxRetries; i++ {
		recipients := shuffleParticipants(participants)
		usedRecipients := make([]bool, len(recipients))
		usedCount := 0

		// TODO: refactor so no longer O(n^2)
		for _, participant := range participants {
			matched := false
			for j := 0; j < len(recipients); j++ {
				possibleRecipient := recipients[j]
				if usedRecipients[j] {
					continue
				}

				if err := participant.UpdateRecipient(possibleRecipient); err == nil {
					usedRecipients[j] = true
					usedCount++
					matched = true
					break
				}
			}
			if !matched {
				break
			}
		}

		if usedCount == len(recipients) {
			break
		}
	}

	return participants, nil
}