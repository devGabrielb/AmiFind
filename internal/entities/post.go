package entities

import "time"

type PostRequest struct {
	Title       string     `json:"title" validate:"required,max=50"`
	Description string     `json:"description" validate:"required,max=500"`
	Date        string     `json:"date" validate:"required,datetime"`
	Pet         PetRequest `json:"pet" validate:"required"`
}

type Post struct {
	Id          int
	Title       string
	Description string
	Date        time.Time
	Pet         Pet
	PetId       int
}

func (po *Post) NewPet(pet PetRequest) {
	po.Pet = Pet{
		ProfilePicture: pet.ProfilePicture,
		Name:           pet.Name,
		Location:       pet.Location,
		Sex:            pet.Sex,
		BreedId:        pet.BreedId,
	}
}
