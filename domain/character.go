package domain

type Character struct {
	ID           int
	Name         string
	Race         string
	Class        string
	Level        int
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
	Skills       string
}

func NewCharacter(
	name, race, class string,
	level, str, dex, con, intel, wis, cha int,
	skills string,
) Character {
	return Character{
		Name:         name,
		Race:         race,
		Class:        class,
		Level:        level,
		Strength:     str,
		Dexterity:    dex,
		Constitution: con,
		Intelligence: intel,
		Wisdom:       wis,
		Charisma:     cha,
		Skills: 	  skills,
	}
}
