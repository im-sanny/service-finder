package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(name, email, password string) error
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(name, email, password string) error {
	// 1. Basic Validation
	if name == "" || email == "" || password == "" {
		return fmt.Errorf("name, email, and password are required: %w", ErrInvalidInput)
	}

	// 2. Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Create the user object
	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	// 4. Save to DB via Repository
	return s.userRepo.Create(user)
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
