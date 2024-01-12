package repositories

import "database/sql"

type PetRepository interface {
	Store()
}

type petRepository struct {
	db *sql.DB
}

func NewPetRepository(db *sql.DB) PetRepository {
	return &petRepository{db: db}
}

func (p *petRepository) Store() {}
