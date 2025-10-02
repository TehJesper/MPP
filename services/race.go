package services
import (
	"test/domain"
	"test/infrastructure"
	"fmt"
)

func checkRace(race string) (error) {
	races, err := infrastructure.LoadRacesAndSubraces()

    if err != nil {
        return err
    }

    if !domain.IsValidRace(race, races) {
        return fmt.Errorf("invalid race: %s", race)
    }

	return nil
}
