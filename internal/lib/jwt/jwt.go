package jwt

import (
	"auth-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user *domain.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"email":  user.Email,
		"exp":    time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
