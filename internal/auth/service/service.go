package service

import (
	"context"
	"errors"
	"strconv"
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
	jwtTTL     time.Duration
}

func New(
	repository repository.UserAuthRepository,
	jwtSecret string,
	jwtTTL int64,
) *UserAuthService {

	return &UserAuthService{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
		jwtTTL:     time.Duration(jwtTTL) * time.Second,
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

	claims := domain.CustomClaims{
		Role: "system", // Change this later, setting up all users as role=system
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}

func (s *UserAuthService) ValidateUserPermission(ctx context.Context, userID int64, userPermission string) (bool, error) {

	// Step 1. Get all user's userPermissions using its ID
	userPermissions, err := s.repository.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// Step 2. Creating a permissions map (with empty structs) for each user permissions
	permissionsMap := make(map[string]struct{})

	for _, permission := range userPermissions {
		permissionsMap[permission.Name] = struct{}{}
	}

	// Step 3. Validating if the user has the required permission
	_, ok := permissionsMap[userPermission]

	return ok, nil
}
