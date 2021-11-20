package impt

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewFixedLengthFormatter(modelType reflect.Type) (*FixedLengthFormatter, error) {
	formatCols, err := GetIndexes(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &FixedLengthFormatter{modelType: modelType, formatCols: formatCols}, nil
}

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
				v.Format = tagValue
			}
			ma[i] = v
		}
	}
	return ma, nil
}
func (f FixedLengthFormatter) ToStruct(ctx context.Context, lines []string) (interface{}, error) {
	line := strings.Join(lines, ``)
	record := reflect.New(f.modelType).Interface()
	err := ScanLineFixLength(line, f.modelType, record, f.formatCols)
	if err != nil {
		return nil, err
	}
	if record != nil {
		return reflect.Indirect(reflect.ValueOf(record)).Interface(), nil
	}
	return record, err
}

func ScanLineFixLength(line string, modelType reflect.Type, record interface{}, formatCols map[int]*FixedLength) error {
	s := reflect.Indirect(reflect.ValueOf(record))
	numFields := modelType.NumField()
	start := 0
	size := len(line)
	for j := 0; j < numFields; j++ {
		field := modelType.Field(j)
		format, ok := formatCols[j]
		if ok && format.Length > 0 {
			end := start + format.Length
			if end > size {
				return errors.New(fmt.Sprintf("scanLineFixLength - exceed range max size . Field name = %v , line = %v ", field.Name, line))
			}
			value := strings.TrimSpace(line[start:end])
			f := s.Field(j)
			if f.IsValid() {
				if f.CanSet() {
					typef := field.Type.String()
					switch typef {
					case "string", "*string":
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(&value))
						} else {
							f.SetString(value)
						}
					case "int64", "*int64":
						value, _ := strconv.ParseInt(value, 10, 64)
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(&value))
						} else {
							f.SetInt(value)
						}
					case "int", "*int":
						value, _ := strconv.Atoi(value)
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(&value))
						} else {
							f.Set(reflect.ValueOf(value))
						}
					case "bool", "*bool":
						boolValue, _ := strconv.ParseBool(value)
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(&boolValue))
						} else {
							f.SetBool(boolValue)
						}
					case "float64", "*float64":
						floatValue, _ := strconv.ParseFloat(value, 64)
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(&floatValue))
						} else {
							f.SetFloat(floatValue)
						}
					case "time.Time", "*time.Time":
						if format, ok := formatCols[j]; ok {
							var fieldDate time.Time
							var err error
							if len(format.Format) > 0 {
							} else {
								fieldDate, err = time.Parse(DateLayout, line)
							}
							if err != nil {
								return err
							}
							if f.Kind() == reflect.Ptr {
								f.Set(reflect.ValueOf(&fieldDate))
							} else {
								f.Set(reflect.ValueOf(fieldDate))
							}
						}
					}
				}
			}
			start = end
		}
	}
	return nil
}
