package service

import (
	"context"
	"fmt"
)

type ErrorChecker struct {
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error)
}

func NewErrorChecker(validate func(context.Context, interface{}) ([]ErrorMessage, error)) *ErrorChecker {
	return &ErrorChecker{validate: validate}
}

func (v *ErrorChecker) Check(ctx context.Context, model interface{}) error {
	errors, err := v.validate(ctx, model)
	if err != nil {
		return err
	}
	if errors != nil && len(errors) > 0 {
		m := fmt.Sprintf("Data: %s . Error: %s", model, errors)
		return fmt.Errorf(m)
	}
	return nil
}
