package export

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FixedLengthFormatter struct {
	modelType  reflect.Type
	formatCols map[int]*FixedLength
}
type FixedLength struct {
	Format string
	Length int
}
func GetIndexes(modelType reflect.Type, tagName string) (map[int]*FixedLength, error) {
	ma := make(map[int]*FixedLength, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		tagLength := field.Tag.Get("length")
		if len(tagLength) > 0 {
			length, err := strconv.Atoi(tagLength)
			if err != nil || length < 0 {
				return ma, err
			}
			v := &FixedLength{Length: length}
			if len(tagValue) > 0 {
				if strings.Contains(tagValue, "dateFormat:") {
					tagValue = strings.ReplaceAll(tagValue, "dateFormat:", "")
				}
				v.Format = tagValue
			}
			ma[i] = v
		}
	}
	return ma, nil
}

func NewFixedLengthFormatter(modelType reflect.Type) (*FixedLengthFormatter, error) {
	formatCols, err := GetIndexes(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &FixedLengthFormatter{modelType: modelType, formatCols: formatCols}, nil
}

func (f *FixedLengthFormatter) Format(ctx context.Context, model interface{}) string {
	return ToFixedLength(model, f.formatCols)
}
func ToFixedLength(model interface{}, formatCols map[int]*FixedLength) string {
	arr := make([]string, 0)
	sumValue := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < sumValue.NumField(); i++ {
		format, ok := formatCols[i]
		if ok {
			field := sumValue.Field(i)
			kind := field.Kind()
			var value string
			if kind == reflect.Ptr && field.IsNil() {
				value = FixedLengthString(format.Length, value)
			} else {
				v := field.Interface()
				if kind == reflect.Ptr {
					v = reflect.Indirect(reflect.ValueOf(v)).Interface()
				}
				d, okD := v.(*time.Time)
				if okD {
					value = d.Format(format.Format)
				} else {
					value = fmt.Sprint(v)
					if len(value) > format.Length {
						value = strings.TrimSpace(value)
					}
					if len(format.Format) > 0 {
						value = fmt.Sprintf(format.Format, value)
					}
					value = FixedLengthString(format.Length, value)
				}
			}
			arr = append(arr, value)
		}
	}
	return strings.Join(arr, "") + "\n"
}
func FixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}
