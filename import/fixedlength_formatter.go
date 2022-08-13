package impt

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
	Scale  int
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
	return ma, nil
}
func (f FixedLengthFormatter) ToStruct(ctx context.Context, line string, res interface{}) (error) {
	err := ScanLineFixLength(line, res, f.formatCols)
	if err != nil {
		return err
	}
	return err
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
								fieldDate, err = time.Parse(format.Format, value)
							} else {
								fieldDate, err = time.Parse(DateLayout, value)
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
						if formatf, ok := formatCols[j]; ok {
							bf := new(big.Float)
							if bfv, ok1 := bf.SetString(value); ok1 {
								if formatf.Scale >= 0 && bfv != nil {
									k := Round(*bf, formatf.Scale)
									bf = &k
								}
								if f.Kind() == reflect.Ptr {
									f.Set(reflect.ValueOf(bf))
								} else {
									if bf != nil {
										f.Set(reflect.ValueOf(*bf))
									}
								}
							}

						}
					case "big.Int", "*big.Int":
						if _, ok := formatCols[j]; ok {
							bf := new(big.Int)
							if bfv, oki := bf.SetString(value, 10); oki {
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
