package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s")))
}
