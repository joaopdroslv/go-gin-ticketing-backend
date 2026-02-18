package handler

import (
	"go-gin-ticketing-backend/internal/auth/schemas"
	"go-gin-ticketing-backend/internal/auth/service"
	"go-gin-ticketing-backend/internal/shared/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAuthHandler struct {
	service *service.UserAuthService
}

func New(service *service.UserAuthService) *UserAuthHandler {

	return &UserAuthHandler{service: service}
}

func (h *UserAuthHandler) RegisterUser(c *gin.Context) {

	var body schemas.UserRegisterBody

	if err := c.ShouldBindJSON(&body); err != nil {
		responses.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.RegisterUser(c, body)
	if err != nil {
		responses.Failed(c, http.StatusInternalServerError, err.Error())
	}

	responses.OK(c, gin.H{"id": user.ID})
}

func (h *UserAuthHandler) LoginUser(c *gin.Context) {

	var body schemas.UserLoginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		responses.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.LoginUser(c, body)
	if err != nil {
		responses.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses.OK(c, gin.H{"token": token})
}
