package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marcinkonwiak/batch-requests-server/handler"
	"github.com/marcinkonwiak/batch-requests-server/validator"
)

func NewServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = validator.NewValidator()

	r := e.Group("")
	h := handler.NewHandler()
	h.RegisterRoutes(r)

	return e
}
