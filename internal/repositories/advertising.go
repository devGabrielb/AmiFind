package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/devGabrielb/AmiFind/internal/entities"
)

type AdvertisingRepository interface {
	Store(ctx context.Context, advadvertising entities.Advertising) (int, error)
	GetAdvertisingsByCriteria(ctx context.Context, queryParams map[string]string) ([]entities.Advertising, error)
}

type advertisingRepository struct {
	db *sql.DB
}

func NewAdvertisingRepository(db *sql.DB) AdvertisingRepository {
	return &advertisingRepository{db: db}
}

type AdvertisingQuery struct {
	Category   string `query:"category"`
	AnimalType string `query:"animalType"`
	Breed      string `query:"breed"`
	Sex        string `query:"sex"`
	Location   string `query:"location"`
	Term       string `query:"term"`
}

func (a *advertisingRepository) GetAdvertisingsByCriteria(ctx context.Context, queryParams map[string]string) ([]entities.Advertising, error) {
	where := "WHERE "
	Advertisings := make([]entities.Advertising, 0)
	category, ok := queryParams["category"]
	if ok {
		where += "ad.category= ?"
	}
	animalType, ok := queryParams["animalType"]
	if ok {
		where += "AND at.name= ? "
	}
	breed, ok := queryParams["breed"]
	if ok {
		where += "AND br.name = ? "
	}
	sex, ok := queryParams["sex"]
	if ok {
		where += "AND pe.sex = ? "
	}
	location, ok := queryParams["location"]
	if ok {
		where += "AND pe.location = ?;"
	}
	query, err := a.db.PrepareContext(ctx,
		`SELECT
			ad.id 
			ad.category,
    		po.title,po.date,
    		pe.profilePicture, pe.location,
		FROM 
    		advertisings ad
    		INNER JOIN posts po ON po.id = ad.postId
    		INNER JOIN pets pe ON po.petId = pe.id
    		INNER JOIN breeds br ON pe.breedId = br.id
    		INNER JOIN animalTypes at ON b.animalTypeId = at.id
`+where)
	if err != nil {
		return nil, err
	}
	rows, err := query.QueryContext(ctx, category, animalType, breed, sex, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		adv := entities.Advertising{}

		if err := rows.Scan(
			&adv.Id, &adv.Category,
			&adv.Post.Title, &adv.Post.Date,
			&adv.Post.Pet.ProfilePicture, &adv.Post.Pet.Location); err != nil {
			return nil, err
		}
		Advertisings = append(Advertisings, adv)
	}
	return Advertisings, nil
}
func (a *advertisingRepository) Store(ctx context.Context, advertising entities.Advertising) (int, error) {
	log.Println(advertising)
	//create transaction
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	pet := advertising.Post.Pet

	petId, err := storePet(tx, ctx, pet)
	if err != nil {
		return 0, err
	}

	post := advertising.Post
	post.PetId = int(petId)
	postId, err := storePost(tx, ctx, post)
	if err != nil {
		return 0, err
	}
	advertising.PostId = int(postId)
	query, err := tx.PrepareContext(ctx, "INSERT INTO advertisings(status,category,userId,postId) VALUES(?, ?,?,?)")
	if err != nil {
		return 0, err
	}
	adv, err := query.ExecContext(ctx, advertising.Status, advertising.Category, advertising.UserId, advertising.PostId)
	if err != nil {
		return 0, err
	}
	advId, err := adv.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(advId), nil
}

func storePost(tx *sql.Tx, ctx context.Context, post entities.Post) (int64, error) {
	query, err := tx.PrepareContext(ctx, "INSERT INTO posts(title,description, date, petId) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	p, err := query.ExecContext(ctx, post.Title, post.Description, post.Date, post.PetId)
	if err != nil {
		return 0, err
	}

	postId, err := p.LastInsertId()
	if err != nil {
		return 0, err
	}
	return postId, nil
}
func storePet(tx *sql.Tx, ctx context.Context, pet entities.Pet) (int64, error) {
	query, err := tx.PrepareContext(ctx, "INSERT INTO pets(profilePicture, name,location, sex, breedId) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	p, err := query.ExecContext(ctx, pet.ProfilePicture, pet.Name, pet.Location, pet.Sex, pet.BreedId)
	if err != nil {
		return 0, err
	}

	petId, err := p.LastInsertId()
	if err != nil {
		return 0, err
	}
	return petId, nil
}
