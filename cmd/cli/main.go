package main

import (
	"encoding/json"
	"github.com/igodwin/secretsanta/pkg/config"
	. "github.com/igodwin/secretsanta/pkg/notifier"
	. "github.com/igodwin/secretsanta/pkg/participant"
	"log"
	"math/rand"
	"os"
)

const maxRetries = 1000

func main() {
	appConfig := config.GetConfig()

	participantFilePath := "participants.json"
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		participantFilePath = envConfigPath
	}
	if len(os.Args) > 1 {
		participantFilePath = os.Args[1]
	}
	log.Printf("using participant file: %s\n", participantFilePath)
	data, err := os.ReadFile(participantFilePath)
	if err != nil {
		log.Fatalf("error encountered while reading participant file: \n%v", err)
		return
	}

	var participants []*Participant
	err = json.Unmarshal(data, &participants)
	if err != nil {
		log.Fatalf("error encountered while loading participants: \n%v", err)
		return
	}

	participants, err = drawNames(participants)
	if err != nil {
		log.Fatalf("error encountered while drawing names: \n%v", err)
		return
	}

	err = sendNotifications(participants, appConfig)
	if err != nil {
		log.Fatalf("error encountered while attempting to send notifications:\n%v", err)
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

func sendNotifications(participants []*Participant, appConfig *config.Config) error {
	var notifier Notifier
	var emailNotifier = &EmailNotifier{
		Host:        appConfig.SMTP.Host,
		Port:        appConfig.SMTP.Port,
		Identity:    appConfig.SMTP.Identity,
		Username:    appConfig.SMTP.Username,
		Password:    appConfig.SMTP.Password,
		FromAddress: appConfig.SMTP.FromAddress,
		FromName:    appConfig.SMTP.FromName,
	}

	for _, participant := range participants {
		switch participant.NotificationType {
		case "email":
			notifier = emailNotifier
		default:
			notifier = &Stdout{}
		}

		err := notifier.IsConfigured()
		if err != nil {
			return err
		}
		err = notifier.SendNotification(participant)
		if err != nil {
			return err
		}
	}
	return nil
}
