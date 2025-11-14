package services

import (
	"fmt"
	"strings"
	skl "test/infrastructure/class"
	rce "test/infrastructure/race"
)

func getSkills(class string, race string) string {
	skillsSlice, amount, err := skl.LoadClassSkillsFromName(class)
	if err != nil {
		fmt.Print(err)
	}

	if len(skillsSlice) > amount {
		skillsSlice = skillsSlice[:amount]
	}

	raceSkills, err := rce.GetRaceSkillsByName(race)

	// Hardcode Dwarf history skill since it is not in the proficiency field of the API
	if strings.ToLower(race) == "dwarf" {
		raceSkills = append(raceSkills, "History")
	}

	skillsSlice = append(skillsSlice, raceSkills...)

	if err != nil {
		fmt.Print(err)
	}

	// Default acolyte skills
	skillsSlice = append(skillsSlice, "Insight", "Religion")

	for i, s := range skillsSlice {
		skillsSlice[i] = strings.ToLower(s)
	}

	skills := strings.Join(skillsSlice, ", ")

	return skills
}
