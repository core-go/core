package validator

import (
	"context"
	s "github.com/core-go/core"
	v "github.com/core-go/core/v10"
)

type Validate[T any] func(ctx context.Context, model T) ([]s.ErrorMessage, error)

type Validator[T any] struct {
	validator *v.Validator
}

func NewChecker[T any](opts ...bool) (*Validator[T], error) {
	return NewValidatorWithMap[T](nil, opts...)
}
func NewCheckerWithMap[T any](mp map[string]string, opts ...bool) (*Validator[T], error) {
	return NewValidatorWithMap[T](mp, opts...)
}
func NewValidator[T any](opts ...bool) (*Validator[T], error) {
	return NewValidatorWithMap[T](nil, opts...)
}
func NewValidatorWithMap[T any](mp map[string]string, opts ...bool) (*Validator[T], error) {
	val, err := v.NewCheckerWithMap(mp, opts...)
	if err != nil {
		return nil, err
	}
	return &Validator[T]{val}, nil
}
func (p *Validator[T]) Check(ctx context.Context, model T) ([]s.ErrorDetail, error) {
	errs, err := p.Validate(ctx, model)
	errors := s.BuildErrorDetails(errs, p.validator.IgnoreField)
	return errors, err
}
func (p *Validator[T]) Validate(ctx context.Context, model T) ([]s.ErrorMessage, error) {
	return p.validator.Validate(ctx, model)
}
