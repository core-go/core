package sql

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type BatchStatement struct {
	Query         string
	Values        []interface{}
	Keys          []string
	Columns       []string
	Attributes    map[string]interface{}
	AttributeKeys map[string]interface{}
}
type FieldDB struct {
	JSON   string
	Column string
	Field  string
	Index  int
	Key    bool
	Update bool
	Insert bool
	Scale  int8
	True   *string
	False  *string
}
type Schema struct {
	SKeys    []string
	SColumns []string
	Keys     []*FieldDB
	Columns  []*FieldDB
	Fields   map[string]*FieldDB
}

func BuildFieldsBySchema(schema *Schema) string {
	columns := make([]string, 0)
	for _, s := range schema.SColumns {
		columns = append(columns, s)
	}
	return strings.Join(columns, ",")
}
func BuildQueryBySchema(table string, schema *Schema) string {
	columns := make([]string, 0)
	for _, s := range schema.SColumns {
		columns = append(columns, s)
	}
	return "select " + strings.Join(columns, ",") + " from " + table + " "
}
func CreateSchema(modelType reflect.Type) *Schema {
	m := modelType
	if m.Kind() == reflect.Ptr {
		m = m.Elem()
	}
	numField := m.NumField()
	scolumns := make([]string, 0)
	skeys := make([]string, 0)
	columns := make([]*FieldDB, 0)
	keys := make([]*FieldDB, 0)
	schema := make(map[string]*FieldDB, 0)
	for idx := 0; idx < numField; idx++ {
		field := m.Field(idx)
		tag, _ := field.Tag.Lookup("gorm")
		if !strings.Contains(tag, IgnoreReadWrite) {
			update := !strings.Contains(tag, "update:false")
			insert := !strings.Contains(tag, "insert:false")
			if has := strings.Contains(tag, "column"); has {
				json := field.Name
				col := json
				str1 := strings.Split(tag, ";")
				num := len(str1)
				for i := 0; i < num; i++ {
					str2 := strings.Split(str1[i], ":")
					for j := 0; j < len(str2); j++ {
						if str2[j] == "column" {
							isKey := strings.Contains(tag, "primary_key")
							col = str2[j+1]
							scolumns = append(scolumns, col)
							jTag, jOk := field.Tag.Lookup("json")
							if jOk {
								tagJsons := strings.Split(jTag, ",")
								json = tagJsons[0]
							}
							f := &FieldDB{
								JSON:   json,
								Column: col,
								Index:  idx,
								Scale:  -1,
								Key:    isKey,
								Update: update,
								Insert: insert,
							}
							if isKey {
								skeys = append(skeys, col)
								keys = append(keys, f)
							}
							tScale, sOk := field.Tag.Lookup("scale")
							if sOk {
								scale, err := strconv.Atoi(tScale)
								if err == nil {
									f.Scale = int8(scale)
								}
							}
							tTag, tOk := field.Tag.Lookup("true")
							if tOk {
								f.True = &tTag
								fTag, fOk := field.Tag.Lookup("false")
								if fOk {
									f.False = &fTag
								}
							}
							columns = append(columns, f)
							schema[col] = f
						}
					}
				}
			}
		}
	}
	s := &Schema{SColumns: scolumns, SKeys: skeys, Columns: columns, Keys: keys, Fields: schema}
	return s
}
func MakeSchema(modelType reflect.Type) ([]*FieldDB, []*FieldDB) {
	m := CreateSchema(modelType)
	return m.Columns, m.Keys
}
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
func QuoteColumnName(str string) string {
	//if strings.Contains(str, ".") {
	//	var newStrs []string
	//	for _, str := range strings.Split(str, ".") {
	//		newStrs = append(newStrs, str)
	//	}
	//	return strings.Join(newStrs, ".")
	//}
	return str
}
func GetDBValue(v interface{}, boolSupport bool, scale int8) (string, bool) {
	switch v.(type) {
	case string:
		s0 := v.(string)
		if len(s0) == 0 {
			return "''", true
		}
		return "", false
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
	case bool:
		if !boolSupport {
			return "", false
		}
		b0 := v.(bool)
		if b0 {
			return "true", true
		} else {
			return "false", true
		}
	case float64:
		if scale >= 0 {
			mt := "%." + strconv.Itoa(int(scale)) + "f"
			return fmt.Sprintf(mt, v), true
		}
		return "", false
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
		return "", false
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
func setValue(model interface{}, index int, value interface{}) (interface{}, error) {
	valueObject := reflect.Indirect(reflect.ValueOf(model))
	switch reflect.ValueOf(model).Kind() {
	case reflect.Ptr:
		{
			valueObject.Field(index).Set(reflect.ValueOf(value))
			return model, nil
		}
	default:
		if modelWithTypeValue, ok := model.(reflect.Value); ok {
			_, err := setValueWithTypeValue(modelWithTypeValue, index, value)
			return modelWithTypeValue.Interface(), err
		}
	}
	return model, nil
}
func setValueWithTypeValue(model reflect.Value, index int, value interface{}) (reflect.Value, error) {
	trueValue := reflect.Indirect(model)
	switch trueValue.Kind() {
	case reflect.Struct:
		{
			val := reflect.Indirect(reflect.ValueOf(value))
			if trueValue.Field(index).Kind() == val.Kind() {
				trueValue.Field(index).Set(reflect.ValueOf(value))
				return trueValue, nil
			} else {
				return trueValue, fmt.Errorf("value's kind must same as field's kind")
			}
		}
	default:
		return trueValue, nil
	}
}
func BuildMapDataAndKeys(model interface{}, update bool) (map[string]interface{}, map[string]interface{}, []string, []string) {
	var mapData = make(map[string]interface{})
	var mapKey = make(map[string]interface{})
	keys := make([]string, 0)
	columns := make([]string, 0)
	mv := reflect.Indirect(reflect.ValueOf(model))
	modelType := mv.Type()
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		if colName, isKey, exist := CheckByIndex(modelType, i, update); exist {
			f := mv.Field(i)
			fieldValue := f.Interface()
			isNil := false
			if f.Kind() == reflect.Ptr {
				if reflect.ValueOf(fieldValue).IsNil() {
					isNil = true
				} else {
					fieldValue = reflect.Indirect(reflect.ValueOf(fieldValue)).Interface()
				}
			}
			if isKey {
				columns = append(columns, colName)
				if !isNil {
					mapKey[colName] = fieldValue
				}
			} else {
				keys = append(keys, colName)
				if !isNil {
					if boolValue, ok := fieldValue.(bool); ok {
						valueS, okS := modelType.Field(i).Tag.Lookup(strconv.FormatBool(boolValue))
						if okS {
							mapData[colName] = valueS
						} else {
							mapData[colName] = boolValue
						}
					} else {
						mapData[colName] = fieldValue
					}
				}
			}
		}
	}
	return mapData, mapKey, keys, columns
}
func CheckByIndex(modelType reflect.Type, index int, update bool) (col string, isKey bool, colExist bool) {
	field := modelType.Field(index)
	tag, _ := field.Tag.Lookup("gorm")
	if strings.Contains(tag, IgnoreReadWrite) {
		return "", false, false
	}
	if update {
		if strings.Contains(tag, "update:false") {
			return "", false, false
		}
	}
	if has := strings.Contains(tag, "column"); has {
		str1 := strings.Split(tag, ";")
		num := len(str1)
		for i := 0; i < num; i++ {
			str2 := strings.Split(str1[i], ":")
			for j := 0; j < len(str2); j++ {
				if str2[j] == "column" {
					isKey := strings.Contains(tag, "primary_key")
					return str2[j+1], isKey, true
				}
			}
		}
	}
	return "", false, false
}
func BuildParamWithNull(colName string) string {
	return QuoteColumnName(colName) + "=null"
}
func BuildSqlParametersAndValues(columns []string, values []interface{}, n *int, start int, joinStr string, buildParam func(int) string) (string, []interface{}, error) {
	arr := make([]string, *n)
	j := start
	var valueParams []interface{}
	for i, _ := range arr {
		columnName := columns[i]
		if values[j] == nil {
			arr[i] = BuildParamWithNull(columnName)
			copy(values[i:], values[i+1:])
			values[len(values)-1] = ""
			values = values[:len(values)-1]
			*n--
		} else {
			arr[i] = fmt.Sprintf("%s = %s", columnName, BuildParametersFrom(j, 1, buildParam))
			valueParams = append(valueParams, values[j])
		}
		j++
	}
	return strings.Join(arr, joinStr), valueParams, nil
}
func BuildParametersFrom(i int, numCol int, buildParam func(int) string) string {
	var arrValue []string
	for j := 0; j < numCol; j++ {
		arrValue = append(arrValue, buildParam(i+j+1))
	}
	return strings.Join(arrValue, ",")
}
func statement() BatchStatement {
	attributes := make(map[string]interface{})
	attributeKeys := make(map[string]interface{})
	return BatchStatement{Keys: []string{}, Columns: []string{}, Attributes: attributes, AttributeKeys: attributeKeys}
}

const (
	t1 = "2006-01-02T15:04:05Z"
	t2 = "2006-01-02T15:04:05-0700"
	t3 = "2006-01-02T15:04:05.0000000-0700"

	l1 = len(t1)
	l2 = len(t2)
	l3 = len(t3)
)

func ToDates(args []interface{}) []int {
	if args == nil || len(args) == 0 {
		ag2 := make([]int, 0)
		return ag2
	}
	var dates []int
	for i, arg := range args {
		if _, ok := arg.(time.Time); ok {
			dates = append(dates, i)
		}
		if _, ok := arg.(*time.Time); ok {
			dates = append(dates, i)
		}
	}
	return dates
}

func ParseDates(args []interface{}, dates []int) []interface{} {
	if args == nil || len(args) == 0 {
		return nil
	}
	if dates == nil || len(dates) == 0 {
		return args
	}
	res := append([]interface{}{}, args...)
	for _, d := range dates {
		if d >= len(args) {
			break
		}
		a := args[d]
		if s, ok := a.(string); ok {
			switch len(s) {
			case l1:
				t, err := time.Parse(t1, s)
				if err == nil {
					res[d] = t
				}
			case l2:
				t, err := time.Parse(t2, s)
				if err == nil {
					res[d] = t
				}
			case l3:
				t, err := time.Parse(t3, s)
				if err == nil {
					res[d] = t
				}
			}
		}
	}
	return res
}
func RoundFloat(num float64, slice int) float64 {
	c := math.Pow10(slice)
	result := math.Ceil(num*c) / c
	return result
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
