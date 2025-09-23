package services

import "test/domain"

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
	char := domain.NewCharacter(name, race, class, level, str, dex, con, intel, wis, cha)
	return s.repo.Save(char)
}

func (s *CharacterService) ViewCharacter(name string) (domain.Character, error) {
	return s.repo.View(name)
}

func (s *CharacterService) DeleteCharacter(name string) {
	s.repo.Delete(name)
}
