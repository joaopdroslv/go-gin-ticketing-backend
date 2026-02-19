package handler

import (
	"database/sql"
	"errors"
	"go-gin-ticketing-backend/internal/shared/responses"
	"go-gin-ticketing-backend/internal/user/schemas"
	userservice "go-gin-ticketing-backend/internal/user/service/user"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *userservice.UserService
	logger      *slog.Logger
}

func New(logger *slog.Logger, userService *userservice.UserService) *UserHandler {

	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
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
		// Do not treat sql in the handler
		if errors.Is(err, sql.ErrNoRows) {
			responses.OK(c, "user not found")
			return
		}

		// TODO: Currently losing the stacktrace,
		// the error is being considered thrown from the handler
		// instead of one of the previous layers.
		h.logger.Error(
			"get user by id",
			"error", err,
			"user_id", id,
		)
		responses.Failed(c, http.StatusInternalServerError, "internal server error")
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
