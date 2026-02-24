package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator validator.Validate
}

// method
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// set custom validator
func NewCustomValidator() *CustomValidator {
	v := *validator.New()

	// register custom validator
	v.RegisterValidation("strongpass", strongPassword)
	return &CustomValidator{validator: v}
}

// set string password validation
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check each condition separately
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLower || !hasUpper || !hasDigit {
		return false
	}

	return true
}
