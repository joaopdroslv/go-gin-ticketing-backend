package auth

import (
	"context"
	"errors"
	"go-gin-ticketing-backend/internal/domain"
	"go-gin-ticketing-backend/internal/shared/enums"
	"go-gin-ticketing-backend/internal/shared/schemas"
	"go-gin-ticketing-backend/internal/shared/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository AuthRepository
	jwtSecret      []byte
	jwtTTL         time.Duration
}

func NewAuthService(
	authRepository AuthRepository,
	jwtSecret string,
	jwtTTL int64,
) *AuthService {

	return &AuthService{
		authRepository: authRepository,
		jwtSecret:      []byte(jwtSecret),
		jwtTTL:         time.Duration(jwtTTL) * time.Second,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, body RegisterBody) error {

	birthdate, err := time.Parse("2006-01-02", body.Birthdate)
	if err != nil {
		return err
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	data := &RegisterUserData{
		UserStatusID: int64(enums.Active),
		Name:         body.Name,
		Birthdate:    birthdate,
		Email:        body.Email,
		PasswordHash: string(passwordHash),
	}

	err = s.authRepository.RegisterUser(ctx, data)
	if err != nil {
		if utils.IsDuplicateKey(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (s *AuthService) LoginUser(ctx context.Context, body LoginBody) (string, error) {

	userCredential, err := s.authRepository.GetUserByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	err = s.validateUserStatus(userCredential)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userCredential.PasswordHash), []byte(body.Password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	claims := schemas.CustomClaims{
		Role: "system", // Change this later, setting up all users as role=system
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userCredential.UserInfo.ID, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) validateUserStatus(userCredential *UserCredential) error {

	switch userCredential.UserInfo.UserStatusID {
	case int64(enums.Inactive):
		return domain.ErrInactiveUser
	case int64(enums.EmailConfirmationPending):
		return domain.ErrUserEmailConfirmationPending
	case int64(enums.PasswordCreationPending):
		return domain.ErrUserPasswordCreationPending
	case int64(enums.Deleted):
		return domain.ErrDeletedUser
	}

	return nil
}
