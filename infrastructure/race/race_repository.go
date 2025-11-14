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
	Traits 	   	   []Trait        `json:"traits"`
}

type AbilityBonus struct {
	AbilityScore struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"ability_score"`
	Bonus int `json:"bonus"`
}

type Subrace struct {
	Index          string         `json:"index"`
	Name           string         `json:"name"`
	AbilityBonuses []AbilityBonus `json:"ability_bonuses"`
}
type Trait struct {
	Index        string       `json:"index"`
	Name         string       `json:"name"`
	Proficiencies []Proficiency `json:"proficiencies"`
}

type Proficiency struct {
	Index string `json:"index"`
	Name  string `json:"name"`
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

func GetRaceSkillsByName(raceName string) ([]string, error) {
	data, err := OpenRaceFile()
	if err != nil {
		return nil, err
	}

	var Races RacesWrapper
	if err := json.Unmarshal(data, &Races); err != nil {
		return nil, err
	}
	var skills []string
	for _, r := range Races.Races {
		if strings.EqualFold(r.Name, raceName) {
			for _, t := range r.Traits {
				if len(t.Proficiencies) != 0 {
					for _, p := range t.Proficiencies {
						name := strings.TrimSpace(p.Name)
						if strings.HasPrefix(strings.ToLower(p.Index), "skill") {
							parts := strings.SplitN(name, ":", 2)
							if len(parts) == 2 {
								skill := strings.TrimSpace(parts[1])
									skills = append(skills, skill)
							}
						}
					}
				} 
			}
		}

	}
	return skills, nil
}