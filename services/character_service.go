package services

import (
	"test/domain"
	"test/infrastructure"
	"fmt"
	"math"
	"strings"
	)

type CharacterService struct {
	repo domain.CharacterRepository
}

func NewService(repo domain.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

func (s *CharacterService) CreateNewCharacter(
	name, class, race string,
	level, str, dex, con, intel, wis, cha int,
	) (domain.Character, error) {

	// Check for valid race
	if err := checkRace(race); err != nil {
        return domain.Character{}, err
    }

	// Check for valid class
	if err := checkClass(class); err != nil {
        return domain.Character{}, err
    }


	skills := getSkills(class)

	char := domain.NewCharacter(name, race, class, level, str, dex, con, intel, wis, cha, skills)

	bonuses, err := infrastructure.GetRaceBonusesByName(race)
	if err != nil {
		return domain.Character{}, err
	}

	for ability, bonus := range bonuses {
		switch strings.ToLower(ability) {
		case "str":
			char.Strength += bonus
		case "dex":
			char.Dexterity += bonus
		case "con":
			char.Constitution += bonus
		case "int":
			char.Intelligence += bonus
		case "wis":
			char.Wisdom += bonus
		case "cha":
			char.Charisma += bonus
		}
	}
	return s.repo.Save(char)
}

func (s *CharacterService) ViewCharacter(name string) (string) {
	c, err := s.repo.View(name)

	if err != nil {
		fmt.Printf(`character "%s" not found`, name)
		return ""
	}
	
	return formatCharacter(c)
}

func (s *CharacterService) DeleteCharacter(name string) {
	s.repo.Delete(name)
}

func formatCharacter(c domain.Character) (string) {

	return fmt.Sprintf(
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
		c.Name, strings.ToLower(c.Class), strings.ToLower(c.Race), c.Level,
		formatStat(c.Strength), formatStat(c.Dexterity), formatStat(c.Constitution),
		formatStat(c.Intelligence), formatStat(c.Wisdom), formatStat(c.Charisma),
		c.Skills,
	)
}

func formatStat(score int) string {
	mod := int(math.Floor(float64(score-10) / 2))
	sign := "+"
	if mod < 0 {
		sign = "" 
	}

	return fmt.Sprintf("%d (%s%d)", score, sign, mod)
}
