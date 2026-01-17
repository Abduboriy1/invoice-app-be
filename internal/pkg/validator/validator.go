package validator

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Validator struct {
	errors []string
}

func New() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, field+" is required")
	}
	return v
}

func (v *Validator) Email(field, value string) *Validator {
	if !emailRegex.MatchString(value) {
		v.errors = append(v.errors, field+" must be a valid email")
	}
	return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors = append(v.errors, field+" must be at least "+string(rune(min+'0'))+" characters")
	}
	return v
}

func (v *Validator) Positive(field string, value float64) *Validator {
	if value <= 0 {
		v.errors = append(v.errors, field+" must be positive")
	}
	return v
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() []string {
	return v.errors
}
