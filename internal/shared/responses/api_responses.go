package responses

import "github.com/gin-gonic/gin"

type APIResponse[T any] struct {
	StatusCode int    `json:"code"`
	Success    bool   `json:"success"`
	Resource   T      `json:"resource,omitempty"`
	Error      string `json:"error,omitempty"`
}

func OK[T any](c *gin.Context, data T) {

	c.JSON(200, APIResponse[T]{
		StatusCode: 200,
		Success:    true,
		Resource:   data,
	})
}

func Created[T any](c *gin.Context, data T) {

	c.JSON(201, APIResponse[T]{
		StatusCode: 201,
		Success:    true,
		Resource:   data,
	})
}

func Failed(c *gin.Context, status int, err string) {

	c.JSON(status, APIResponse[any]{
		StatusCode: status,
		Success:    false,
		Error:      err,
	})
}
