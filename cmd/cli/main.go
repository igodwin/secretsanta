package main

import (
	"encoding/json"
	"fmt"
	"github.com/igodwin/secretsanta/internal/individual"
	"math/rand"
	"os"
	"strings"
)

func main() {
	configPath := "individuals.json"
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

	var individuals []individual.Individual
	err = json.Unmarshal(data, &individuals)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Shuffle the slice
	rand.Shuffle(len(individuals), func(i, j int) {
		individuals[i], individuals[j] = individuals[j], individuals[i]
	})

	for _, i := range individuals {
		fmt.Println(fmt.Sprintf("Found individual named %s with the email addresses %s. Individual also has the following exclusion(s): %s", i.Name, strings.Join(i.Email, ", "), strings.Join(i.Exclusions, ", ")))
	}
}
