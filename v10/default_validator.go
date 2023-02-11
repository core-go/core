package validator

import (
	"context"
	"fmt"
	"reflect"
	"unicode"

	ut "github.com/go-playground/universal-translator"

	s "github.com/core-go/core"
	"github.com/go-playground/validator/v10"
)

const (
	method = "method"
	patch  = "patch"
)

type DefaultValidator struct {
	validate           *validator.Validate
	Trans              *ut.Translator
	CustomValidateList []CustomValidate
}
func NewChecker(opts ...bool) (*DefaultValidator, error) {
	return NewValidator(opts...)
}
func NewValidator(opts ...bool) (*DefaultValidator, error) {
	register := true
	if len(opts) > 0 {
		register = opts[0]
	}
	uValidate, uTranslator, err := NewDefaultValidator()
	if err != nil {
		return nil, err
	}
	list := GetCustomValidateList()
	validator := &DefaultValidator{validate: uValidate, Trans: &uTranslator, CustomValidateList: list}
	if register {
		err2 := validator.RegisterCustomValidate()
		if err2 != nil {
			return validator, err2
		}
	}
	return validator, nil
}
func NewDefaultChecker() (*validator.Validate, ut.Translator, error) {
	return NewDefaultValidator()
}
func NewDefaultValidator() (*validator.Validate, ut.Translator, error) {
	validate := validator.New()
	var transl ut.Translator
	if trans != nil {
		transl = *trans
	} else {
		list := GetCustomValidateList()
		for _, v := range list {
			err := validate.RegisterValidation(v.Tag, v.Fn)
			if err != nil {
				return nil, nil, err
			}
		}
		ptr, err := RegisterTranslatorEn(validate)
		if err != nil {
			return nil, nil, err
		}
		transl = ptr
	}
	return validate, transl, nil
}
func (p *DefaultValidator) Check(ctx context.Context, model interface{}) ([]s.ErrorDetail, error) {
	errs, err := p.Validate(ctx, model)
	errors := s.BuildErrorDetails(errs)
	return errors, err
}
func (p *DefaultValidator) Validate(ctx context.Context, model interface{}) ([]s.ErrorMessage, error) {
	errors := make([]s.ErrorMessage, 0)
	err := p.validate.Struct(model)

	if err != nil {
		errors, err = p.MapErrors(err)
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

func getTagName(err validator.FieldError) string {
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
func (p *DefaultValidator) RegisterCustomValidate() error {
	for _, v := range p.CustomValidateList {
		err := p.validate.RegisterValidation(v.Tag, v.Fn)
		if err != nil {
			return err
		}
	}
	if p.Trans != nil && p.validate != nil {
		// register default translate
		for _, validate := range p.CustomValidateList {
			if text, ok := translations[validate.Tag]; ok {
				err := AddMessage(p.validate, *p.Trans, validate.Tag, text, true)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *DefaultValidator) MapErrors(err error) (list []s.ErrorMessage, err1 error) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		err1 = fmt.Errorf("InvalidValidationError")
		return
	}
	tr := *p.Trans
	for _, err := range err.(validator.ValidationErrors) {
		code := getTagName(err)
		list = append(list, s.ErrorMessage{Field: s.FormatErrorField(err.Namespace()), Code: code, Message: err.Translate(tr), Param: err.Param()})
	}
	return
}
