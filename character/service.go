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