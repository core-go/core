package validator

import sv "github.com/core-go/core"

func NewErrorChecker() *sv.ErrorChecker {
	v := NewValidator()
	return sv.NewErrorChecker(v.Validate)
}
