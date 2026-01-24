package main

import (
	"fmt"

	s "github.com/verywelloo/3-go-echo-task-management/app/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var port = s.GetEnv("SEVER_PORT", "5004")

func main() {
	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogHost:     true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				fmt.Printf("== REQ : Time : %v, URI : %v, Method : %v, Status : %v, Host : %v\n\n", v.StartTime, v.URI, v.Method, v.Status, v.Host, v.Error)
			} else {
				fmt.Printf("!! !! ERR : Time : %v, URI : %v, Method : %v, Status : %v, Host : %v, Err : %v\n\n", v.StartTime, v.Method, v.URI, v.Status, v.Host, v.Error)
			}
			return nil
		},
	}))
	//e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.GetEnv("DB_HOST", port))))
}
