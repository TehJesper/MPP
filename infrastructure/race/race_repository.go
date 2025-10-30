package race

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Race struct {
	Index          string         `json:"index"`
	Name           string         `json:"name"`
	Subraces       []Subrace      `json:"subraces"`
	AbilityBonuses []AbilityBonus `json:"ability_bonuses"`
}

type AbilityBonus struct {
	AbilityScore struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"ability_score"`
	Bonus int `json:"bonus"`
}

type Subrace struct {
	Index string `json:"index"`
	Name  string `json:"name"`
    AbilityBonuses []AbilityBonus `json:"ability_bonuses"`
}

type RacesWrapper struct {
	Races []Race `json:"races"`
}

func OpenRaceFile() ([]byte, error) {
	return os.ReadFile("races.json")
}

func LoadRacesAndSubraces() ([]string, error) {
	data, err := OpenRaceFile()
	if err != nil {
		return nil, err
	}

	var races RacesWrapper
	if err := json.Unmarshal(data, &races); err != nil {
		return nil, err
	}

	var all []string
	for _, r := range races.Races {
		all = append(all, r.Name)
		for _, sub := range r.Subraces {
			all = append(all, sub.Name)
		}
	}
	return all, nil
}

func GetRaceBonusesByName(raceName string) (map[string]int, error) {
    data, err := OpenRaceFile()
    if err != nil {
        return nil, err
    }

    var races RacesWrapper
    if err := json.Unmarshal(data, &races); err != nil {
        return nil, err
    }

    normalizedInput := strings.ToLower(strings.ReplaceAll(raceName, " ", "-"))

    var foundRace *Race
    var parentRace *Race

    for _, r := range races.Races {
        normalizedR := strings.ToLower(strings.ReplaceAll(r.Name, " ", "-"))
        if normalizedR == normalizedInput {
            foundRace = &r
            break
        }

        for _, sub := range r.Subraces {
            normalizedSub := strings.ToLower(strings.ReplaceAll(sub.Name, " ", "-"))
            if normalizedSub == normalizedInput {
                foundRace = &Race{
                    Index:          sub.Index,
                    Name:           sub.Name,
                    AbilityBonuses: sub.AbilityBonuses,
                }
                parentRace = &r
                break
            }
        }

        if foundRace != nil {
            break
        }
    }

    if foundRace == nil {
        return nil, fmt.Errorf("race not found: %s", raceName)
    }

    bonuses := make(map[string]int)

    if parentRace != nil {
        for _, b := range parentRace.AbilityBonuses {
            key := strings.ToLower(b.AbilityScore.Index)
            bonuses[key] += b.Bonus
        }
    }

    for _, b := range foundRace.AbilityBonuses {
        key := strings.ToLower(b.AbilityScore.Index)
        bonuses[key] += b.Bonus
    }

    return bonuses, nil
}


