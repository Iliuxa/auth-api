package repository

import (
	"auth-api/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (err error) {
	const operation = "repository.CreateUser"

	stmt, err := r.db.Prepare(`insert into users (fullName, email, password) values ($1, $2, $3) RETURNING id`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.FullName, user.Email, user.PassHash)
	err = row.Scan(&user.Id)
	if err != nil {
		err = fmt.Errorf("%s: %w", operation, err)
	}

	return
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const operation = "repository.FindByEmail"

	stmt, err := r.db.Prepare(`SELECT id, fullName, email, password as passHash FROM users WHERE email = $1`)
	if err != nil {
		return &domain.User{}, fmt.Errorf("%s: %w", operation, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	var user domain.User
	err = row.Scan(&user.Id, &user.FullName, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.User{}, fmt.Errorf("%s: %w", operation, domain.ErrUserNotFound)
		}

		return &domain.User{}, fmt.Errorf("%s: %w", operation, err)
	}

	return &user, nil
}
