package services

import (
	"fmt"
	"strings"
	skl "test/infrastructure/class"
)

func getSkills(class string) string {
	skillsSlice, amount, err := skl.LoadClassSkillsFromName(class)
	if err != nil {
		fmt.Print(err)
	}

	if len(skillsSlice) > amount {
		skillsSlice = skillsSlice[:amount]
	}
	// Default acolyte skills
	skillsSlice = append(skillsSlice, "Insight", "Religion")

	for i, s := range skillsSlice {
		skillsSlice[i] = strings.ToLower(s)
	}

	skills := strings.Join(skillsSlice, ", ")

	return skills
}
