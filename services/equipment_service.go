package services

import (
	"fmt"
	"test/domain"
	"test/infrastructure/equipment"
	"strings"
)

func (s *CharacterService) EquipCharacter(
	name, weapon, slot, armor, shield string,
) (error) {
	char, err := s.repo.View(name)
	if err != nil {
		return err
	}

	if weapon != "" && slot != "" {
		switch slot {
		case "main hand":
			if char.Equipment.Mainhand != "" {
				fmt.Printf("%s already occupied", slot)
			}
			char.EquipWeapon(slot, weapon)
		case "off hand":
			if char.Equipment.Offhand != "" {
				fmt.Printf("%s already occupied", slot)
			}
			char.EquipWeapon(slot, weapon)
		default:
			fmt.Printf("invalid slot: %s", slot)
		}
	}

	if armor != "" {
		if char.Equipment.Armor != "" {
			fmt.Print("armor already equipped")
		}
		char.EquipArmor(armor)
	}

	if shield != "" {
		if char.Equipment.Shield != "" {
			fmt.Print("shield already equipped")
		}
		char.EquipShield(shield)
	}

	if err := s.repo.SaveEquipment(char); err != nil {
    	fmt.Print("failed to save equipment: %w", err)
	}

	return nil
}
type EquipmentService struct {
	Repo equipment.EquipmentRepository
}

func (s EquipmentService) GetArmorAC(armorName string, dexMod int) (int, error) {
	equipments, err := s.Repo.LoadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to load equipments: %w", err)
	}

	armor, err := domain.GetArmorByName(equipments, armorName)
	if err != nil {
		return 0, err
	}

	return armor.CalculateAC(dexMod), nil
}
func (s EquipmentService) LoadArmorByName(name string) (*domain.Equipments, error) {
    equipments, err := s.Repo.LoadAll()
    if err != nil {
        return nil, fmt.Errorf("failed to load equipment list: %w", err)
    }

    normalizedInput := strings.ToLower(strings.TrimSpace(name))
    normalizedInput = strings.ReplaceAll(normalizedInput, " armor", "") // remove trailing "armor" for matching

    for _, e := range equipments {
        normalizedEquip := strings.ToLower(strings.TrimSpace(e.Name))
        normalizedEquip = strings.ReplaceAll(normalizedEquip, " armor", "")

        if normalizedInput == normalizedEquip {
            return &e, nil
        }
    }

    return nil, fmt.Errorf("armor not found: %s", name)
}

