package handler

import (
	"net/http"
	"strconv"
	"ticket-io/internal/shared/enums"
	"ticket-io/internal/shared/errors"
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
		response.Fail(c, http.StatusInternalServerError, string(enums.ErrInternal))
		return
	}

	formatted_users := mapper.UsersToResponse(users, statusMap)

	response.OK(c,
		dto.GetAllResponse{
			Total: total,
			Items: formatted_users,
		},
	)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, string(enums.ErrInvalidID))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		response.Fail(c, http.StatusNotFound, string(enums.ErrNotFound))
		return
	}

	response.OK(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var body dto.UserCreateBody

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, string(enums.ErrBadRequest))
		return
	}

	birthdate, err := time.Parse("2006-01-02", body.Birthdate)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "The provided birthdate is invalid.")
		return
	}

	user, err := h.userService.Create(
		c.Request.Context(),
		body.Email,
		body.Name,
		birthdate,
		body.StatusID,
	)
	if err != nil {
		response.Fail(c, 500, string(enums.ErrInternal))
		return
	}

	// TODO: return formatted user

	response.OK(c, user)
}

func (h *UserHandler) UpdateByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, string(enums.ErrInvalidID))
		return
	}

	var body dto.UserUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateByID(c, id, body)
	if err != nil {
		switch err {
		case errors.ErrNothingToUpdate:
			response.Fail(c, http.StatusBadRequest, string(enums.ErrBadRequest))
		case errors.ErrZeroRowsAffected:
			response.Fail(c, http.StatusBadRequest, string(enums.ErrZeroRowsAffected))
		default:
			response.Fail(c, http.StatusInternalServerError, string(enums.ErrInternal))
		}
		return
	}

	// TODO: return formatted user

	response.OK(c, user)
}

func (h *UserHandler) DeleteByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, string(enums.ErrInvalidID))
		return
	}

	success, err := h.userService.DeleteByID(c.Request.Context(), int64(id))
	if err != nil {
		response.Fail(c, http.StatusNotFound, string(enums.ErrNotFound))
		return
	}

	response.OK(c, dto.UserDeleteResponse{
		ID:      id,
		Deleted: success,
	})
}
