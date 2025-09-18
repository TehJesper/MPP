package character

type Character struct {
	Name  string
	Class string
	Race string
	Level int
	Strength int
	Dexterity int
	Constitution int
	Intelligence int
	Wisdom int 
	Charisma int
	Skills string
}
type Repository interface {
	Create(c *Character) error
	View(name string) (*Character, error)
	Delete(name string)
}