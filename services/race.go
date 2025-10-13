package services
import (
	"test/domain"
	rce "test/infrastructure/race"
	"fmt"
)

func checkRace(race string) (error) {
	races, err := rce.LoadRacesAndSubraces()

    if err != nil {
        return err
    }

    if !domain.IsValidRace(race, races) {
        return fmt.Errorf("invalid race: %s", race)
    }

	return nil
}
