package class

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Class struct {
	Index              string              `json:"index"`
	Name               string              `json:"name"`
	Subclasses         []Subclass          `json:"subclasses"`
	ProficiencyChoices []ProficiencyChoice `json:"proficiency_choices"`
	Spells 			   []Spells			   `json:"spells"`
	ClassLevel 		   []ClassLevel		   `json:"class_levels"`
}

type ProficiencyChoice struct {
	Choose int    `json:"choose"`
	From   From   `json:"from"`
}

type From struct {
	Options       []Option `json:"options"`
}

type Option struct {
	Item       Item   `json:"item"`
}

type Item struct {
	Index string `json:"index"`
	Name  string `json:"name"` 	
}

type Subclass struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type Spells struct {
	Name string `json:"name"`
	Level int `json:"level"` 
}

type ClassLevel struct {
	Level        int           `json:"level,omitempty"`
	Spellcasting *Spellcasting `json:"spellcasting,omitempty"`
}

type Spellcasting struct {
	SpellSlotsLevel0 int `json:"cantrips_known"`
	SpellSlotsLevel1 int `json:"spell_slots_level_1"`
	SpellSlotsLevel2 int `json:"spell_slots_level_2"`
	SpellSlotsLevel3 int `json:"spell_slots_level_3"`
	SpellSlotsLevel4 int `json:"spell_slots_level_4"`
	SpellSlotsLevel5 int `json:"spell_slots_level_5"`
	SpellSlotsLevel6 int `json:"spell_slots_level_6"`
	SpellSlotsLevel7 int `json:"spell_slots_level_7"`
	SpellSlotsLevel8 int `json:"spell_slots_level_8"`
	SpellSlotsLevel9 int `json:"spell_slots_level_9"`
}

type ClassesWrapper struct {
	Classes []Class `json:"classes"`
}

func OpenClassFile() ([]byte, error) {
	return os.ReadFile("classes.json")
}

func LoadClassesAndSubclasses() ([]string, error) {
	data, err := OpenClassFile()
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Classes []Class `json:"classes"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	var all []string
	for _, r := range wrapper.Classes {
		all = append(all, r.Name)
		for _, sub := range r.Subclasses {
			all = append(all, sub.Name)
		}
	}

	return all, nil
}

func LoadClassSkillsFromName(name string) ([]string, int, error) {
	data, err := OpenClassFile()
	if err != nil {
		return nil, 0, err
	}

	var classes ClassesWrapper
	if err := json.Unmarshal(data, &classes); err != nil {
		return nil, 0, err
	}

	nameLower := strings.ToLower(name)

	for _, c := range classes.Classes {
		if strings.ToLower(c.Name) == nameLower {
			return extractSkills(c), c.ProficiencyChoices[0].Choose, nil
		}
		for _, sc := range c.Subclasses {
			if strings.ToLower(sc.Name) == nameLower {
				return extractSkills(c), c.ProficiencyChoices[0].Choose, nil
			}
		}
	}

	return nil, 0, nil
}

func extractSkills(c Class) []string {
	var skills []string
	for _, pc := range c.ProficiencyChoices {
		for _, opt := range pc.From.Options {
			name := strings.TrimPrefix(opt.Item.Name, "Skill: ")
			skills = append(skills, name)
		}
	}
	return skills
}

func CanCastSpells(name string) (bool, error) {
	data, err := OpenClassFile()
	if err != nil {
		return false, err
	}

	var wrapper ClassesWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return false, err
	}

	nameLower := strings.ToLower(name)

	for _, c := range wrapper.Classes {
		if strings.ToLower(c.Name) == nameLower || strings.ToLower(c.Index) == nameLower {
			if len(c.Spells) > 0 {
				return true, nil
			}
			return false, nil
		}
		for _, sc := range c.Subclasses {
			if strings.ToLower(sc.Name) == nameLower || strings.ToLower(sc.Index) == nameLower {

				if len(c.Spells) > 0 {
					return true, nil
				}
				return false, nil
			}
		}
	}

	return false, nil
}
func GetSpellcastingForClassAndLevel(name string, level int) (*Spellcasting, error) {
	data, err := OpenClassFile()
	if err != nil {
		return nil, err
	}

	var wrapper ClassesWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	nameLower := strings.ToLower(name)

	for _, c := range wrapper.Classes {
		if strings.ToLower(c.Name) == nameLower || strings.ToLower(c.Index) == nameLower {
			if level <= 0 || level > len(c.ClassLevel) {
				fmt.Print("the spell has higher level than the available spell slots")
				return nil, nil
			}
			classLevel := c.ClassLevel[level-1]
			if classLevel.Spellcasting != nil {
				return classLevel.Spellcasting, nil
			}
			return nil, nil
		}
	}

	return nil, nil
}
func GetHighestSpellSlotForClassAndLevel(name string, level int) (int, error) {
	sc, err := GetSpellcastingForClassAndLevel(name, level)
	if err != nil {
		return 0, err
	}
	if sc == nil {
		return 0, nil
	}

	switch {
	case sc.SpellSlotsLevel9 > 0:
		return 9, nil
	case sc.SpellSlotsLevel8 > 0:
		return 8, nil
	case sc.SpellSlotsLevel7 > 0:
		return 7, nil
	case sc.SpellSlotsLevel6 > 0:
		return 6, nil
	case sc.SpellSlotsLevel5 > 0:
		return 5, nil
	case sc.SpellSlotsLevel4 > 0:
		return 4, nil
	case sc.SpellSlotsLevel3 > 0:
		return 3, nil
	case sc.SpellSlotsLevel2 > 0:
		return 2, nil
	case sc.SpellSlotsLevel1 > 0:
		return 1, nil
	case sc.SpellSlotsLevel0 > 0:
		return 0, nil
	default:
		return 0, nil
	}
}