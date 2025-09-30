package services

import (
	"test/domain"
	"test/infrastructure"
	"fmt"
	"log"
	"math"
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

	// Check the skills
	skills, err := infrastructure.LoadClassSkillsFromName(class)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(skills)



	char := domain.NewCharacter(name, race, class, level, str, dex, con, intel, wis, cha)

	return s.repo.Save(char)
}

func (s *CharacterService) ViewCharacter(name string) (string, error) {
	c, err := s.repo.View(name)

	if err != nil {
		return "", fmt.Errorf("character %s not found", name)
	}

	return formatCharacter(c)
}

func (s *CharacterService) DeleteCharacter(name string) {
	s.repo.Delete(name)
}

func formatCharacter(c domain.Character) (string, error) {

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
		c.Name, c.Class, c.Race, c.Level,
		formatStat(c.Strength), formatStat(c.Dexterity), formatStat(c.Constitution),
		formatStat(c.Intelligence), formatStat(c.Wisdom), formatStat(c.Charisma),
		c.Skills,
	), nil
}

func formatStat(score int) string {
	mod := int(math.Floor(float64(score-10) / 2))
	sign := "+"
	if mod < 0 {
		sign = "" 
	}
	return fmt.Sprintf("%d (%s%d)", score, sign, mod)
}
