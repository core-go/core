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

const (
	DateLayout string = "2006-01-02 15:04:05 +0700 +07"
)

type DelimiterFormatter struct {
	Delimiter  string
	modelType  reflect.Type
	formatCols map[int]Delimiter
}

type Delimiter struct {
	Format string
	Scale  int
}

func NewDelimiterFormatter(modelType reflect.Type, opts ...string) (*DelimiterFormatter, error) {
	sep := "|"
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
func ToTextWithDelimiter(ctx context.Context, model interface{}, delimiter string, formatCols map[int]Delimiter) string {
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
				v := field.Interface()
				if kind == reflect.Ptr {
					v = reflect.Indirect(reflect.ValueOf(v)).Interface()
				}
				if s, okS := v.(string); okS {
					if strings.Contains(value, `"`) || strings.Contains(s, delimiter) {
						s = strings.ReplaceAll(s, `"`, `""`)
						value = "\"" + s + "\""
					} else {
						value = s
					}
				} else if d, okD := v.(time.Time); okD {
					if len(format.Format) > 0 {
						value = d.Format(format.Format)
					} else {
						value = d.Format(DateLayout)
					}
				} else {
					kind = reflect.Indirect(field).Kind()
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
								if format.Scale >= 0 {
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
						}
					} else {
						value = fmt.Sprint(v)
					}
				}
			}
			arr = append(arr, value)
		}
	}
	return strings.Join(arr, delimiter) + "\n"
}
func GetIndexesByTag(modelType reflect.Type, tagName string, skipTag string) (map[int]Delimiter, error) {
	ma := make(map[int]Delimiter)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		if tagValue != "-" {
			skipValue := field.Tag.Get(skipTag)
			v := Delimiter{}
			tagScale, sOk := field.Tag.Lookup("scale")
			if sOk {
				scale, err := strconv.Atoi(tagScale)
				if err == nil {
					v.Scale = scale
				}
			}
			if len(skipValue) > 0 {
				if len(tagValue) > 0 {
					if strings.Contains(tagValue, "dateFormat:") {
						tagValue = strings.ReplaceAll(tagValue, "dateFormat:", "")
						v.Format = tagValue
					}
				}
			} else {
				if len(tagValue) > 0 {
					if strings.Contains(tagValue, "dateFormat:") {
						tagValue = strings.ReplaceAll(tagValue, "dateFormat:", "")
						v.Format = tagValue
					} else if sOk == false && strings.Contains(tagValue, "scale:") {
						tagValue = strings.ReplaceAll(tagValue, "scale:", "")
						scale, err1 := strconv.Atoi(tagValue)
						if err1 != nil {
							return ma, err1
						}
						v.Scale = scale
					}
				} else {

				}
			}
			ma[i] = v
		}
	}
	return ma, nil
}
