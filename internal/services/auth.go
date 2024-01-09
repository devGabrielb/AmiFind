package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/devGabrielb/AmiFind/internal/dtos"
	"github.com/devGabrielb/AmiFind/internal/entities"
	"github.com/devGabrielb/AmiFind/internal/repositories"
	"github.com/devGabrielb/AmiFind/internal/utils"
	"github.com/devGabrielb/AmiFind/pkg/env"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidDto         = errors.New("invalid request")
	ErrNotFound           = errors.New("user not found")
	ErrGetSecretKey       = errors.New("error getting secret key")
	ErrCannotCreateUser   = errors.New("cannot create user")
	ErrEncryptPassword    = errors.New("error while encrypting password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrGenerateToken      = errors.New("error generating token")
	ErrInvalidParams      = errors.New("invalid parameters")
)

type AuthService interface {
	Register(ctx context.Context, registerRequest dtos.RegisterRequest) (int64, error)
	Login(ctx context.Context, loginRequest dtos.LoginRequest) (dtos.LoginResponse, error)
}

type authService struct {
	repository repositories.UserRepository
}

func NewService(repo repositories.UserRepository) AuthService {
	return &authService{
		repository: repo,
	}
}

func (s *authService) Login(ctx context.Context, loginRequest dtos.LoginRequest) (dtos.LoginResponse, error) {

	user, err := s.repository.FindByEmail(ctx, loginRequest.Email)

	if err != nil {

		return dtos.LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {

		return dtos.LoginResponse{}, ErrInvalidCredentials
	}

	env, err := env.TryGetEnv("SECRET_KEY")

	if err != nil {

		return dtos.LoginResponse{}, errors.Join(ErrGetSecretKey, err)
	}

	t := NewToken(env)
	token, err := t.GenerateToken(strconv.Itoa(user.Id))

	if err != nil {

		return dtos.LoginResponse{}, errors.Join(ErrGenerateToken, err)
	}
	return dtos.LoginResponse{Email: user.Email, Id: user.Id, Token: token}, nil
}

func (s *authService) Register(ctx context.Context, registerRequest dtos.RegisterRequest) (int64, error) {

	pass, err := utils.EncryptPassword(registerRequest.Password)

	if err != nil {

		return 0, err
	}

	user := entities.User{
		Profile_picture: registerRequest.Profile_picture,
		Name:            registerRequest.Name,
		Email:           registerRequest.Email,
		Password:        string(pass),
		PhoneNumber:     registerRequest.PhoneNumber,
		Location:        registerRequest.Location,
	}
	id, err := s.repository.Create(ctx, user)

	if err != nil {
		return 0, err
	}
	return id, nil
}
