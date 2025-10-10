package domain

type Character struct {
	ID    int
	Name  string
	Race  string
	Class string
	Level int
	AbilityScore AbilityScores
	AbilityModifiers AbilityModifiers
	Skills string
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

func NewCharacter(
	name, race, class string,
	level int,
	abilities AbilityScores,
	modifiers AbilityModifiers,
	skills string,
) Character {
	return Character{
		Name:         name,
		Race:         race,
		Class:        class,
		Level:        level,
		AbilityScore: abilities,
		AbilityModifiers: modifiers,
		Skills:       skills,
	}
}
