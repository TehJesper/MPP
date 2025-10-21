package class

import (
	"encoding/json"
	"os"
	"strings"
)

type Class struct {
	Index              string              `json:"index"`
	Name               string              `json:"name"`
	Subclasses         []Subclass          `json:"subclasses"`
	ProficiencyChoices []ProficiencyChoice `json:"proficiency_choices"`
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
