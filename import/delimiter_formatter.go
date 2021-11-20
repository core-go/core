package impt

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	DateLayout string = "2006-01-02 15:04:05 +0700 +07"
)

func GetIndexesByTag(modelType reflect.Type, tagName string) (map[int]Delimiter, error) {
	ma := make(map[int]Delimiter, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		v := Delimiter{}
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
	return ma, nil
}
func NewDelimiterFormatter(modelType reflect.Type) (*DelimiterFormatter, error) {
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &DelimiterFormatter{modelType: modelType, formatCols: formatCols}, nil
}

type DelimiterFormatter struct {
	modelType  reflect.Type
	formatCols map[int]Delimiter
}

type Delimiter struct {
	Format string
	Scale  int
}

func (f DelimiterFormatter) ToStruct(ctx context.Context, lines []string) (interface{}, error) {
	record := reflect.New(f.modelType).Interface()
	err := ScanLine(lines, f.modelType, record, f.formatCols)
	if err != nil {
		return nil, err
	}
	if record != nil {
		return reflect.Indirect(reflect.ValueOf(record)).Interface(), nil
	}
	return record, err
}

func ScanLine(lines []string, modelType reflect.Type, record interface{}, formatCols map[int]Delimiter) error {
	s := reflect.Indirect(reflect.ValueOf(record))
	numFields := s.NumField()
	for i := 0; i < numFields; i++ {
		field := modelType.Field(i)
		typef := field.Type.String()
		line := lines[i]
		f := s.Field(i)
		if f.CanSet() {
			switch typef {
			case "string", "*string":
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&line))
				} else {
					f.SetString(line)
				}
			case "int64", "*int64":
				value, _ := strconv.ParseInt(line, 10, 64)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&value))
				} else {
					f.SetInt(value)
				}
			case "int", "*int":
				value, _ := strconv.Atoi(line)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&value))
				} else {
					f.Set(reflect.ValueOf(value))
				}
			case "bool", "*bool":
				boolValue, _ := strconv.ParseBool(line)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&boolValue))
				} else {
					f.SetBool(boolValue)
				}
			case "float64", "*float64":
				floatValue, _ := strconv.ParseFloat(line, 64)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&floatValue))
				} else {
					f.SetFloat(floatValue)
				}
			case "time.Time", "*time.Time":
				if format, ok := formatCols[i]; ok {
					var fieldDate time.Time
					var err error
					if len(format.Format) > 0 {
						fieldDate, err = time.Parse(format.Format, line)
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
			case "big.Float", "*big.Float":
				if formatf, ok := formatCols[i]; ok {
					bf := new(big.Float)
					if bfv, ok := bf.SetString(line);ok{
						if formatf.Scale >= 0 && bfv != nil {
							k := Round(*bf, formatf.Scale)
							bf = &k
						}
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(bfv))
						} else {
							if bfv != nil {
								f.Set(reflect.ValueOf(*bfv))
							}
						}
					}
				}
			case "big.Int", "*big.Int":
				if _, ok := formatCols[i]; ok {
					bf := new(big.Int)
					if bfv, oki := bf.SetString(line, 10); oki {
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(bfv))
						} else {
							if bfv != nil {
								f.Set(reflect.ValueOf(*bfv))
							}
						}
					}
				}
			}
		}
	}
	return nil
}
