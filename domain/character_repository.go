package domain

type CharacterRepository interface {
	Save(c Character) (Character, error)
	View(name string) (Character, error)
	Delete(name string)
	SaveEquipment(c Character) error
	UpdateDerivedStats(c Character, ac int, initiative int, passive int) error
}
