package spells

import (
	"os"
	"fmt"
	"strings"
	"encoding/json"
)

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Spells []Spell `json:"spells"`
}

type Spell struct {
	Name    string    `json:"name"`
	Index   string    `json:"index"`
	Classes []Class   `json:"classes"`
}

type Class struct {
	Name          string          `json:"name"`
	Index         string          `json:"index"`
	Spellcasting  *Spellcasting   `json:"spellcasting,omitempty"`
	Spells        []ClassSpell    `json:"spells,omitempty"`
}

type Spellcasting struct {
	Info               []SpellInfo         `json:"info"`
	SpellcastingAbility SpellcastingAbility `json:"spellcasting_ability"`
}

type SpellInfo struct {
	Name string `json:"name"`
}

type SpellcastingAbility struct {
	Name string `json:"name"`
}

type ClassSpell struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type SpellsWrapper struct {
	Spells []Spell `json:"spells"`
}

func OpenSpellsFile() ([]byte, error) {
	return os.ReadFile("spells.json")
}

func GetSpellLevel(spellName string, className string) (int, error) {
	data, err := OpenSpellsFile()
	if err != nil {
		return -1, fmt.Errorf("failed to read spells.json: %w", err)
	}

	var wrapper SpellsWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return -1, fmt.Errorf("failed to parse JSON: %w", err)
	}

	spellName = strings.ToLower(spellName)
	className = strings.ToLower(className)
	
	for _, spell := range wrapper.Spells {
		for _, class := range spell.Classes {
			if strings.ToLower(class.Name) == className {
				for _, cs := range class.Spells {
					if strings.ToLower(cs.Name) == spellName {
						return cs.Level, nil
					}
				}
			}
		}
	}

	return -1, fmt.Errorf("spell %q not found for class %q", spellName, className)
}