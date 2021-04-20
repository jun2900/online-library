package controllers

import (
	"github.com/go-playground/validator"
)

//Error response for request input
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func HandlingInput(input interface{}) []*ErrorResponse {
	var inputErrors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			inputErrors = append(inputErrors, &element)
		}
	}
	return inputErrors
}
