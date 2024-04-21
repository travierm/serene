package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	RegisterDatabaseRoutes(e)

	e.Logger.Fatal(e.Start(":2700"))
}