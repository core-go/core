package writer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
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
	Scale  int
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
			if length > 0 {
				v := &FixedLength{Length: length}
				tagScale, sOk := field.Tag.Lookup("scale")
				if sOk {
					scale, err := strconv.Atoi(tagScale)
					if err == nil {
						v.Scale = scale
					}
				}
				if len(tagValue) > 0 {
					if strings.Contains(tagValue, "dateFormat:") {
						tagValue = strings.ReplaceAll(tagValue, "dateFormat:", "")
					} else if sOk == false && strings.Contains(tagValue, "scale:") {
						tagValue = strings.ReplaceAll(tagValue, "scale:", "")
						scale, err1 := strconv.Atoi(tagValue)
						if err1 != nil {
							return ma, err1
						}
						v.Scale = scale
					}
					v.Format = tagValue
				}
				ma[i] = v
			}
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
					kind = field.Elem().Kind()
				}
				if s, okS := v.(string); okS {
					value = FixedLengthString(format.Length, s)
				} else {
					if d, okD := v.(time.Time); okD {
						value = d.Format(format.Format)
					} else {
						if kind == reflect.Struct {
							if v2 := reflect.Indirect(reflect.ValueOf(v)); v2.NumField() == 1 {
								f := v2.Field(0)
								fv := f.Interface()
								k := f.Kind()
								if k == reflect.Ptr {
									fv = reflect.Indirect(reflect.ValueOf(fv)).Interface()
								}
								if sv, ok := fv.(big.Float); ok {
									prec := 2
									if format.Scale > 0 {
										prec = format.Scale
									}
									value = sv.Text('f', prec)
								} else if svi, ok := fv.(big.Int); ok {
									value = svi.Text(10)
								} else {
									value = fmt.Sprint(v)
								}
							} else {
								value = fmt.Sprint(v)
								if len(value) > format.Length {
									value = strings.TrimSpace(value)
								}
								if len(format.Format) > 0 {
									value = fmt.Sprintf(format.Format, value)
								}
							}
						} else {
							value = fmt.Sprint(v)
							if len(value) > format.Length {
								value = strings.TrimSpace(value)
							}
							if len(format.Format) > 0 {
								value = fmt.Sprintf(format.Format, value)
							}
						}
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
