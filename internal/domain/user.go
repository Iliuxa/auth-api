package domain

import (
	"context"
	"errors"
)

type User struct {
	Id       int
	FullName string
	Email    string
	PassHash []byte
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	//FindByID(id int) (*User, error)
	//Create(user *User) error
	//Update(user *User) error
	//Delete(id int) error
}
