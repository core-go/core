package convert

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)
const layout = "2006-01-02"

func ToCamelCase(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	p := make(map[string]string)
	for key, element := range m {
		n := ToCamel(key)
		p[n] = element
	}
	return p
}
func ToCamel(s string) string {
	s2 := strings.ToUpper(s)
	s1 := string(s[0])
	for i := 1; i < len(s); i++ {
		if string(s[i-1]) == "_" {
			s1 = s1[:len(s1)-1]
			s1 += string(s2[i])
		} else {
			s1 += string(s[i])
		}
	}
	return s1
}

func TimeToMilliseconds(time string) (int64, error) {
	var h, m, s int
	_, err := fmt.Sscanf(time, "%02d:%02d:%02d", &h, &m, &s)
	if err != nil {
		return 0, err
	}
	return int64(h * 3600000 + m * 60000 + s * 1000), nil
}
func DateToUnixTime(s string) (int64, error) {
	date, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	return date.Unix() * 1000, nil
}
func DateToUnixNano(s string) (int64, error) {
	date, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	return date.UnixNano(), nil
}
func UnixTime(m int64) string {
	dateUtc := time.Unix(0, m* 1000000)
	return dateUtc.Format("2006-01-02")
}
func MillisecondsToTimeString(milliseconds int) string {
	hourUint := 3600000 //60 * 60 * 1000 = 3600000
	minuteUint := 60000 //60 * 1000 = 60000
	secondUint := 1000
	hour := milliseconds / hourUint
	milliseconds = milliseconds % hourUint
	minute := milliseconds / minuteUint
	milliseconds = milliseconds % minuteUint
	second := milliseconds / secondUint
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}
func StringToAvroDate(date *string) (*int, error) {
	if date == nil {
		return nil, nil
	}
	d, err := time.Parse(layout, *date)
	if err != nil {
		return nil, err
	}
	i := int(d.Unix() / 86400)
	return &i, nil
}
func ToAvroDate(date *time.Time) *int {
	if date == nil {
		return nil
	}
	i := int(date.Unix() / 86400)
	return &i
}
func RoundFloat(num float64, slice int) float64 {
	c := math.Pow10(slice)
	result := math.Ceil(num*c) / c
	return result
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

func Clone(origin interface{}) interface{} {
	originValue := reflect.Indirect(reflect.ValueOf(origin))
	originType := reflect.TypeOf(origin)

	resultType := reflect.TypeOf(origin)
	result := reflect.New(resultType)
	numFields := originType.NumField()
	for i := 0; i < numFields; i++ {
		field := originType.Field(i)
		value := originValue.FieldByName(field.Name)
		f := result.Elem().Field(i)
		if value.Kind() == reflect.String {
			f.SetString(value.String())
		} else if value.Kind() == reflect.Int {
			f.SetInt(value.Int())
		} else if value.Kind() == reflect.Float64 {
			f.SetFloat(value.Float())
		} else if value.Kind() == reflect.Bool {
			f.SetBool(value.Bool())
		} else if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			} else {
				val := value.Interface()
				switch val.(type) {
				case *string:
					strVal, ok := val.(*string)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&strVal)))
					}
				case *int:
					intVal, ok := val.(*int)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&intVal)))
					}
				case *float64:
					floatVal, ok := val.(*float64)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&floatVal)))
					}
				case *bool:
					boolVal, ok := val.(*bool)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&boolVal)))
					}
				}
			}
		} else if value.Kind() == reflect.Struct {
			data := Clone(value.Interface())
			f.Set(reflect.Indirect(reflect.ValueOf(data)))
		}
	}
	return result.Interface()
}
func ToMap(in interface{}, ignoreFields ...string) map[string]interface{} {
	return ToMapWithTag(in, "json", ignoreFields...)
}
func ToMapWithTag(in interface{}, tagName string, ignoreFields ...string) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fv := f.Interface()
		k := f.Kind()
		if k == reflect.Ptr {
			if f.IsNil() {
				continue
			} else {
				fv = reflect.Indirect(reflect.ValueOf(fv)).Interface()
			}
		} else if k == reflect.Slice {
			if f.IsNil() {
				continue
			}
		}
		n := getTag(typ.Field(i), tagName)
		out[n] = fv
	}
	for _, v := range ignoreFields {
		if _, ok := out[v]; ok {
			delete(out, v)
		}
	}
	return out
}
func getTag(fi reflect.StructField, tag string) string {
	if tagv := fi.Tag.Get(tag); tagv != "" {
		arrValue := strings.Split(tagv, ",")
		if len(arrValue) > 0 {
			return arrValue[0]
		} else {
			return tagv
		}
	}
	return fi.Name
}
func ToMapOmitEmpty(model interface{}, checkOmit bool, ignoreFields ...string) map[string]interface{} {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}
	numFields := modelType.NumField()
	fields := make(map[string]interface{})
	for i := 0; i < numFields; i++ {
		tag, ok := modelType.Field(i).Tag.Lookup("json")
		if ok {
			name := strings.Split(tag, ",")
			if checkOmit {
				if !modelValue.Field(i).IsZero() {
					fields[name[0]] = modelValue.Field(i).Interface()
				}
			} else {
				fields[name[0]] = modelValue.Field(i).Interface()
			}

		}
	}
	for _, v := range ignoreFields {
		if _, ok := fields[v]; ok {
			delete(fields, v)
		}
	}
	return fields
}
func ToObject(ms map[string]interface{}, res interface{}) error {
	b, err := json.Marshal(ms)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &res)
}
func Copy(src interface{}, des interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &des)
}
func ValueOf(m interface{}, path string) interface{} {
	arr := strings.Split(path, ".")
	i := 0
	var c interface{}
	c = m
	l1 := len(arr) - 1
	for i < len(arr) {
		key := arr[i]
		m2, ok := c.(map[string]interface{})
		if ok {
			c = m2[key]
		}
		if !ok || i >= l1 {
			return c
		}
		i++
	}
	return c
}
func Merge(m map[string]interface{}, sub map[string]interface{}, opts...bool) map[string]interface{} {
	if m == nil {
		return sub
	}
	if len(opts) > 0 && opts[0] == true {
		for k, v := range sub {
			m[k] = v
		}
	} else {
		for k, v := range sub {
			_, ok := m[k]
			if ok {
				m[k] = v
			}
		}
	}
	return m
}
