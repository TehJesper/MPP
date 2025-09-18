package character

import (
	"database/sql"
	"strings"
	"fmt"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(c *Character) error {
	// skillsJSON, _ := json.Marshal(c.Skills)
	result, err := r.db.Exec(
		`INSERT INTO characters (name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma, skills) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Name, c.Race, c.Class, c.Level, c.Strength, c.Dexterity, c.Constitution, c.Intelligence, c.Wisdom, c.Charisma, c.Skills,
	)
	if  err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return fmt.Errorf("character with name %s already exists", c.Name)
		}
		return err	}
	shit := result
	if shit != nil{
		return err
	}
	return err
}
func (r *SQLRepository) View(name string) (*Character, error) {
	row := r.db.QueryRow(`
		SELECT name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma, skills
		FROM characters
		WHERE name = ?
	`, name)

	var c Character
	err := row.Scan(
		&c.Name, &c.Race, &c.Class, &c.Level,
		&c.Strength, &c.Dexterity, &c.Constitution,
		&c.Intelligence, &c.Wisdom, &c.Charisma, &c.Skills,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(`character "%s" not found`, name)
		}
		return nil, err
	}
	return &c, nil
}
func (r *SQLRepository) Delete(name string) {
    r.db.Exec("DELETE FROM characters WHERE name = ?", name)
}
