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
	res, err := r.db.Exec(`
		INSERT INTO characters
		(name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma, skills)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Name, c.Race, c.Class, c.Level,
		c.Strength, c.Dexterity, c.Constitution,
		c.Intelligence, c.Wisdom, c.Charisma,
		c.Skills,
	)
	if err != nil {
		return c, err
	}

	id, _ := res.LastInsertId()
	c.ID = int(id)
	return c, nil
}

func (r *SQLRepository) View(name string) (domain.Character, error) {
	var c domain.Character
	row := r.db.QueryRow(`
		SELECT id, name, race, class, level,
		       strength, dexterity, constitution, intelligence, wisdom, charisma, skills
		FROM characters
		WHERE name = ?`, name)

	err := row.Scan(
		&c.ID, &c.Name, &c.Race, &c.Class, &c.Level,
		&c.Strength, &c.Dexterity, &c.Constitution,
		&c.Intelligence, &c.Wisdom, &c.Charisma,
		&c.Skills,
	)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (r *SQLRepository) Delete(name string) {
    r.db.Exec("DELETE FROM characters WHERE name = ?", name)
}
