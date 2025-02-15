package services

import (
	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Register(username, password string) (string, error)
	Login(username, password string) (string, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func HashPassword(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func (a *authService) generateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authService) Register(username, password string) (string, error) {
	exists, err := a.userRepo.ExistsByUsername(username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("user '%s' already exists", username)
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: hashed}

	if err := a.userRepo.CreateUser(user); err != nil {
		return "", err
	}

	token, err := a.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) Login(username, password string) (string, error) {
	user, err := a.userRepo.GetUserByName(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user '%s' not found", username)
	}

	if err := CheckPassword(user.PasswordHash, password); err != nil {
		return "", err
	}

	token, err := a.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
