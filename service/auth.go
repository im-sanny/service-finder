package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/im-sanny/service-finder/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (string, error) // returns jwt token
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(email, password string) error {
	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// save user
	return &s.userRepo.CreateUser(email, string(hash))
}

func (s *authService) Login(email, password string) (string, error) {
	// get user by email
	hash, err := s.userRepo.GetPasswordHashByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials: %w", ErrNotFound)
	}

	// compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", ErrInvalidInput)
	}

	// generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
