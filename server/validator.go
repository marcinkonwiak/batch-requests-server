package server

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{validator: validator.New()}
	_ = v.validator.RegisterValidation("relativeUri", isRelativeUri)
	return v
}

func (rq *Validator) Validate(i interface{}) error {
	if err := rq.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return nil
}

func isRelativeUri(fl validator.FieldLevel) bool {
	u, err := url.Parse(fl.Field().String())
	return err == nil && u.IsAbs() == false
}
