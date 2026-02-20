package middleware

import (
	"fmt"
	"net/http"
	"strings"

	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	res "github.com/verywelloo/3-go-echo-task-management/app/dto/response"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"

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

		// decode
		accessToken := tokenParts[1]
		claims, err := s.DecodeAccessToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "invalid token",
			})
		}

		// get session from redis
		sessionKey := fmt.Sprintf("session:%s", claims.ID)
		var session req.CacheSession
		if err := s.GetRedis(c, sessionKey, &session); err != nil {
			return c.JSON(http.StatusUnauthorized, res.Result{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
				Details: err.Error(),
			})
		}

		//validate session
		if session.ID != claims.ID || session.Ip != c.RealIP() || session.Agent != c.Request().UserAgent() {
			return c.JSON(htt)
		}
		return nil
	}
}
