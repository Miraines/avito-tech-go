package integration

import (
	"fmt"
	"testing"
	"time"

	"avito-tech-go/internal/repositories"
	"avito-tech-go/internal/services"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

const jwtSecret = "mysecret"

func TestIntegration_Auth_Register_Login(t *testing.T) {
	db := setupIntegrationDB(t)
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtSecret)

	t.Run("Successful registration", func(t *testing.T) {
		token, err := authService.Register("newuser", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		assert.NoError(t, err)
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			assert.Equal(t, "newuser", claims["username"])
			exp := int64(claims["exp"].(float64))
			assert.True(t, exp > time.Now().Add(70*time.Hour).Unix())
		} else {
			t.Error("failed to parse token claims")
		}
	})

	t.Run("Duplicate registration", func(t *testing.T) {
		_, err := authService.Register("duplicate", "password123")
		assert.NoError(t, err)
		_, err = authService.Register("duplicate", "password123")
		assert.Error(t, err)
	})

	t.Run("Successful login", func(t *testing.T) {
		_, err := authService.Register("loginuser", "securepwd")
		assert.NoError(t, err)
		token, err := authService.Login("loginuser", "securepwd")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		assert.NoError(t, err)
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			assert.Equal(t, "loginuser", claims["username"])
		} else {
			t.Error("failed to parse token claims")
		}
	})

	t.Run("Login for non-existing user", func(t *testing.T) {
		_, err := authService.Login("nonexistent", "anyPassword")
		assert.Error(t, err)
	})

	t.Run("Login with invalid password", func(t *testing.T) {
		_, err := authService.Register("wrongpwd", "correctpwd")
		assert.NoError(t, err)
		_, err = authService.Login("wrongpwd", "incorrectpwd")
		assert.Error(t, err)
	})
}
