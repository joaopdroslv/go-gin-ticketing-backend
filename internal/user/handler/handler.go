package handler

import (
	"net/http"
	"strconv"
	"ticket-io/internal/shared/responses"
	"ticket-io/internal/user/schemas"
	userservice "ticket-io/internal/user/service/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *userservice.UserService
}

func New(s *userservice.UserService) *UserHandler {

	return &UserHandler{userService: s}
}

func (h *UserHandler) ListUsers(c *gin.Context) {

	resp, err := h.userService.ListUsers(c.Request.Context())
	if err != nil {
		responses.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses.OK(c, &resp)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		responses.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		responses.Failed(c, http.StatusNotFound, err.Error())
		return
	}

	responses.OK(c, &user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var body schemas.UserCreateBody

	if err := c.ShouldBindJSON(&body); err != nil {
		responses.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), body)
	if err != nil {
		responses.Failed(c, 500, err.Error())
		return
	}

	responses.OK(c, &user)
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		responses.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	var body schemas.UserUpdateBody

	if err := c.ShouldBindJSON(&body); err != nil {
		responses.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateUserByID(c.Request.Context(), id, body)
	if err != nil {
		responses.Failed(c, 500, err.Error())
		return
	}

	responses.OK(c, &user)
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		responses.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.userService.DeleteUserByID(c.Request.Context(), id)
	if err != nil {
		responses.Failed(c, http.StatusNotFound, err.Error())
		return
	}

	responses.OK(c, &resp)
}
