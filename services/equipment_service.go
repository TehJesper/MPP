package services

import (
	"fmt"
	"strings"
	"test/domain"
	"test/infrastructure/equipment"
)

type EquipmentService struct{}

func (s *CharacterService) EquipCharacter(
	name, weapon, slot, armor, shield string,
) error {
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
	ac, err := s.CalculateArmor(&char)
	if err != nil {
		fmt.Println("Failed to calculate AC:", err)
	} else {
		char.ArmorClass = ac
		fmt.Println("New AC for", char.Name, "is", ac)
	}

	if err := s.repo.SaveEquipment(char); err != nil {
		fmt.Print("failed to save equipment: %w", err)
	}

	return nil
}

func (s *EquipmentService) LoadArmorByName(name string) (*domain.Equipments, error) {
	equipments, err := equipment.LoadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load equipment list: %w", err)
	}

	normalizedInput := strings.ToLower(strings.TrimSpace(name))
	normalizedInput = strings.ReplaceAll(normalizedInput, " armor", "")

	for _, e := range equipments {
		normalizedEquip := strings.ToLower(strings.TrimSpace(e.Name))
		normalizedEquip = strings.ReplaceAll(normalizedEquip, " armor", "")
		if normalizedInput == normalizedEquip {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("armor not found: %s", name)
}

func (s *EquipmentService) GetArmorAC(armorName string, dexMod int) (int, error) {
	armor, err := s.LoadArmorByName(armorName)
	if err != nil {
		return 0, err
	}

	return armor.CalculateAC(dexMod), nil
}

func (s *CharacterService) CalculateArmor(c *domain.Character) (int, error) {
	ac := 0
	dex := c.AbilityModifiers.Dexterity

	if c.Equipment.Armor != "" {
		armor, err := s.eq.LoadArmorByName(c.Equipment.Armor)
		if err != nil {
			fmt.Println("Armor not found, default AC")
			ac = 10 + dex
		} else {
			ac = armor.CalculateAC(dex)
		}
	} else {
		switch strings.ToLower(c.Class) {
		case "barbarian":
			ac = 10 + dex + c.AbilityModifiers.Constitution
		case "monk":
			if c.Equipment.Shield == "" {
				ac = 10 + dex + c.AbilityModifiers.Wisdom
			} else {
				ac = 10 + dex
			}
		default:
			ac = 10 + dex
		}
	}

	if c.Equipment.Shield != "" {
		shield, err := s.eq.LoadArmorByName(c.Equipment.Shield)
		if err == nil && shield.ArmorClass != nil {
			ac += shield.ArmorClass.Base
		}
	}

	return ac, nil
}
