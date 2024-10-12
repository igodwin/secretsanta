package main

import (
	"encoding/json"
	"fmt"
	"github.com/igodwin/secretsanta/pkg/config"
	. "github.com/igodwin/secretsanta/pkg/notifier"
	. "github.com/igodwin/secretsanta/pkg/participant"
	"log"
	"math/rand"
	"os"
)

const maxRetries = 1000

func main() {
	config := config.GetConfig()

	participantFilePath := "participants.json"
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		participantFilePath = envConfigPath
	}
	if len(os.Args) > 1 {
		participantFilePath = os.Args[1]
	}
	log.Printf("using config file: %s\n", participantFilePath)
	data, err := os.ReadFile(participantFilePath)
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

	err = sendNotifications(participants, config)
	if err != nil {
		log.Fatalf("error encountered while attempting to send notifications:\n%v", err)
		return
	}
	log.Println("secret santa drawing complete")
}

func shuffleParticipants(participants []*Participant) []*Participant {
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})
	return participants
}

func drawNames(participants []*Participant) ([]*Participant, error) {
	log.Println("start drawing names")
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

func sendNotifications(participants []*Participant, config *config.Config) error {
	log.Println("start sending notifications")
	var notifier Notifier
	var emailNotifier *EmailNotifier
	if config.SMTPIsConfigured() {
		emailNotifier = &EmailNotifier{
			Host:        config.SMTP.Host,
			Port:        config.SMTP.Port,
			Identity:    config.SMTP.Identity,
			Username:    config.SMTP.Username,
			Password:    config.SMTP.Password,
			FromAddress: config.SMTP.FromAddress,
			FromName:    config.SMTP.FromName,
		}
	}

	for _, participant := range participants {
		switch participant.NotificationType {
		case "email":
			if !config.SMTPIsConfigured() {
				return fmt.Errorf("smtp is not configured, but a participant has an email notification set")
			}
			notifier = emailNotifier
		case "stdout":
			notifier = &StdOut{}
		default:
			return fmt.Errorf("unsupported notification type")
		}

		err := notifier.SendNotification(participant)
		if err != nil {
			return err
		}
	}
	return nil
}
