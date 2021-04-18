package validator

import (
	"context"
	"fmt"
	"reflect"
	"unicode"

	s "github.com/common-go/service"
	"gopkg.in/go-playground/validator.v9"
)

const (
	method = "method"
	patch  = "patch"
)

type DefaultValidator struct {
	validate           *validator.Validate
	CustomValidateList []CustomValidate
}

func NewValidator() *DefaultValidator {
	list := GetCustomValidateList()
	return &DefaultValidator{CustomValidateList: list}
}

func (p *DefaultValidator) Validate(ctx context.Context, model interface{}) ([]s.ErrorMessage, error) {
	errors := make([]s.ErrorMessage, 0)
	if p.validate == nil {
		validate := validator.New()
		validate = p.RegisterCustomValidate(validate)
		p.validate = validate
	}
	err := p.validate.Struct(model)

	if err != nil {
		errors, err = MapErrors(err)
	}
	v := ctx.Value(method)
	if v != nil {
		v2, ok := v.(string)
		if ok {
			if v2 == patch {
				errs := s.RemoveRequiredError(errors)
				return errs, nil
			}
		}
	}
	return errors, err
}

var alias = map[string]string{
	"max":      "maxlength",
	"min":      "minlength",
	"gtefield": "minfield",
	"ltefield": "maxfield",
}

func MapErrors(err error) (list []s.ErrorMessage, err1 error) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		err1 = fmt.Errorf("InvalidValidationError")
		return
	}
	for _, err := range err.(validator.ValidationErrors) {
		code := formatCodeMsg(err)
		list = append(list, s.ErrorMessage{Field: s.FormatErrorField(err.Namespace()), Code: code})
	}
	return
}

func formatCodeMsg(err validator.FieldError) string {
	var code string
	if aliasTag, ok := alias[err.Tag()]; ok {
		if (err.Tag() == "max" || err.Tag() == "min") && err.Kind() != reflect.String {
			code = err.Tag()
		} else {
			code = aliasTag
		}
	} else {
		code = err.Tag()
	}
	if err.Param() != "" {
		code += ":" + lcFirstChar(err.Param())
	}
	return code
}
func lcFirstChar(s string) string {
	if len(s) > 0 {
		runes := []rune(s)
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}
	return s
}
func (p *DefaultValidator) RegisterCustomValidate(validate *validator.Validate) *validator.Validate {
	for _, v := range p.CustomValidateList {
		validate.RegisterValidation(v.Tag, v.Fn)
	}
	return validate
}
