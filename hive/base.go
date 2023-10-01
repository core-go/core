package hive

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const t0 = "2006-01-02 15:04:05"

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func WrapString(v string) string {
	if strings.Index(v, `'`) >= 0 {
		return join(`'`, strings.Replace(v, "'", "''", -1), `'`)
	}
	return join(`'`, v, `'`)
}

func GetDBValue(v interface{}, scale int8, layoutTime string) (string, bool) {
	switch v.(type) {
	case string:
		s0 := v.(string)
		if len(s0) == 0 {
			return "''", true
		}

		return WrapString(s0), true
	case bool:
		b0 := v.(bool)
		if b0 {
			return "true", true
		} else {
			return "false", true
		}
	case int:
		return strconv.Itoa(v.(int)), true
	case int64:
		return strconv.FormatInt(v.(int64), 10), true
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10), true
	case big.Int:
		var z1 big.Int
		z1 = v.(big.Int)
		return z1.String(), true
	case float64:
		if scale >= 0 {
			mt := "%." + strconv.Itoa(int(scale)) + "f"
			return fmt.Sprintf(mt, v), true
		}
		return fmt.Sprintf("'%f'", v), true
	case time.Time:
		tf := v.(time.Time)
		if len(layoutTime) > 0 {
			f := tf.Format(layoutTime)
			return WrapString(f), true
		}
		f := tf.Format(t0)
		return WrapString(f), true
	case big.Float:
		n1 := v.(big.Float)
		if scale >= 0 {
			n2 := Round(n1, int(scale))
			return fmt.Sprintf("%v", &n2), true
		} else {
			return fmt.Sprintf("%v", &n1), true
		}
	case big.Rat:
		n1 := v.(big.Rat)
		if scale >= 0 {
			return RoundRat(n1, scale), true
		} else {
			return n1.String(), true
		}
	case float32:
		if scale >= 0 {
			mt := "%." + strconv.Itoa(int(scale)) + "f"
			return fmt.Sprintf(mt, v), true
		}
		return fmt.Sprintf("'%f'", v), true
	default:
		if scale >= 0 {
			v2 := reflect.ValueOf(v)
			if v2.Kind() == reflect.Ptr {
				v2 = v2.Elem()
			}
			if v2.NumField() == 1 {
				f := v2.Field(0)
				fv := f.Interface()
				k := f.Kind()
				if k == reflect.Ptr {
					if f.IsNil() {
						return "null", true
					} else {
						fv = reflect.Indirect(reflect.ValueOf(fv)).Interface()
						sv, ok := fv.(big.Float)
						if ok {
							return sv.Text('f', int(scale)), true
						} else {
							return "", false
						}
					}
				} else {
					sv, ok := fv.(big.Float)
					if ok {
						return sv.Text('f', int(scale)), true
					} else {
						return "", false
					}
				}
			} else {
				return "", false
			}
		} else {
			return "", false
		}
	}
	return "", false
}
func Round(num big.Float, scale int) big.Float {
	marshal, _ := num.MarshalText()
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
func RoundRat(rat big.Rat, scale int8) string {
	digits := int(math.Pow(float64(10), float64(scale)))
	floatNumString := rat.RatString()
	sl := strings.Split(floatNumString, "/")
	a := sl[0]
	b := sl[1]
	c, _ := strconv.Atoi(a)
	d, _ := strconv.Atoi(b)
	intNum := c / d
	surplus := c - d*intNum
	e := surplus * digits / d
	r := surplus * digits % d
	if r >= d/2 {
		e += 1
	}
	res := strconv.Itoa(intNum) + "." + strconv.Itoa(e)
	return res
}
