package services

import (
	"context"
	"errors"
	"strconv"

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
	Register(ctx context.Context, user entities.User) (int64, error)
	Login(ctx context.Context, email string, pass string) (entities.User, error)
	GenerateToken(user entities.User) (string, error)
}

type authService struct {
	repository repositories.UserRepository
}

func NewService(repo repositories.UserRepository) AuthService {
	return &authService{
		repository: repo,
	}
}

func (s *authService) Login(ctx context.Context, email string, pass string) (entities.User, error) {

	user, err := s.repository.FindByEmail(ctx, email)

	if err != nil {

		return entities.User{}, ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {

		return entities.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *authService) Register(ctx context.Context, user entities.User) (int64, error) {

	pass, err := utils.EncryptPassword(user.Password)

	if err != nil {

		return 0, err
	}

	user.Password = string(pass)
	id, err := s.repository.Create(ctx, user)

	if err != nil {
		return 0, ErrNotFound
	}
	return id, nil
}

func (s *authService) GenerateToken(user entities.User) (string, error) {
	env, err := env.TryGetEnv("SECRET_KEY")

	if err != nil {

		return "", errors.Join(ErrGetSecretKey, err)
	}

	t := NewToken(env)
	token, err := t.GenerateToken(strconv.Itoa(user.Id))

	if err != nil {

		return "", errors.Join(ErrGenerateToken, err)
	}

	return token, nil
}
