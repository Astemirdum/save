package handler

import (
	validator "github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

func (val *Validator) Validate(i interface{}) error {
	return val.validator.Struct(i)
}
