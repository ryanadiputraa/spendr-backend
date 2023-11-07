package validator

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(val any) (error, map[string]string)
}

type validation struct {
	validator *validator.Validate
}

func NewValidator() Validator {
	return &validation{
		validator: validator.New(),
	}
}

func (v *validation) Validate(val any) (error, map[string]string) {
	err := v.validator.Struct(val)
	errorsMap := make(map[string]string)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldErr := range validationErrors {
			field := fieldToSnakeCase(fieldErr.Field())
			errorsMap[field] = FieldErrMsg(fieldErr)

		}
		return errors.New("invalid params"), errorsMap
	}

	return nil, nil
}

func FieldErrMsg(err validator.FieldError) string {
	field := fieldToSnakeCase(err.Field())
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "max":
		return fmt.Sprintf("%s should have a maximum length of %s", field, err.Param())
	case "min":
		return fmt.Sprintf("%s should have a minimum length of %s", field, err.Param())
	case "email":
		return fmt.Sprintf("%s should be a valid email address", field)
	case "http_url":
		return fmt.Sprintf("%s should be a valid http url", field)
	default:
		return err.Error()
	}
}

func fieldToSnakeCase(input string) string {
	var result []rune
	for i, char := range input {
		if i > 0 && unicode.IsUpper(char) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(char))
	}
	return string(result)
}
