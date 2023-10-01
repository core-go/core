package validator

import sv "github.com/core-go/core"

func NewErrorChecker() *sv.ErrorChecker {
	v, _ := NewValidatorWithMap()
	return sv.NewErrorChecker(v.Validate)
}
