package character

import (
	"database/sql"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(c *Character) error {
	result, err := r.db.Exec(
		`INSERT INTO characters (name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Name, c.Race, c.Class, c.Level, c.Strength, c.Dexterity, c.Constitution, c.Intelligence, c.Wisdom, c.Charisma,
	)
	if  err != nil {
		return err
	}
	shit := result
	if shit != nil{
		return err
	}
	return err
}