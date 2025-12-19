package service

import (
	"context"
	"errors"
	"ticket-io/internal/auth/domain"
	"ticket-io/internal/auth/dto"
	"ticket-io/internal/auth/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthService struct {
	repository repository.UserAuthRepository
	jwtSecret  []byte
}

func New(repository repository.UserAuthRepository, jwtSecret string) *UserAuthService {

	return &UserAuthService{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
	}
}

func (s *UserAuthService) RegisterUser(ctx context.Context, body dto.UserRegisterBody) (*domain.UserAuth, error) {

	birthdate, err := time.Parse("2006-01-02", body.Birthdate)
	if err != nil {
		return nil, err
	}

	defaultStatusID := 1 // Creating all users with "active" status by default
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	user, err := domain.NewUserAuth(body.Name, birthdate, int64(defaultStatusID), body.Email, string(hash))
	if err != nil {
		return nil, err
	}

	return s.repository.RegisterUser(ctx, user)
}

func (s *UserAuthService) LoginUser(ctx context.Context, body dto.UserLoginBody) (string, error) {

	user, err := s.repository.GetUserByEmail(ctx, body.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)) != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}
