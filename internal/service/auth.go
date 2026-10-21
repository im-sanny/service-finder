package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/internal/model"
	"github.com/im-sanny/service-finder/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(u *model.User) error
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(u *model.User) error {
	// Hash password
	if u.PasswordHash == "" {
		return fmt.Errorf("password required: %w", ErrInvalidInput)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = string(hash)

	return s.userRepo.Create(u)
}

func (s *authService) Login(email, password string) (string, error) {
	// 1. Get the stored hash for this email
	hash, err := s.userRepo.GetPasswordHashByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials: %w", ErrInvalidInput)
	}

	// 2. Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", ErrInvalidInput)
	}

	// 3. Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
