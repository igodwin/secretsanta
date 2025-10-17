package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/igodwin/secretsanta/internal/draw"
	"github.com/igodwin/secretsanta/internal/notification"
	"github.com/igodwin/secretsanta/pkg/config"
	"github.com/igodwin/secretsanta/pkg/participant"
)

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

	var participants []*participant.Participant
	err = json.Unmarshal(data, &participants)
	if err != nil {
		log.Fatalf("error encountered while loading participants: \n%v", err)
		return
	}

	participants, err = draw.Names(participants)
	if err != nil {
		log.Fatalf("error encountered while drawing names: \n%v", err)
		return
	}

	err = notification.Send(participants, appConfig)
	if err != nil {
		log.Fatalf("error encountered while attempting to send notifications:\n%v", err)
		return
	}
}

