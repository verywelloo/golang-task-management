package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
				fmt.Printf("== REQ : Time : %v, URI : %v, Method : %v, Status : %v, Host : %v, Res Time : %v\n\n", v.StartTime, v.URI, v.Method, v.Status, v.Host, v.Latency)
			} else {
				// Latency = duration time
				fmt.Printf("!! !! ERR : Time : %v, URI : %v, Method : %v, Status : %v, Host : %v, Err : %v, Res Time : %v\n\n", v.StartTime, v.Method, v.URI, v.Status, v.Host, v.Error, v.Latency)
			}
			return nil
		},
	}))

	// setup gracefully shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// routes
	routes.ApiRouter(e)

	// start server in goroutine
	go func() {
		// start server
		fmt.Printf("Server starting on prot %s...\n", port)
		if err := e.Start(":" + port); err != nil {
			e.Logger.Fatal("Server error: ", err)
		}
	}()

	// wait for shutdown signal
	ctx.Done()
	fmt.Println("!! Shutdown signal received, shutting down server...")

	// gracefully shutdown down
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctxShutdown); err != nil {
		e.Logger.Fatal("Cannot shutdown, err: ", err)
	}

	//e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.GetEnv("DB_HOST", port))))
}
