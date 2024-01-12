package entities

type RegisterRequest struct {
	ProfilePicture string `json:"profilePicture" validate:"required,http_url"`
	Name           string `json:"name" validate:"required,max=24"`
	Email          string `json:"email" validate:"email,required,max=24"`
	Password       string `json:"password" validate:"required,max=12"`
	PhoneNumber    string `json:"phoneNumber" validate:"e164,required,max=20"`
	Location       string `json:"location" validate:"required,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=24"`
	Password string `json:"password" validate:"required,max=12"`
}

type LoginResponse struct {
	Id    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
}

type User struct {
	ProfilePicture string
	Name           string
	Email          string
	Password       string
	PhoneNumber    string
	Location       string
	Id             int
}
