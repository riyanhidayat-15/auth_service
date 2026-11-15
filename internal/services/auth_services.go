package services

import (
	"auth-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type authService struct {
	userRepo domain.UserRepository
}

func NewAuthService (userRepo domain.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register (ctx context.Context, username, password string) (*domain.User, error) {
	existing, _ := s.userRepo.FindByUsername(ctx, username)
	if existing != nil {
		return nil, errors.New("Username sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Gagal hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, username, string(hashedPassword))
	if err != nil {
		return  nil, err
	}

	return user, nil
}

var ErrInvalidCredentials = errors.New("username atau password salah")

func (s *authService) Login (ctx context.Context, username, password string) (string, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.RegisteredClaims{
		Subject:strconv.FormatInt(user.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}