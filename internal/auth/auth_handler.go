package auth

import (
	"errors"
	"go-gin-ticketing-backend/internal/domain"
	"go-gin-ticketing-backend/internal/shared/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {

	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {

	var body RegisterBody

	if err := c.ShouldBindJSON(&body); err != nil {
		schemas.Failed(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.authService.RegisterUser(c, body)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			schemas.Failed(c, http.StatusConflict, "this email address is already in use")
			return
		}

		schemas.Failed(c, http.StatusInternalServerError, "sorry, something went wrong")
		return
	}

	schemas.OK(c, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {

	var body LoginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		schemas.Failed(c, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.LoginUser(c, body)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			schemas.Failed(c, http.StatusUnauthorized, err.Error())
			return
		}

		if domain.IsUserStatusRelated(err) {
			schemas.Failed(c, http.StatusForbidden, err.Error())
			return
		}

		schemas.Failed(c, http.StatusInternalServerError, "sorry, something went wrong")
		return
	}

	schemas.OK(c, gin.H{"message": "logged in successfully", "token": token})
}
