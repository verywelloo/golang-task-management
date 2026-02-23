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
	matched, _ := regexp.MatchString(
		`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$`,
		fl.Field().String(),
	)
	return matched
}
