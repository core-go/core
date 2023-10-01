package validator

import sv "github.com/core-go/core"

func NewErrorChecker() (*sv.ErrorChecker, error) {
	v, err := NewValidator()
	if err != nil {
		return nil, err
	}
	return sv.NewErrorChecker(v.Validate), nil
}
