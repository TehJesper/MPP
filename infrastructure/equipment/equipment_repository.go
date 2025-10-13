package equipment
import (
	"database/sql"
	"test/domain"
)
type SQLRepository struct {
	db *sql.DB
}
func (r *SQLRepository) SaveEquipment(c domain.Character) error {
	_, err := r.db.Exec(`
		UPDATE characters
		SET mainhand = ?, offhand = ?, shield = ?, armor = ?
		WHERE name = ?
	`, c.Equipment.Mainhand, c.Equipment.Offhand, c.Equipment.Shield, c.Equipment.Armor, c.Name)

	return err
}
