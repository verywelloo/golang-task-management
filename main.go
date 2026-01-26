package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
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

	// for cors
	allowOrigin := os.Getenv("ALLOW_ORIGIN")
	if allowOrigin == "" {
		log.Fatal("XX ALLOW_ORIGIN environment variable must be set")
	}

	// global middleware
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 16 << 10, // stack size = บันทึกประวัติการทำงานของโค้ดก่อนที่จะพัง
		LogErrorFunc: func(c echo.Context, err error, _ []byte) error {
			file, line, fn := topAppFrame()
			method := c.Request().Method
			path := c.Request().URL.String()
			reqID := c.Response().Header().Get(echo.HeaderXRequestID) // request id from header

			c.Logger().Errorf("[PANIC] %s %s req_id=%s at %s:%d (%s): %v", method, path, reqID, file, line, fn, err)

			return err
		},
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{allowOrigin}, // don't set *, cause cors will block download request
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Disposition"}, // allow frontend to read file name
		AllowCredentials: true,                            // cookies, session
	}))

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
}

func topAppFrame() (file string, line int, function string) {
	const maxDepth = 64
	pcs := make([]uintptr, maxDepth)
	// Skip frames: Callers, topAppFrame, LogErrorFunc wrapper, recover trampoline
	n := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		// f = current, more = next
		f, more := frames.Next()
		fn := f.Function
		if !strings.HasPrefix(fn, "runtime.") && !strings.HasPrefix(fn, "net/http.") &&
			!strings.Contains(fn, "github.com/labstack/echo") {
			return f.File, f.Line, fn
		}

		// if not more, break
		if !more {
			break
		}
	}
	return "?", 0, "?"
}
