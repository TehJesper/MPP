package character

import (
	"fmt"
	// "io"
	// "log"
	_"net/http"
	"os"
	"strings"
	"encoding/json"
	"math"
)

type Service struct {
	repo Repository
}

type AbilityBonus struct {
	AbilityScore struct {
		Name string `json:"name"`
	} `json:"ability_score"`
	Bonus int `json:"bonus"`
}

type Race struct {
	Index          string         `json:"index"`           
	Name           string         `json:"name"`            
	AbilityBonuses []AbilityBonus `json:"ability_bonuses"`
}
type Class struct {
	Index              string               `json:"index"`
	Name               string               `json:"name"`
	ProficiencyChoices []ProficiencyChoice  `json:"proficiency_choices"`
}

type ProficiencyChoice struct {
	Choose int     `json:"choose"`
	From   Options `json:"from"`
}

type Options struct {
	Options []Option `json:"options"`
}

type Option struct {
	Item Item `json:"item"`
}

type Item struct {
	Name string `json:"name"`
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateNewCharacter(name string, class string, race string, level int, str int, dex int, con int, intl int, wis int, cha int) (*Character, error) {
	// Fetch racial ability score bonuses
	bonuses, err := checkRaceBonus(strings.ReplaceAll(race, " ", "-"))
	if err != nil {
		return nil, err
	}

	// Load all classes
	classFile, err := os.Open("classes.json")
	if err != nil {
		return nil, err
	}
	defer classFile.Close()

	var allClasses []Class
	if err := json.NewDecoder(classFile).Decode(&allClasses); err != nil {
		return nil, err
	}
	classData, err := findClassByName(allClasses, class)
	if err != nil {
		return nil, err
	}
	skills := getSkillsForClass(*classData)

	// Apply bonuses
	for _, ab := range bonuses {
		switch ab.AbilityScore.Name {
		case "STR":
			str += ab.Bonus
		case "DEX":
			dex += ab.Bonus
		case "CON":
			con += ab.Bonus
		case "INT":
			intl += ab.Bonus
		case "WIS":
			wis += ab.Bonus
		case "CHA":
			cha += ab.Bonus
		}
	}
	char := &Character{
		Name: name,
		Class: strings.ToLower(class),
		Race: strings.ToLower(race),
		Level: level,
		Strength: str,
		Dexterity: dex,
		Constitution: con,
		Intelligence: intl,
		Wisdom: wis,
		Charisma: cha,
		Skills: skills,

	}
	if err := s.repo.Create(char); err != nil {
		return nil, err
	}
	// fmt.Println("Character created!")
	return char, nil
}

func checkRaceBonus(race string) ([]AbilityBonus, error) {
	file, err := os.Open("races.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allRaces []Race
	if err := json.NewDecoder(file).Decode(&allRaces); err != nil {
		return nil, err
	}

	// Find the race matching the parameter
	for _, r := range allRaces {
    if strings.EqualFold(r.Index, race) {
        return r.AbilityBonuses, nil
    }
}

	return nil, fmt.Errorf("race not found: %s", race)
}

func findClassByName(allClasses []Class, className string) (*Class, error) {
	
	for _, c := range allClasses {
		if strings.EqualFold(c.Index, className) {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("class not found: %s", className)
}
func getSkillsForClass(classData Class) string {
    var skills []string

    // Add the first N class skills (based on Choose)
    if len(classData.ProficiencyChoices) > 0 {
        options := classData.ProficiencyChoices[0].From.Options
        choose := classData.ProficiencyChoices[0].Choose
        if choose > len(options) {
            choose = len(options)
        }

        for i := 0; i < choose; i++ {
            name := options[i].Item.Name
            // Strip the "Skill: " prefix if present
            name = strings.TrimPrefix(name, "Skill: ")
            // Convert to lowercase for consistency
            name = strings.ToLower(name)
            skills = append(skills, name)
        }
    }

    // Add Acolyte background skills last
    skills = append(skills, "insight", "religion")

    // Return as a comma-separated string
    return strings.Join(skills, ", ")
}



func (s *Service) ViewCharacterByName(name string) (string, error) {
	c, err := s.repo.View(name)
	if err != nil {
		return "", err
	}

	// helper function to format ability score with modifier
	formatStat := func(score int) string {
	mod := int(math.Floor(float64(score-10) / 2))
	sign := "+"
	if mod < 0 {
		sign = "" 
	}
	return fmt.Sprintf("%d (%s%d)", score, sign, mod)
}

	// Format skill proficiencies if stored as slice
	result := fmt.Sprintf(
		`Name: %s
Class: %s
Race: %s
Background: acolyte
Level: %d
Ability scores:
  STR: %s
  DEX: %s
  CON: %s
  INT: %s
  WIS: %s
  CHA: %s
Proficiency bonus: +2
Skill proficiencies: %s`,
		c.Name, c.Class, c.Race, c.Level,
		formatStat(c.Strength),
		formatStat(c.Dexterity),
		formatStat(c.Constitution),
		formatStat(c.Intelligence),
		formatStat(c.Wisdom),
		formatStat(c.Charisma),
		c.Skills,
	)

	return result, nil
}
func (s *Service) DeleteCharacterByName(name string) {
    s.repo.Delete(name)
    // if err != nil {
    //     return "", err
    // }
    // return fmt.Sprintf("Character %s deleted", name), nil
}