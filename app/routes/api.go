package routes

import (
	"github.com/labstack/echo/v4"
	api "github.com/verywelloo/3-go-echo-task-management/app/controllers"
)

func ApiRouter(e *echo.Echo) {
	apiGroup := e.Group("/api")

	v1Group := apiGroup.Group("/v1")

	auth := v1Group.Group("/auth")
	auth.POST("/register", api.Register)

	user := v1Group.Group("/user")
	user.GET("", api.GetAllUser)
}
