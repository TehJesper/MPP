package domain

type Character struct {
	ID                int
	Name              string
	Race              string
	Class             string
	Level             int
	ProficiencyBonus  int
	AbilityScore      AbilityScores
	AbilityModifiers  AbilityModifiers
	Equipment         Equipment
	Skills            string
	ArmorClass        int
	Initiative        int
	PassivePerception int
}
type AbilityScores struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

type AbilityModifiers struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

type Equipment struct {
	Mainhand string
	Offhand  string
	Shield   string
	Armor    string
}

func NewCharacter(
	name, race, class string,
	level, proficiencybonus int,
	abilities AbilityScores,
	modifiers AbilityModifiers,
	skills string,
	initiative, passivePerception int,
) Character {
	return Character{
		Name:              name,
		Race:              race,
		Class:             class,
		Level:             level,
		ProficiencyBonus:  proficiencybonus,
		AbilityScore:      abilities,
		AbilityModifiers:  modifiers,
		Skills:            skills,
		Initiative:        initiative,
		PassivePerception: passivePerception,
	}
}

func (c *Character) EquipWeapon(slot, weapon string) {
	switch slot {
	case "main hand", "mainhand":
		c.Equipment.Mainhand = weapon
	case "off hand", "offhand":
		c.Equipment.Offhand = weapon
	}
}

func (c *Character) EquipArmor(armor string) {
	c.Equipment.Armor = armor
}

func (c *Character) EquipShield(shield string) {
	c.Equipment.Shield = shield
}
