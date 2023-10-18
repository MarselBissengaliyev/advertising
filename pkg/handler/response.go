package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c echo.Context, statusCode int, message string) {
	logrus.Error(message)
	err := c.JSON(statusCode, errorResponse{message})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}
