package validator

import sv "github.com/core-go/core/import"

func NewErrorChecker() *sv.ErrorChecker {
	v := NewValidator()
	return sv.NewErrorChecker(v.Validate)
}
