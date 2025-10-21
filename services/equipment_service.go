package services

import "fmt"

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
