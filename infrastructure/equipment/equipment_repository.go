package equipment

type EquipmentWrapper struct {
	Data EquipmentData `json:"data"`
}

type EquipmentData struct {
	Equipments []Equipment `json:"equipments,omitempty"`
}

type Equipment struct {
	Index             string            `json:"index,omitempty"`
	Name              string            `json:"name,omitempty"`
	EquipmentCategory *Category         `json:"equipment_category,omitempty"`
	Properties        []Property        `json:"properties,omitempty"`
	ArmorClass        *ArmorClass       `json:"armor_class,omitempty"`
	ArmorCategory     string            `json:"armor_category,omitempty"` // e.g., Light, Medium, Heavy
}

type Category struct {
	Name  string `json:"name"`
	Index string `json:"index"`
}

type Property struct {
	Name  string `json:"name"`
	Index string `json:"index"`
}

type ArmorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}
