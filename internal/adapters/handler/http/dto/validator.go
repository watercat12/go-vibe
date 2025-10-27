package dto

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9@!$%^]{12,50}$`)

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 12 || len(password) > 50 {
		return false
	}

	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasDigit := strings.ContainsAny(password, "0123456789")
	hasSpecial := strings.ContainsAny(password, "@!$%^")

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func ValidateTeam(fl validator.FieldLevel) bool {
	team := fl.Field().String()
	validTeams := []string{"FRONT_END", "BACK_END", "QA", "ADMIN", "BRSE", "DESIGN", "OTHERS"}
	for _, validTeam := range validTeams {
		if team == validTeam {
			return true
		}
	}
	return false
}

func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("password", ValidatePassword)
	v.RegisterValidation("team", ValidateTeam)
	v.RegisterValidation("phone", ValidatePhone)
}