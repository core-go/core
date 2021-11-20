package export

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	DateLayout string = "2006-01-02 15:04:05 +0700 +07"
)

type DelimiterFormatter struct {
	Delimiter  string
	modelType  reflect.Type
	formatCols map[int]string
}

func NewDelimiterFormatter(modelType reflect.Type, opts ...string) (*DelimiterFormatter, error) {
	sep := ","
	if len(opts) > 0 && len(opts[0]) > 0 {
		sep = opts[0]
	}
	skipTag := ""
	if len(opts) > 1 && len(opts[1]) > 0 {
		skipTag = opts[1]
	}
	formatCols, err := GetIndexesByTag(modelType, "format", skipTag)
	if err != nil {
		return nil, err
	}
	return &DelimiterFormatter{modelType: modelType, formatCols: formatCols, Delimiter: sep}, nil
}

func (f *DelimiterFormatter) Format(ctx context.Context, model interface{}) string {
	return ToTextWithDelimiter(ctx, model, f.Delimiter, f.formatCols)
}
func ToTextWithDelimiter(ctx context.Context, model interface{}, delimiter string, formatCols map[int]string) string {
	arr := make([]string, 0)
	sumValue := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < sumValue.NumField(); i++ {
		format, ok := formatCols[i]
		if ok {
			field := sumValue.Field(i)
			kind := field.Kind()
			var value string
			if kind == reflect.Ptr && field.IsNil() {
				value = ""
			} else {
				value = fmt.Sprint(field.Interface())
				v := field.Interface()
				if kind == reflect.Ptr {
					v = reflect.Indirect(reflect.ValueOf(v)).Interface()
				}
				d, okD := v.(time.Time)
				if okD {
					if len(format) > 0 {
						value = d.Format(format)
					} else {
						value = d.Format(DateLayout)
					}
				} else {
					s, okS := v.(string)
					if okS {
						if strings.Contains(value, `"`) || strings.Contains(s, delimiter) {
							s = strings.ReplaceAll(s, `"`, `""`)
							value = "\"" + s + "\""
						}

					}
				}
			}
			arr = append(arr, value)
		}
	}
	return strings.Join(arr, delimiter) + "\n"
}
func GetIndexesByTag(modelType reflect.Type, tagName string, skipTag string) (map[int]string, error) {
	ma := make(map[int]string, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		skipValue := field.Tag.Get(skipTag)
		if len(skipValue) > 0 {
			if len(tagValue) > 0 {
				if strings.Contains(tagValue, "dateFormat:") {
					tagValue = strings.ReplaceAll(tagValue, "dateFormat:", "")
				}
				ma[i] = tagValue
			}
		} else {
			if len(tagValue) > 0 {
				if strings.Contains(tagValue, "dateFormat:") {
					tagValue = strings.ReplaceAll(tagValue, "dateFormat:", "")
				}
				ma[i] = tagValue
			} else {
				ma[i] = ""
			}
		}
	}
	return ma, nil
}
