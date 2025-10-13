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
		char.EquipWeapon(slot, weapon)
	}

	if armor != "" {
		char.EquipArmor(armor)
	}

	if shield != "" {
		char.EquipShield(shield)
	}

	if err := s.repo.SaveEquipment(char); err != nil {
    	return fmt.Errorf("failed to save equipment: %w", err)
	}

	return nil
}
