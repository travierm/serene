package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterDatabaseRoutes(e *echo.Echo) {
	e.GET("/databases", GetDatabases)
	e.GET("/users/:id", GetDatabaseById)
	e.POST("/users", CreateDatabase)
}

func GetDatabases(c echo.Context) error {
	return c.JSON(http.StatusOK, databases)
}

func GetDatabaseById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, user := range databases {
		if user.ID == id {
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func CreateDatabase(c echo.Context) error {
	user := new(Database)
	if err := c.Bind(user); err != nil {
		return err
	}
	user.ID = len(databases) + 1
	databases = append(databases, *user)
	return c.JSON(http.StatusCreated, user)
}
