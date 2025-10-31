package services

import (
	"fmt"
	"slices"
	"strings"
	cls "test/infrastructure/class"
	"test/infrastructure/spells"
)

func (s *CharacterService) PrepareSpell(name string, spell string) {
	c, err := s.repo.View(name)

	if err != nil {
		return
	}

	canCast, err := cls.CanCastSpells(c.Class)

	if err != nil {
		return
	}

	if !canCast {
		fmt.Print("this class can't cast spells")
		return
	}
	if !canPrepare(spell, c.Class, c.Level) {

		return
	}

	fmt.Print("Prepared spell ", spell)
}

func (s *CharacterService) LearnSpell(name string, spell string) {
	c, err := s.repo.View(name)
	if err != nil {
		return
	}

	canCast, err := cls.CanCastSpells(c.Class)
	if err != nil || !canCast {
		fmt.Print("this class can't cast spells")
		return
	}

	if !canLearn(spell, c.Class, c.Level) {
		return
	}

	fmt.Print("Learned spell ", spell)
}

// Check if class can learn
// Classes that can prepare spells. Source: https://www.thegamer.com/dungeons-dragons-dnd-preparing-spells-guide/
func canLearn(spell string, class string, level int) bool {
	classesLearn := []string{"sorcerer", "bard", "warlock"}
	if !slices.Contains(classesLearn, strings.ToLower(class)) {
		fmt.Print("this class prepares spells and can't learn them")
		return false
	}

	spellLevel, _ := spells.GetSpellLevel(spell, class)
	highestSlot, _ := cls.GetHighestSpellSlotForClassAndLevel(class, level)

	if spellLevel > highestSlot {
		fmt.Println("the spell has higher level than the available spell slots")
		return false
	}

	return true
}

// Check if class can prepare
// Classes that can prepare spells. Source: https://www.thegamer.com/dungeons-dragons-dnd-preparing-spells-guide/
func canPrepare(spell string, class string, level int) bool {
	classesPrepare := []string{"cleric", "druid", "paladin", "ranger", "wizard"}
	if !slices.Contains(classesPrepare, strings.ToLower(class)) {
		fmt.Print("this class learns spells and can't prepare them")
		return false
	}
	spellLevel, _ := spells.GetSpellLevel(spell, class)
	highestSlot, _ := cls.GetHighestSpellSlotForClassAndLevel(class, level)

	if spellLevel > highestSlot {
		fmt.Println("the spell has higher level than the available spell slots")
		return false
	}

	return true
}
