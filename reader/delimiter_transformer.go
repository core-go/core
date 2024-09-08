package reader

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Delimiter struct {
	TypeName string
	Format   string
	Scale    int
	Handle   func(f reflect.Value, line string, format string, scale int) error
}

func NewDelimiterTransformer(modelType reflect.Type, options ...string) (*DelimiterTransformer, error) {
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		return nil, err
	}
	separator := ""
	if len(options) > 0 {
		separator = options[0]
	} else {
		separator = "|"
	}
	return &DelimiterTransformer{modelType: modelType, formatCols: formatCols, separator: separator}, nil
}

type DelimiterTransformer struct {
	modelType  reflect.Type
	formatCols map[int]Delimiter
	separator  string
}

func (f DelimiterTransformer) Transform(ctx context.Context, lineStr string, res interface{}) error {
	lines := strings.Split(lineStr, f.separator)
	err := ScanLine(lines, res, f.formatCols)
	return err
}

func GetIndexesByTag(modelType reflect.Type, tagName string) (map[int]Delimiter, error) {
	ma := make(map[int]Delimiter)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		if tagValue != "-" {
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
	return ma, nil
}
func Min(n1 int, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}
func ScanLine(lines []string, record interface{}, formatCols map[int]Delimiter) error {
	s := reflect.Indirect(reflect.ValueOf(record))
	numFields := s.NumField()
	l := len(formatCols)
	le := Min(numFields, l)
	for i := 0; i < le; i++ {
		line := lines[i]
		f := s.Field(i)
		if f.CanSet() {
			if format, ok := formatCols[i]; ok {
				err := format.Handle(f, line, format.Format, format.Scale)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
