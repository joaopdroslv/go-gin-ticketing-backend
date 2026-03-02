package user

import (
	"errors"
	"go-gin-ticketing-backend/internal/domain"
	sharedschemas "go-gin-ticketing-backend/internal/shared/schemas"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {

	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {

	var query sharedschemas.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		sharedschemas.Failed(c, http.StatusBadRequest, "invalid query params")
		return
	}
	query.NormalizePagination()

	response, err := h.userService.GetAllUsers(c.Request.Context(), query)
	if err != nil {
		sharedschemas.Failed(c, http.StatusInternalServerError, "sorry, something went wrong")
		return
	}

	sharedschemas.OK(c, &response)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		sharedschemas.Failed(c, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			sharedschemas.OK(c, err.Error())
			return
		}

		sharedschemas.Failed(c, http.StatusInternalServerError, "internal server error")
		return
	}

	sharedschemas.OK(c, &user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var body CreateUserBody

	if err := c.ShouldBindJSON(&body); err != nil {
		sharedschemas.Failed(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), body)
	if err != nil {
		sharedschemas.Failed(c, http.StatusInternalServerError, "sorry, something went wrong")
		return
	}

	sharedschemas.OK(c, &user)
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		sharedschemas.Failed(c, http.StatusBadRequest, "invalid id")
		return
	}

	var body UpdateUserBody

	if err := c.ShouldBindJSON(&body); err != nil {
		sharedschemas.Failed(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.UpdateUserByID(c.Request.Context(), id, body)
	if err != nil {
		if errors.Is(err, domain.ErrNothingToUpdate) {
			sharedschemas.OK(c, err.Error())
			return
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			sharedschemas.Failed(c, http.StatusBadRequest, err.Error())
			return
		}

		sharedschemas.Failed(c, http.StatusInternalServerError, "sorry, something went wrong")
		return
	}

	sharedschemas.OK(c, &user)
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		sharedschemas.Failed(c, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.userService.DeleteUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			sharedschemas.Failed(c, http.StatusBadRequest, "user not found")
			return
		}

		sharedschemas.Failed(c, http.StatusNotFound, "sorry, something went wrong")
		return
	}

	sharedschemas.OK(c, &resp)
}
