package handler

import (
	"net/http"
	"strconv"
	"ticket-io/internal/shared/response"
	"ticket-io/internal/user/handler/dto"
	"ticket-io/internal/user/handler/mapper"
	"ticket-io/internal/user/service"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, total, statusMap, err := h.userService.GetAllWithStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users_formatted := mapper.UsersToResponse(users, statusMap)

	response.OK(c, dto.GetAllResponse{
		Total: total,
		Items: users_formatted,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	birthdate, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birthdate"})
		return
	}

	user, err := h.userService.Create(
		c.Request.Context(),
		req.Email,
		req.Name,
		birthdate,
		req.StatusID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
