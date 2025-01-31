package test

import (
	"auth-api/test/suite"
	auth "github.com/Iliuxa/protos/gen/proto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const passDefaultLen = 10

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := getFakePassword()

	resp, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Name: "1",
		Login: &auth.LoginInfo{
			Email:    email,
			Password: pass,
		},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetJwt())

	token := resp.GetJwt()
	require.NotEmpty(t, token)
	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	const deltaSec = 1
	assert.Equal(t, email, claims["email"].(string))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSec)

}

func getFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}
