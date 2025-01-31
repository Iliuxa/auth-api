package test

import (
	"auth-api/test/suite"
	auth "github.com/Iliuxa/protos/gen/proto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"testing"
	"time"
)

const passDefaultLen = 10

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := getFakePassword()

	loginInfo := &auth.LoginInfo{
		Email:    email,
		Password: pass,
	}

	resp, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Name:  "1",
		Login: loginInfo,
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

	loginResp, err := st.AuthClient.Login(ctx, loginInfo)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResp.GetJwt())

	tokenParsed2, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)

	claims2, ok := tokenParsed2.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, claims["email"].(string), claims2["email"].(string))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims2["exp"].(float64), deltaSec)
}

func TestRegisterLogin_Register(t *testing.T) {
	ctx, st := suite.New(t)

	type testType struct {
		testName    string
		expectedErr string
		data        auth.RegisterRequest
	}

	tests := []testType{
		{
			"empty name",
			"Validation Error",
			auth.RegisterRequest{
				Name:  "",
				Login: &auth.LoginInfo{Email: gofakeit.Email(), Password: getFakePassword()},
			},
		},
		{
			"empty email",
			"Validation Error",
			auth.RegisterRequest{
				Name:  "Name",
				Login: &auth.LoginInfo{Email: "", Password: getFakePassword()},
			},
		},
		{
			"invalid email",
			"Validation Error",
			auth.RegisterRequest{
				Name:  "Name",
				Login: &auth.LoginInfo{Email: "1234", Password: getFakePassword()},
			},
		},
		{
			"empty password",
			"Validation Error",
			auth.RegisterRequest{
				Name:  "Name",
				Login: &auth.LoginInfo{Email: gofakeit.Email(), Password: ""},
			},
		},
	}

	wg := &sync.WaitGroup{}
	toProc := make(chan testType, 5)

	for i := 0; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if test, ok := <-toProc; ok {
					_, err := st.AuthClient.Register(ctx, &test.data)
					require.NoError(t, err)
					require.Contains(t, err.Error(), test.expectedErr)
				} else {
					return
				}
			}
		}()
	}

	go func() {
		for _, test := range tests {
			toProc <- test
		}
		close(toProc)
	}()

	wg.Wait()
}

func getFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}
