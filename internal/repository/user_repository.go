package repository

import (
	"auth-api/internal/domain"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *domain.User) (err error) {
	_, err = r.db.Exec(`insert into users (fullName, email, password) values ($1, $2, $3)`, user.FullName, user.Email, user.PassHash)
	return
}

func (r *userRepository) FindByEmail(email string) (user *domain.User, err error) {

	return nil, nil
}
