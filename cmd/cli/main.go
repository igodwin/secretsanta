package main

import (
	"encoding/json"
	"fmt"
	. "github.com/igodwin/secretsanta/internal/participant"
	"math/rand"
	"os"
	"strings"
)

func shuffleParticipantSlice(participants []Participant) []Participant {
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})
	return participants
}

func drawNames(participants []Participant) error {

	participants = shuffleParticipantSlice(participants)
	recipients := make([]Participant, len(participants))
	for i, p := range participants {
		recipients[i] = p
	}
	recipients = shuffleParticipantSlice(recipients)

	for _, i := range participants {
		fmt.Println(fmt.Sprintf("Found individual named %s with the email addresses %s. Participant also has the following exclusion(s): %s", i.Name, strings.Join(i.Email, ", "), strings.Join(i.Exclusions, ", ")))
	}

	return nil
}

func main() {
	configPath := "participants.json"
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		configPath = envConfigPath
	}
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	fmt.Printf("Using config file: %s\n", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var participants []Participant
	err = json.Unmarshal(data, &participants)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = drawNames(participants)
	if err != nil {
		fmt.Println(err)
		return
	}
}
