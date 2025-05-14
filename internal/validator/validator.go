package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field string, message string) {
	_, exists := v.Errors[field]
	if !exists {
		v.Errors[field] = message
	}
}

func (v *Validator) Check(ok bool, field string, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func IsValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}

func IsValidPassword(password string) bool {
	// Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one digit, and one special character
	if len(password) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// IsValidQuestionType checks if the provided question type is valid.
func IsValidQuestionType(questionType string) bool {
	validTypes := []string{"text", "checkbox", "radio", "scale"}
	for _, validType := range validTypes {
		if questionType == validType {
			return true
		}
	}
	return false
}
func Valid(value string) bool {
	return NotBlank(value) && MinLength(value, 1) && MaxLength(value, 1000)
}

// Validate response
func ValidateResponse(response string) bool {
	return NotBlank(response) && MinLength(response, 1) && MaxLength(response, 1000)
}
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsValidRole(role string) bool {
	validRoles := []string{"coach", "teacher"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// validate credentials
func validateLogCredentials(email, password string) bool {
	return NotBlank(email) && IsValidEmail(email) && NotBlank(password) && IsValidPassword(password)
}
