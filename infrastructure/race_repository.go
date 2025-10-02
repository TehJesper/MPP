package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
    "strings"
)

type Race struct {
    Index    string    `json:"index"`
    Name     string    `json:"name"`
    Subraces []Subrace `json:"subraces"`
    AbilityBonuses []AbilityBonus `json:"ability_bonuses"`
}

type AbilityBonus struct {
	AbilityScore struct {
		Index string `json:"index"`
	} `json:"ability_score"`
	Bonus int `json:"bonus"`
}

type Subrace struct {
    Index string `json:"index"`
    Name  string `json:"name"`
    AbilityBonuses []AbilityBonus `json:"ability_bonuses"`
}

func OpenRaceFile() ([]byte, error) {
	return os.ReadFile("races.json")
}

func LoadRacesAndSubraces() ([]string, error) {
    data, err := OpenRaceFile()
    if err != nil {
        return nil, err
    }

    var races []Race
    if err := json.Unmarshal(data, &races); err != nil {
        return nil, err
    }

    var all []string
    for _, r := range races {
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

	var races []Race
	if err := json.Unmarshal(data, &races); err != nil {
		return nil, err
	}

	flat := make(map[string]Race)
    for _, r := range races {
        flat[strings.ToLower(r.Name)] = r
        for _, sub := range r.Subraces {
            flat[strings.ToLower(sub.Name)] = Race{
                Index:          sub.Index,
                Name:           sub.Name,
                AbilityBonuses: sub.AbilityBonuses,
            }
        }
    }

	race, ok := flat[strings.ToLower(raceName)]
	if !ok {
		return nil, fmt.Errorf("race not found: %s", raceName)
	}

	bonuses := make(map[string]int)
	for _, b := range race.AbilityBonuses {
		key := strings.ToLower(b.AbilityScore.Index)
		bonuses[key] = b.Bonus
	}

	return bonuses, nil
}