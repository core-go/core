package validator

import (
	sv "github.com/common-go/service"
)

func NewErrorChecker() *sv.ErrorChecker {
	v := NewValidator()
	return sv.NewErrorChecker(v.Validate)
}
