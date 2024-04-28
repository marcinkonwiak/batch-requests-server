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
	_ = v.validator.RegisterValidation("relativeUrl", isRelativeUrl)
	return v
}

func (rq *Validator) Validate(i interface{}) error {
	if err := rq.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return nil
}

func isRelativeUrl(fl validator.FieldLevel) bool {
	if fl.Field().String()[0] != '/' {
		return false
	}

	u, err := url.Parse(fl.Field().String())
	if err != nil {
		return false
	}

	return u.Scheme == "" && u.Host == ""
}
