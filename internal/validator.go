package internal

import (
	validator "gopkg.in/go-playground/validator.v9"
)

// CustomValidator is struct for custom validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate is method implementation for validating struct
func (cv CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewValidator is function to init custom validator
func NewValidator() CustomValidator {
	return CustomValidator{validator: validator.New()}
}
