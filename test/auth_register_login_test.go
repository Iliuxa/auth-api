package test

import (
	"auth-api/test/suite"
	auth "github.com/Iliuxa/protos/gen/proto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	type testType struct {
		testName    string
		expectedErr string
		name        string
		email       string
		password    string
	}

	tests := []testType{
		{
			testName:    "empty name",
			expectedErr: "Validation Error",
			name:        "",
			email:       gofakeit.Email(),
			password:    getFakePassword(),
		},
		{
			testName:    "empty email",
			expectedErr: "Validation Error",
			name:        "Name",
			email:       "",
			password:    getFakePassword(),
		},
		{
			testName:    "invalid email",
			expectedErr: "Validation Error",
			name:        "Name",
			email:       "1234",
			password:    getFakePassword(),
		},
		{
			testName:    "empty password",
			expectedErr: "Validation Error",
			name:        "Name",
			email:       gofakeit.Email(),
			password:    "",
		},
	}

	wg := sync.WaitGroup{}

	for _, test := range tests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Run(test.testName, func(t *testing.T) {
				_, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
					Name:  test.name,
					Login: &auth.LoginInfo{Email: test.email, Password: test.password},
				})
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedErr)
			})
		}()
	}

	wg.Wait()
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	type testType struct {
		testName    string
		expectedErr string
		email       string
		password    string
	}

	tests := []testType{
		{
			testName:    "empty email",
			expectedErr: "Validation Error",
			email:       "",
			password:    getFakePassword(),
		},
		{
			testName:    "invalid email",
			expectedErr: "Validation Error",
			email:       "1234",
			password:    getFakePassword(),
		},
		{
			testName:    "empty password",
			expectedErr: "Validation Error",
			email:       gofakeit.Email(),
			password:    "",
		},
	}

	wg := sync.WaitGroup{}

	for _, test := range tests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Run(test.testName, func(t *testing.T) {
				_, err := st.AuthClient.Login(ctx, &auth.LoginInfo{
					Email:    test.email,
					Password: test.password,
				})
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedErr)
			})
		}()
	}

	email := gofakeit.Email()
	password := getFakePassword()
	_, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Name:  "test",
		Login: &auth.LoginInfo{Email: email, Password: password},
	})
	require.NoError(t, err)

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Run("Invalid password", func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &auth.LoginInfo{
				Email:    email,
				Password: "12345",
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), "Invalid email or password")
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Run("User not found", func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &auth.LoginInfo{
				Email:    gofakeit.Email(),
				Password: password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), "Invalid email or password")
		})
	}()

	wg.Wait()
}

func getFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}
