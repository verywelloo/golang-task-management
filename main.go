package main

import (
	"fmt"

	"github.com/verywelloo/3-go-echo-task-management/app/routes"
	s "github.com/verywelloo/3-go-echo-task-management/app/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var port = s.GetEnv("SEVER_PORT", "5004")

func main() {
	// start echo
	e := echo.New()

	// set log details
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogHost:     true,
		LogError:    true,
		HandleError: true,
		// set custom log
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				// Latency = duration time
				fmt.Printf("== REQ : Time : %v, URI : %v, Method : %v, Status : %v, Host : %v, Res Time : %v\n\n", v.StartTime, v.URI, v.Method, v.Status, v.Host, v.Error, v.Latency)
			} else {
				// Latency = duration time
				fmt.Printf("!! !! ERR : Time : %v, URI : %v, Method : %v, Status : %v, Host : %v, Err : %v, Res Time : %v\n\n", v.StartTime, v.Method, v.URI, v.Status, v.Host, v.Error, v.Latency)
			}
			return nil
		},
	}))

	routes.ApiRouter(e)

	// start server
	fmt.Printf("Server starting on prot %s...\n", port)
	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal("Server error: ", err)
	}
	//e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.GetEnv("DB_HOST", port))))
}
