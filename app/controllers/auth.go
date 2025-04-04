package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
)

func Register(c echo.Context) error {
	var payload req.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	return c.JSON(http.StatusOK, "OK")
}
