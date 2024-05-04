package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) processBindingError(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Badly-formed JSON request body (at position %d)", syntaxError.Offset)
		return echo.NewHTTPError(http.StatusBadRequest, msg)

	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		return echo.NewHTTPError(http.StatusBadRequest, msg)

	default:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
}
