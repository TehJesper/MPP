package services

import (
	"fmt"
	"test/domain"
	rce "test/infrastructure/race"
)

func checkRace(race string) error {
	races, err := rce.LoadRacesAndSubraces()

	if err != nil {
		return err
	}

	if !domain.IsValidRace(race, races) {
		return fmt.Errorf("invalid race: %s", race)
	}

	return nil
}
