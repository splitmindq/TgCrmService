package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/ttacon/libphonenumber"
	"regexp"
	"strings"
)

var (
	Validate      = validator.New()
	phoneRegex    = regexp.MustCompile(`^\+?[0-9]{7,}$`)
	defaultRegion = "RU"
)

func init() {
	err := Validate.RegisterValidation("phone", validateContact)
	if err != nil {
		panic("failed to register contact validator: " + err.Error())
	}
}

func validateContact(fl validator.FieldLevel) bool {
	contact := strings.TrimSpace(fl.Field().String())
	if contact == "" {
		return false
	}

	return isValidPhoneNumber(contact)
}

func isValidPhoneNumber(phone string) bool {

	if !phoneRegex.MatchString(phone) {
		return false
	}
	if !strings.HasPrefix(phone, "+") {
		num, err := libphonenumber.Parse(phone, defaultRegion)
		if err == nil && libphonenumber.IsValidNumber(num) {
			return true
		}
	}

	num, err := libphonenumber.Parse(phone, "")
	if err != nil {
		return false
	}
	return libphonenumber.IsValidNumber(num)
}
