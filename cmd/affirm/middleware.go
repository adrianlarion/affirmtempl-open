package main

import (
	"net/http"

	"github.com/adrianlarion/affirmtempl-open/internal/model"
	"github.com/labstack/echo/v4"
)

func checkAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if model.UserIsAuthenticated(c) {
			return next(c)
		} else {
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}
