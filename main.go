package main

import (
	"fmt"

	s "github.com/verywelloo/3-go-echo-task-management/app/services"

	"github.com/labstack/echo/v4"
)

var port = s.GetEnv("SEVER_PORT", "5004")

func main() {
	e := echo.New()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.GetEnv("DB_HOST", port))))
}
