package equipment

import (
	"encoding/json"
	"fmt"
	"os"

	"test/domain"
)

type EquipmentRepository struct {
	FilePath string
}

func (r EquipmentRepository) LoadAll() ([]domain.Equipments, error) {
	data, err := os.ReadFile(r.FilePath)

	
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var equipments []domain.Equipments
	if err := json.Unmarshal(data, &equipments); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return equipments, nil
}
