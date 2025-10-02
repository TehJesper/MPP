package services

import (
	"fmt"
	"strings"
	"test/infrastructure"
)

func getSkills(class string) (string) {
	skillsSlice, err := infrastructure.LoadClassSkillsFromName(class)
	if err != nil {
		fmt.Print(err)
	}

	if (len(skillsSlice) > 2){
		skillsSlice = skillsSlice[:2]
	}
	// Default acolyte skills
	skillsSlice  = append(skillsSlice, "Insight", "Religion")
	
	for i, s := range skillsSlice {
		skillsSlice[i] = strings.ToLower(s)
	}
	
	skills := strings.Join(skillsSlice, ", ")

	return skills
}