package repositories

import (
	"context"
	"database/sql"

	"github.com/devGabrielb/AmiFind/internal/entities"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (entities.User, error)
	Create(ctx context.Context, user entities.User) (int64, error)
}

type userRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	query := "SELECT id,profile_picture,name,email,password,phoneNumber,location FROM users WHERE email=?;"
	row := r.db.QueryRowContext(ctx, query, email)
	user := entities.User{}

	err := row.Scan(&user.Id, &user.Profile_picture, &user.Name, &user.Email, &user.Password, &user.PhoneNumber, &user.Location)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user entities.User) (int64, error) {
	u, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users(profile_picture,name, email, password, phoneNumber, location) VALUES (?,?,?,?,?,?);",
		user.Profile_picture, user.Name, user.Email, user.Password, user.PhoneNumber, user.Location)
	if err != nil {
		return 0, err
	}
	id, err := u.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
