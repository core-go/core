package reader

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type FixedLength struct {
	TypeName string
	Format   string
	Length   int
	Scale    int
	Handle   func(f reflect.Value, line string, format string, scale int) error
}

func NewFixedLengthTransformer(modelType reflect.Type) (*FixedLengthTransformer, error) {
	formatCols, err := GetIndexes(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &FixedLengthTransformer{modelType: modelType, formatCols: formatCols}, nil
}

type FixedLengthTransformer struct {
	modelType  reflect.Type
	formatCols map[int]*FixedLength
}

func (f FixedLengthTransformer) Transform(ctx context.Context, line string, res interface{}) error {
	err := ScanLineFixLength(line, res, f.formatCols)
	return err
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
				v.TypeName = field.Type.String()
				fn, ok := funcMap[v.TypeName]
				if ok {
					v.Handle = fn
				} else {
					v.Handle = HandleUnknown
				}
				ma[i] = v
			}
		}
	}
	return ma, nil
}

func ScanLineFixLength(line string, record interface{}, formatCols map[int]*FixedLength) error {
	modelType := reflect.TypeOf(record).Elem()
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
					err := format.Handle(f, value, format.Format, format.Scale)
					if err != nil {
						return err
					}
				}
			}
			start = end
		}
	}
	return nil
}

func Round(num big.Float, scale int) big.Float {
	marshal, _ := num.MarshalText()
	if strings.IndexRune(string(marshal), '.') == -1 {
		return num
	}
	fmt.Println(marshal)
	var dot int
	for i, v := range marshal {
		if v == 46 {
			dot = i + 1
			break
		}
	}
	a := marshal[:dot]
	b := marshal[dot : dot+scale+1]
	c := b[:len(b)-1]

	if b[len(b)-1] >= 53 {
		c[len(c)-1] += 1
	}
	var r []byte
	r = append(r, a...)
	r = append(r, c...)
	num.UnmarshalText(r)
	return num
}
