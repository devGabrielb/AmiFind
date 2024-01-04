package dtos

type RegisterRequest struct {
	Profile_picture_url string `json:"profilePictureUrl" validate:"required,max=255"`
	Name                string `json:"name" validate:"required,max=24"`
	Email               string `json:"email" validate:"required,max=24"`
	Password            string `json:"password" validate:"required,max=12"`
	PhoneNumber         string `json:"phoneNumber" validate:"required,max=20"`
	Location            string `json:"location" validate:"required,max=255"`
}

func (u *RegisterRequest) Validate() error {
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=24"`
	Password string `json:"password" validate:"required,max=12"`
}

func (u *LoginRequest) Validate() error {
	return nil
}

type LoginResponse struct {
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
	Id    int    `json:"id,omitempty"`
}
