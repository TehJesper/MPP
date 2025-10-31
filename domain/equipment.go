package domain

import "fmt"

type Equipments struct {
	Index             string            `json:"index"`
	Name              string            `json:"name"`
	EquipmentCategory EquipmentCategory `json:"equipment_category"`
	Properties        []Property        `json:"properties,omitempty"`
	ArmorClass        *ArmorClass       `json:"armor_class,omitempty"`
	ArmorCategory     string            `json:"armor_category,omitempty"`
}

type EquipmentCategory struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type Property struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type ArmorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}

func (e Equipments) CalculateAC(dexMod int) int {
	baseAC := 10
	if e.ArmorClass != nil {
		baseAC = e.ArmorClass.Base
	}

	switch e.ArmorCategory {
	case "Light":
		if e.ArmorClass != nil && e.ArmorClass.DexBonus {
			return baseAC + dexMod
		}
		return baseAC
	case "Medium":
		mod := dexMod
		if mod > 2 {
			mod = 2
		}
		if e.ArmorClass != nil && e.ArmorClass.DexBonus {
			return baseAC + mod
		}
		return baseAC
	case "Heavy":
		return baseAC
	default:
		return 10 + dexMod
	}
}

func GetArmorByName(equipments []Equipments, armorName string) (*Equipments, error) {
	for _, eq := range equipments {
		if eq.EquipmentCategory.Name == "Armor" && eq.Name == armorName {
			return &eq, nil
		}
	}
	return nil, fmt.Errorf("armor %s not found", armorName)
}
