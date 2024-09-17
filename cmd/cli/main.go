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
	if len(os.Args) < 2 {
		fmt.Println("Usage: secretsanta <filename>")
		return
	}
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
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
