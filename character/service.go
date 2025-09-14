package character

import "fmt"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateNewCharacter (name string, class string, race string, level int, str int, dex int, con int, intl int, wis int, cha int) (*Character, error) {
	char := &Character{
		Name: name,
		Class: class,
		Race: race,
		Level: level,
		Strength: str,
		Dexterity: dex,
		Constitution: con,
		Intelligence: intl,
		Wisdom: wis,
		Charisma: cha,

	}
	if err := s.repo.Create(char); err != nil {
		return nil, err
	}
	fmt.Println("Character created!")
	return char, nil
}
func (s *Service) ViewCharacterByName(name string) (string, error) {
	c, err := s.repo.View(name)
	if err != nil {
		return "", err
	}

	// Return a formatted string
	return fmt.Sprintf(
		"Name: %s, Race: %s, Class: %s, Level: %d, STR: %d, DEX: %d, CON: %d, INT: %d, WIS: %d, CHA: %d",
		c.Name, c.Race, c.Class, c.Level,
		c.Strength, c.Dexterity, c.Constitution,
		c.Intelligence, c.Wisdom, c.Charisma,
	), nil
}