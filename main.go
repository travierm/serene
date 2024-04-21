package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	file, err := os.OpenFile("storage/wal.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening WAL file:", err)
		return
	}
	defer file.Close()

	// Create a new WAL instance
	//wal := &WAL{file: file}

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &SereneContext{Databases: databases}
			return next(cc)
		}
	})

	RegisterDatabaseRoutes(e)

	e.Logger.Fatal(e.Start(":2700"))
}
