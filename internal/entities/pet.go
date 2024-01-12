package entities

type PetRequest struct {
	ProfilePicture string `json:"profilePicture" validate:"required,http_url"`
	Name           string `json:"name" validate:"required,max=24"`
	Location       string `json:"location" validate:"required,max=24"`
	Sex            string `josn:"sex" validate:"required,max=3"`
	BreedId        int    `json:"breedId" validate:"required"`
}
type Pet struct {
	Id             int
	Name           string
	Location       string
	Sex            string
	ProfilePicture string
	BreedId        int
}
