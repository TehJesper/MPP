package infrastructure

import "test/domain"

type CharacterRepository struct {
	storage []domain.Character
}

func NewCharacterRepository() *CharacterRepository {
	return &CharacterRepository{
		storage: []domain.Character{},
	}
}

func (r *CharacterRepository) Save(c domain.Character) {
	r.storage = append(r.storage, c)
}
