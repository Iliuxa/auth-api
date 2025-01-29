package usecase

import (
	"auth-api/internal/domain"
	"auth-api/internal/lib/jwt"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AuthUsecase interface {
	Login(ctx context.Context, email string, password string) (jwt string, err error)
	Register(ctx context.Context, email string, password string, name string) (jwt string, err error)
}

type authUsecase struct {
	userRepository domain.UserRepository
	log            *slog.Logger
}

func NewAuthUsecase(userRepository domain.UserRepository, log *slog.Logger) AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
		log:            log,
	}
}

func (a *authUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	const operation = "AuthUsecase.Login"

	log := a.log.With(
		slog.String("operation", operation),
		slog.String("email", email),
	)
	log.Info("Attempting to login user")

	user, err := a.userRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			a.log.Warn("User not found", slog.StringValue(err.Error()))
			return "", fmt.Errorf("%s: %w", operation, domain.ErrUserNotFound)
		}

		a.log.Error("Failed to get user", slog.StringValue(err.Error()))
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("Invalid credentials", slog.StringValue(err.Error()))
		return "", fmt.Errorf("%s: %w", operation, domain.ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, time.Hour)
	if err != nil {
		a.log.Error("Failed to create token", slog.StringValue(err.Error()))
		return "", fmt.Errorf("%s: %w", operation, err)
	}
	return token, nil
}

func (a *authUsecase) Register(ctx context.Context, email string, password string, name string) (string, error) {
	const operation = "AuthUsecase.Register"

	log := a.log.With(
		slog.String("operation", operation),
		slog.String("email", email),
	)
	log.Info("Register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate password", slog.StringValue(err.Error()))
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	err = a.userRepository.CreateUser(&domain.User{
		FullName: name,
		Email:    email,
		PassHash: passHash,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			a.log.Warn("User already exists", slog.StringValue(err.Error()))
			return "", fmt.Errorf("%s: %w", operation, domain.ErrUserAlreadyExists)
		}

		log.Error("Failed to create user", slog.StringValue(err.Error()))
		return "", fmt.Errorf("%s: %w", operation, err)
	}
	return "", nil
}
