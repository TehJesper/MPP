package infrastructure

import (
	"database/sql"
	"test/domain"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Save(c domain.Character) (domain.Character, error) {
	abilities := []any{
		c.AbilityScore.Strength,
		c.AbilityScore.Dexterity,
		c.AbilityScore.Constitution,
		c.AbilityScore.Intelligence,
		c.AbilityScore.Wisdom,
		c.AbilityScore.Charisma,
	}

	modifiers := []any{
		c.AbilityModifiers.Strength,
		c.AbilityModifiers.Dexterity,
		c.AbilityModifiers.Constitution,
		c.AbilityModifiers.Intelligence,
		c.AbilityModifiers.Wisdom,
		c.AbilityModifiers.Charisma,
	}

	args := []any{
		c.Name,
		c.Race,
		c.Class,
		c.Level,
		c.ProficiencyBonus,
	}
	args = append(args, abilities...)
	args = append(args, modifiers...)
	args = append(args, c.Skills)

	res, err := r.db.Exec(`
		INSERT INTO characters
		(
			name, race, class, level, proficiency_bonus, 
			strength, dexterity, constitution, intelligence, wisdom, charisma,
			strength_mod, dexterity_mod, constitution_mod, intelligence_mod, wisdom_mod, charisma_mod,
			skills
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, args...)
	if err != nil {
		return c, err
	}

	id, _ := res.LastInsertId()
	c.ID = int(id)
	return c, nil
}


func (r *SQLRepository) View(name string) (domain.Character, error) {
	var c domain.Character
	var mainhand, offhand, shield, armor sql.NullString
	row := r.db.QueryRow(`
		SELECT id, name, race, class, level, proficiency_bonus,
		       strength, dexterity, constitution, intelligence, wisdom, charisma,
		       strength_mod, dexterity_mod, constitution_mod, intelligence_mod, wisdom_mod, charisma_mod,
		       skills,
			   mainhand, offhand, shield, armor
		FROM characters
		WHERE name = ?`, name)

	err := row.Scan(
		&c.ID, &c.Name, &c.Race, &c.Class, &c.Level, &c.ProficiencyBonus,
		&c.AbilityScore.Strength, &c.AbilityScore.Dexterity, &c.AbilityScore.Constitution,
		&c.AbilityScore.Intelligence, &c.AbilityScore.Wisdom, &c.AbilityScore.Charisma,
		&c.AbilityModifiers.Strength, &c.AbilityModifiers.Dexterity, &c.AbilityModifiers.Constitution,
		&c.AbilityModifiers.Intelligence, &c.AbilityModifiers.Wisdom, &c.AbilityModifiers.Charisma,
		&c.Skills,
		&mainhand, &offhand, &shield, &armor,
		// &c.Equipment.Mainhand, &c.Equipment.Offhand, &c.Equipment.Shield, &c.Equipment.Armor,
	)
	if err != nil {
		return c, err
	}
	if mainhand.Valid {
		c.Equipment.Mainhand = mainhand.String
	}
	if offhand.Valid {
		c.Equipment.Offhand = offhand.String
	}
	if shield.Valid {
		c.Equipment.Shield = shield.String
	}
	if armor.Valid {
		c.Equipment.Armor = armor.String
	}
	return c, nil
}

func (r *SQLRepository) SaveEquipment(c domain.Character) error {
	_, err := r.db.Exec(`
		UPDATE characters
		SET mainhand = ?, offhand = ?, shield = ?, armor = ?
		WHERE name = ?
	`, c.Equipment.Mainhand, c.Equipment.Offhand, c.Equipment.Shield, c.Equipment.Armor, c.Name)

	return err
}

func (r *SQLRepository) Delete(name string) {
	r.db.Exec("DELETE FROM characters WHERE name = ?", name)
}

