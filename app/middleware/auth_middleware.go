package middleware

import (
	"net/http"
	"strings"

	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "invalid authorization format",
			})
		}

		claims, err := 

		// get session from redis

		return nil
	}
}
