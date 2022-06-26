package core

import (
	"context"
	"strings"
	"unicode"
)

type ResultInfo struct {
	Status  int            `mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	Errors  []ErrorMessage `mapstructure:"errors" json:"errors,omitempty" gorm:"column:errors" bson:"errors,omitempty" dynamodbav:"errors,omitempty" firestore:"errors,omitempty"`
	Value   interface{}    `mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Message string         `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

type ErrorMessage struct {
	Field   string `mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
type Validator interface {
	Validate(ctx context.Context, model interface{}) ([]ErrorMessage, error)
}
type MapValidator interface {
	Validate(ctx context.Context, model map[string]interface{}) ([]ErrorMessage, error)
}

func RemoveRequiredError(errors []ErrorMessage) []ErrorMessage {
	if errors == nil || len(errors) == 0 {
		return errors
	}
	errs := make([]ErrorMessage, 0)
	for _, s := range errors {
		if s.Code != "required" && !strings.HasPrefix(s.Code, "minlength") {
			errs = append(errs, s)
		} else if strings.Index(s.Field, ".") >= 0 {
			errs = append(errs, s)
		}
	}
	return errs
}
func FormatErrorField(s string) string {
	splitField := strings.Split(s, ".")
	length := len(splitField)
	if length == 1 {
		return lcFirstChar(splitField[0])
	} else if length > 1 {
		var tmp []string
		for _, v := range splitField[1:] {
			tmp = append(tmp, lcFirstChar(v))
		}
		return strings.Join(tmp, ".")
	}
	return s
}
func lcFirstChar(s string) string {
	if len(s) > 0 {
		runes := []rune(s)
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}
	return s
}
