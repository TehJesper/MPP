package services

import (
	"fmt"
	"math"
	"strings"
	"test/domain"
	rce "test/infrastructure/race"
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

	if err := checkRace(race); err != nil {
		return domain.Character{}, err
	}

	if err := checkClass(class); err != nil {
		return domain.Character{}, err
	}

	skills := getSkills(class)

	abilities := domain.AbilityScores{
		Strength:     str,
		Dexterity:    dex,
		Constitution: con,
		Intelligence: intel,
		Wisdom:       wis,
		Charisma:     cha,
	}

	bonuses, err := rce.GetRaceBonusesByName(race)
	if err != nil {
		return domain.Character{}, err
	}

	for ability, bonus := range bonuses {
		switch strings.ToLower(ability) {
		case "str":
			abilities.Strength += bonus
		case "dex":
			abilities.Dexterity += bonus
		case "con":
			abilities.Constitution += bonus
		case "int":
			abilities.Intelligence += bonus
		case "wis":
			abilities.Wisdom += bonus
		case "cha":
			abilities.Charisma += bonus
		}
	}

	modifiers := domain.AbilityModifiers{
		Strength:     calculateStat(abilities.Strength),
		Dexterity:    calculateStat(abilities.Dexterity),
		Constitution: calculateStat(abilities.Constitution),
		Intelligence: calculateStat(abilities.Intelligence),
		Wisdom:       calculateStat(abilities.Wisdom),
		Charisma:     calculateStat(abilities.Charisma),
	}

	char := domain.NewCharacter(name, race, class, level, abilities, modifiers, skills)

	return s.repo.Save(char)
}

func (s *CharacterService) ViewCharacter(name string) string {
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

func formatCharacter(c domain.Character) string {
	base := fmt.Sprintf(
		`Name: %s
Class: %s
Race: %s
Background: acolyte
Level: %d
Ability scores:
  STR: %d (%+d)
  DEX: %d (%+d)
  CON: %d (%+d)
  INT: %d (%+d)
  WIS: %d (%+d)
  CHA: %d (%+d)
Proficiency bonus: +2
Skill proficiencies: %s`,
		c.Name,
		strings.ToLower(c.Class),
		strings.ToLower(c.Race),
		c.Level,
		c.AbilityScore.Strength, c.AbilityModifiers.Strength,
		c.AbilityScore.Dexterity, c.AbilityModifiers.Dexterity,
		c.AbilityScore.Constitution, c.AbilityModifiers.Constitution,
		c.AbilityScore.Intelligence, c.AbilityModifiers.Intelligence,
		c.AbilityScore.Wisdom, c.AbilityModifiers.Wisdom,
		c.AbilityScore.Charisma, c.AbilityModifiers.Charisma,
		c.Skills,
	)
	if c.Equipment.Mainhand != "" {
		base += fmt.Sprintf("\nMain hand: %s", c.Equipment.Mainhand)
	}
	if c.Equipment.Offhand != "" {
		base += fmt.Sprintf("\nOff hand: %s", c.Equipment.Offhand)
	}
	if c.Equipment.Armor != "" {
		base += fmt.Sprintf("\nArmor: %s", c.Equipment.Armor)
	}
	if c.Equipment.Shield != "" {
		base += fmt.Sprintf("\nShield: %s", c.Equipment.Shield)
	}

	return base
}

func calculateStat(score int) int {
	mod := int(math.Floor(float64(score-10) / 2))
	return mod
}
