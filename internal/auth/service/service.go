package service

import (
	"context"
	"errors"
	"go-gin-ticketing-backend/internal/auth/domain"
	authrepository "go-gin-ticketing-backend/internal/auth/repository/auth"
	permissionrepository "go-gin-ticketing-backend/internal/auth/repository/permission"
	"go-gin-ticketing-backend/internal/auth/schemas"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthService struct {
	userAuthRepository   authrepository.UserAuthRepository
	permissionRepository permissionrepository.PermissionRepository
	jwtSecret            []byte
	jwtTTL               time.Duration
}

func New(
	userAuthRepository authrepository.UserAuthRepository,
	permissionRepository permissionrepository.PermissionRepository,
	jwtSecret string,
	jwtTTL int64,
) *UserAuthService {

	return &UserAuthService{
		userAuthRepository:   userAuthRepository,
		permissionRepository: permissionRepository,
		jwtSecret:            []byte(jwtSecret),
		jwtTTL:               time.Duration(jwtTTL) * time.Second,
	}
}

func (s *UserAuthService) RegisterUser(ctx context.Context, body schemas.UserRegisterBody) (*domain.UserAuth, error) {

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

	return s.userAuthRepository.RegisterUser(ctx, user)
}

func (s *UserAuthService) LoginUser(ctx context.Context, body schemas.UserLoginBody) (string, error) {

	user, err := s.userAuthRepository.GetUserByEmail(ctx, body.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)) != nil {
		return "", errors.New("invalid credentials")
	}

	claims := schemas.CustomClaims{
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

func (s *UserAuthService) HasThisPermission(ctx context.Context, userID int64, userPermission string) (bool, error) {

	// Step 1. Get all user's userPermissions using its ID
	userPermissions, err := s.permissionRepository.GetPermissionsByUserID(ctx, userID)
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
