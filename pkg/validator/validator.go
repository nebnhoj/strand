package validator

import (
	"github.com/go-playground/validator"
	"github.com/iancoleman/strcase"
)

var v = validator.New()

// FieldError describes a single validation failure.
type FieldError struct {
	Error bool        `json:"error"`
	Field string      `json:"field"`
	Tag   string      `json:"tag"`
	Value interface{} `json:"value"`
}

// Validate runs struct validation and returns all field errors.
func Validate(data interface{}) []FieldError {
	var errs []FieldError
	if err := v.Struct(data); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, FieldError{
				Error: true,
				Field: strcase.ToSnake(e.Field()),
				Tag:   e.Tag(),
				Value: e.Value(),
			})
		}
	}
	return errs
}
