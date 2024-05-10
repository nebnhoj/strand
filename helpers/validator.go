package helpers

import (
	"github.com/go-playground/validator"
	"github.com/iancoleman/strcase"
)

var Validate = validator.New()

type ErrorResponse struct {
	Error       bool        `json:"error"`
	FailedField string      `json:"field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value" value:"error"`
}

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Validator(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := Validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse
			elem.FailedField = strcase.ToSnake(err.Field()) // Export struct field name
			elem.Tag = err.Tag()                            // Export struct tag
			elem.Value = err.Value()                        // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
