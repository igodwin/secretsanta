package main

import (
	"encoding/json"
	"fmt"
	. "github.com/igodwin/secretsanta/internal/participant"
	"log"
	"math/rand"
	"os"
)

const maxRetries = 1000

func main() {
	configPath := "participants.json"
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		configPath = envConfigPath
	}
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	log.Printf("using config file: %s\n", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var participants []*Participant
	err = json.Unmarshal(data, &participants)
	if err != nil {
		fmt.Println(err)
		return
	}

	participants, err = drawNames(participants)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = notify(participants)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func shuffleParticipants(participants []*Participant) []*Participant {
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})
	return participants
}

func drawNames(participants []*Participant) ([]*Participant, error) {
	participants = shuffleParticipants(participants)

	for i := 0; i < maxRetries; i++ {
		recipients := shuffleParticipants(participants)
		usedRecipients := make([]bool, len(recipients))
		usedCount := 0

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

func notify(participants []*Participant) error {
	for _, participant := range participants {
		fmt.Printf("%s is buying for %s\n", participant.Name, participant.Recipient.Name)
	}
	return nil
}
