package individual

import "fmt"

type Individual struct {
	Name       string   `json:"name"`
	Email      []string `json:"email"`
	Exclusions []string `json:"exclusions"`
	Match      *Individual
}

func (i *Individual) SetMatch(individual *Individual) error {
	for _, a := range i.Exclusions {
		if individual.Name == a {
			return fmt.Errorf("individual %s is excluded", individual.Name)
		}
	}
	i.Match = individual

	return nil
}
