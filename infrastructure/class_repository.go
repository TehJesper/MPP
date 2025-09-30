package infrastructure

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
    Desc   string `json:"desc"`
    Choose int    `json:"choose"`
    Type   string `json:"type"`
    From   From   `json:"from"`
}

type From struct {
    OptionSetType string   `json:"option_set_type"`
    Options       []Option `json:"options"`
}

type Option struct {
    OptionType string `json:"option_type"`
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

func OpenClassFile() ([]byte, error) {
	return os.ReadFile("classes.json")
}

func LoadClassesAndSubclasses() ([]string, error) {
	data, err := OpenClassFile()

	if err != nil {
		return nil, err
	}

    var classes []Class
    if err := json.Unmarshal(data, &classes); err != nil {
        return nil, err
    }

    var all []string
    for _, r := range classes {
        all = append(all, r.Name)
        for _, sub := range r.Subclasses {
            all = append(all, sub.Name)
        }
    }
	
    return all, nil
}

func LoadClassSkillsFromName(name string) ([]string, error) {
    data, err := OpenClassFile()
    if err != nil {
        return nil, err
    }

    var classes []Class
    if err := json.Unmarshal(data, &classes); err != nil {
        return nil, err
    }

    nameLower := strings.ToLower(name)


    for _, c := range classes {
        if strings.ToLower(c.Name) == nameLower {
            return extractSkills(c), nil
        }
        for _, sc := range c.Subclasses {
            if strings.ToLower(sc.Name) == nameLower {
                return extractSkills(c), nil
            }
        }
    }

    return nil, nil
}

// Helper to extract all skill names from a class
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