package equipment

import (
	"encoding/json"
	"fmt"
	"os"
	"test/domain"
)

func LoadAll() ([]domain.Equipments, error) {
	data, err := os.ReadFile("5e-SRD-Equipment.json")

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var equipments []domain.Equipments
	if err := json.Unmarshal(data, &equipments); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return equipments, nil
}
