package handler

import "github.com/gin-gonic/gin"

type APIError struct {
	Message string
	Code    int
}

func HandleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, APIError{
		Message: err.Error(),
		Code:    statusCode,
	})
}
