package services

import (
	"fmt"
	"math"
	"strings"
	"test/domain"
	rce "test/infrastructure/race"
	cls "test/infrastructure/class"
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

	// Proficiency bonus
	bonus := 0
		switch {
		case level >= 17:
			bonus = 6
		case level >= 13:
			bonus = 5
		case level >= 9:
			bonus = 4
		case level >= 5:
			bonus = 3
		default:
			bonus = 2
		}
	

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

	char := domain.NewCharacter(name, race, class, level, bonus, abilities, modifiers, skills)

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
Proficiency bonus: %+d
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
		c.ProficiencyBonus,
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

canCast, err := cls.CanCastSpells(c.Class)
if err != nil {

}

if canCast {
	spellcasting, err := cls.GetSpellcastingForClassAndLevel(c.Class, c.Level)
	if err == nil && spellcasting != nil {
		spellSlotOutput := "\nSpell slots:"

		if spellcasting.SpellSlotsLevel0 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 0: %d", spellcasting.SpellSlotsLevel0)
		}
		if spellcasting.SpellSlotsLevel1 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 1: %d", spellcasting.SpellSlotsLevel1)
		}
		if spellcasting.SpellSlotsLevel2 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 2: %d", spellcasting.SpellSlotsLevel2)
		}
		if spellcasting.SpellSlotsLevel3 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 3: %d", spellcasting.SpellSlotsLevel3)
		}
		if spellcasting.SpellSlotsLevel4 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 4: %d", spellcasting.SpellSlotsLevel4)
		}
		if spellcasting.SpellSlotsLevel5 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 5: %d", spellcasting.SpellSlotsLevel5)
		}
		if spellcasting.SpellSlotsLevel6 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 6: %d", spellcasting.SpellSlotsLevel6)
		}
		if spellcasting.SpellSlotsLevel7 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 7: %d", spellcasting.SpellSlotsLevel7)
		}
		if spellcasting.SpellSlotsLevel8 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 8: %d", spellcasting.SpellSlotsLevel8)
		}
		if spellcasting.SpellSlotsLevel9 > 0 {
			spellSlotOutput += fmt.Sprintf("\n  Level 9: %d", spellcasting.SpellSlotsLevel9)
		}

		base += spellSlotOutput
	}
}
	initiative, ac, passivePerception := CalculateDerivedStats(c)
	base += fmt.Sprintf("\nArmor class: %d\nInitiative bonus: %d\nPassive perception: %d", ac, initiative, passivePerception)
	return base
}

func calculateStat(score int) int {
	mod := int(math.Floor(float64(score-10) / 2))
	return mod
}

func CalculateDerivedStats(c domain.Character) (initiative int, ac int, passivePerception int) {
	// 1. Initiative = DEX modifier
	initiative = c.AbilityModifiers.Dexterity

	// 2. Armor Class (AC)
	baseAC := 10
	if c.Equipment.Armor != "" {
		switch strings.ToLower(c.Equipment.Armor) {
		case "leather armor":
			baseAC = 11
			ac = baseAC + c.AbilityModifiers.Dexterity
		case "studded leather":
			baseAC = 12
			ac = baseAC + c.AbilityModifiers.Dexterity
		case "chain shirt":
			baseAC = 13
			// Medium armor: max DEX modifier = 2
			mod := c.AbilityModifiers.Dexterity
			if mod > 2 {
				mod = 2
			}
			ac = baseAC + mod
		case "scale mail":
			baseAC = 14
			mod := c.AbilityModifiers.Dexterity
			if mod > 2 {
				mod = 2
			}
			ac = baseAC + mod
		case "plate":
			baseAC = 18
			ac = baseAC // Heavy armor ignores DEX
		default:
			ac = baseAC + c.AbilityModifiers.Dexterity
		}
	} else {
		// No armor
		ac = 10 + c.AbilityModifiers.Dexterity
	}

	// Shield adds +2 AC
	if c.Equipment.Shield != "" {
		ac += 2
	}

	// 3. Passive Perception = 10 + WIS modifier + proficiency if proficient
	
	passivePerception = 10 + c.AbilityModifiers.Wisdom
	skills := strings.Split(c.Skills, ",") // assuming c.Skills is "acrobatics,athletics,perception"
	for _, skill := range skills {
		skill = strings.TrimSpace(skill) // remove any extra spaces
		if strings.ToLower(skill) == "perception" {
			passivePerception += c.ProficiencyBonus
			break
		}
	}

	return initiative, ac, passivePerception
}