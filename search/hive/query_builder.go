package query

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	s "github.com/core-go/core/search"
)

const (
	d0 = "2006-01-02"
	t0 = "2006-01-02 15:04:05"
	t1 = "2006-01-02T15:04:05Z"
	t2 = "2006-01-02T15:04:05-0700"
	t3 = "2006-01-02T15:04:05.0000000-0700"

	l1 = len(t1)
	l2 = len(t2)
	l3 = len(t3)

	desc = "desc"
	asc  = "asc"
)

type Builder[T any, F any] struct {
	TableName string
	ModelType reflect.Type
}

func UseQuery[T any, F any](tableName string) func(F) string {
	b := NewBuilder[T, F](tableName)
	return b.BuildQuery
}
func NewBuilder[T any, F any](tableName string) *Builder[T, F] {
	var t T
	resultModelType := reflect.TypeOf(t)
	if resultModelType.Kind() == reflect.Ptr {
		resultModelType = resultModelType.Elem()
	}
	return &Builder[T, F]{TableName: tableName, ModelType: resultModelType}
}
func (b *Builder[T, F]) BuildQuery(filter F) string {
	return Build(filter, b.TableName, b.ModelType)
}

const (
	like             = "like"
	greaterEqualThan = ">="
	greaterThan      = ">"
	lessEqualThan    = "<="
	lessThan         = "<"
	in               = "in"
)

func getStringFromTag(typeOfField reflect.StructField, tagName string, key string) *string {
	tag := typeOfField.Tag
	properties := strings.Split(tag.Get(tagName), ";")
	for _, property := range properties {
		if strings.HasPrefix(property, key) {
			column := property[len(key):]
			return &column
		}
	}
	return nil
}

func getJoinFromSqlBuilderTag(typeOfField reflect.StructField) *string {
	return getStringFromTag(typeOfField, "sql_builder", "join:")
}

func getColumnNameFromSqlBuilderTag(typeOfField reflect.StructField) *string {
	return getStringFromTag(typeOfField, "sql_builder", "column:")
}

func Build(filter interface{}, tableName string, modelType reflect.Type) string {
	s1 := ""
	rawConditions := make([]string, 0)
	// queryValues := make([]interface{}, 0)
	qQueryValues := make([]string, 0)
	qCols := make([]string, 0)
	rawJoin := make([]string, 0)
	sortString := ""
	fields := make([]string, 0)
	var excluding []string
	var keyword string
	value := reflect.Indirect(reflect.ValueOf(filter))
	filterType := value.Type()
	numField := value.NumField()
	var idCol string
	// marker := 0
	for i := 0; i < numField; i++ {
		columnName := getColumn(filterType, i)
		if columnName == "-" {
			continue
		}
		field := value.Field(i)
		kind := field.Kind()
		x := field.Interface()
		tf := value.Type().Field(i)
		fieldTypeName := tf.Type.String()
		typeOfField := value.Type().Field(i) // ???
		var psv string
		isContinue := false
		// param := buildParam(marker + 1)
		if kind == reflect.Ptr {
			if field.IsNil() {
				if fieldTypeName != "*string" {
					continue
				} else {
					isContinue = true
				}
			} else {
				field = field.Elem()
				kind = field.Kind()
				x = field.Interface()
			}
		}
		if !isContinue {
			s0, ok0 := x.(string)
			if ok0 {
				if len(s0) == 0 {
					isContinue = true
				}
				psv = s0
			}
		}
		if len(columnName) == 0 {
			_, _, columnName = getFieldByJson(modelType, tf.Name)
		}
		columnNameFromSqlBuilderTag := getColumnNameFromSqlBuilderTag(typeOfField)
		if columnNameFromSqlBuilderTag != nil {
			columnName = *columnNameFromSqlBuilderTag
		}

		joinFromSqlBuilderTag := getJoinFromSqlBuilderTag(typeOfField)
		if joinFromSqlBuilderTag != nil {
			rawJoin = append(rawJoin, *joinFromSqlBuilderTag)
		}
		if isContinue {
			if len(keyword) > 0 {
				qMatch, isQ := tf.Tag.Lookup("q")
				if isQ {
					if qMatch == "=" {
						qQueryValues = append(qQueryValues, keyword)
					} else if qMatch == "like" {
						qQueryValues = append(qQueryValues, buildQ(keyword))
					} else {
						qQueryValues = append(qQueryValues, prefix(keyword))
					}
					qCols = append(qCols, columnName)
				}
			}
			continue
		}
		if v, ok := x.(s.Filter); ok {
			if len(v.Fields) > 0 {
				for _, key := range v.Fields {
					i, _, columnName := getFieldByJson(modelType, key)
					if len(columnName) < 0 {
						fields = fields[len(fields):]
						break
					} else if i > -1 {
						fields = append(fields, columnName)
					}
				}
			}
			if len(fields) > 0 {
				s1 = `select ` + strings.Join(fields, ",") + ` from ` + tableName
			}
			if len(v.Sort) > 0 {
				sortString = buildSort(v.Sort, modelType)
			}
			if v.Excluding != nil && len(v.Excluding) > 0 {
				index, _, columnName := getFieldByBson(value.Type(), "_id")
				if !(index == -1 || columnName == "") {
					idCol = columnName
					excluding = v.Excluding
				}
			}
			if len(v.Q) > 0 {
				keyword = strings.TrimSpace(v.Q)
			}
			continue
		} else if len(psv) > 0 {
			key, ok := tf.Tag.Lookup("operator")
			if !ok {
				key, _ = tf.Tag.Lookup("q")
			}
			if key == "=" {
				param := WrapString(psv)
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %s", columnName, "=", param))
			} else {
				if key == "like" {
					rawConditions = append(rawConditions, fmt.Sprintf("%s %s %s", columnName, like, AllWrapString(psv)))
				} else {
					rawConditions = append(rawConditions, fmt.Sprintf("%s %s %s", columnName, like, PrefixWrapString(psv)))
				}
			}
		} else if kind == reflect.Slice {
			l := field.Len()
			if field.Len() > 0 {
				var arrValue []string
				for i := 0; i < l; i++ {
					model := field.Index(i).Addr()
					v, ok := GetDBValue(model, 2, "")
					if ok {
						arrValue = append(arrValue, v)
					}
				}
				format := fmt.Sprintf("(%s)", strings.Join(arrValue, ","))
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %s", columnName, in, format))
			}
		} else if dateTime, ok := x.(s.TimeRange); ok {
			if dateTime.Min != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s '%s'", columnName, greaterEqualThan, dateTime.Min.Format(t0)))
			}
			if dateTime.Max != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s '%s'", columnName, lessEqualThan, dateTime.Max.Format(t0)))
			} else if dateTime.Top != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %s", columnName, lessThan, dateTime.Top.Format(t0)))
			}
		} else if numberRange, ok := x.(s.NumberRange); ok {
			if numberRange.Min != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %.6f", columnName, greaterEqualThan, *numberRange.Min))
			} else if numberRange.Bottom != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %.6f", columnName, greaterThan, *numberRange.Bottom))
			}
			if numberRange.Max != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %.6f", columnName, lessEqualThan, *numberRange.Max))
			} else if numberRange.Top != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %.6f", columnName, lessThan, *numberRange.Top))
			}
		} else if numberRange, ok := x.(s.Int64Range); ok {
			if numberRange.Min != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, greaterEqualThan, numberRange.Min))
			} else if numberRange.Bottom != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, greaterThan, numberRange.Bottom))
			}
			if numberRange.Max != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, lessEqualThan, numberRange.Max))
			} else if numberRange.Top != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, lessThan, numberRange.Top))
			}
		} else if numberRange, ok := x.(s.IntRange); ok {
			if numberRange.Min != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, greaterEqualThan, numberRange.Min))
			} else if numberRange.Bottom != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, greaterThan, numberRange.Bottom))
			}
			if numberRange.Max != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, lessEqualThan, numberRange.Max))
			} else if numberRange.Top != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, lessThan, numberRange.Top))
			}
		} else if numberRange, ok := x.(s.Int32Range); ok {
			if numberRange.Min != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, greaterEqualThan, numberRange.Min))
			} else if numberRange.Bottom != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, greaterThan, numberRange.Bottom))
			}
			if numberRange.Max != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, lessEqualThan, numberRange.Max))
			} else if numberRange.Top != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %d", columnName, lessThan, numberRange.Top))
			}
		} else if dateRange, ok := x.(s.DateRange); ok {
			if dateRange.Min != nil {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s '%s'", columnName, greaterEqualThan, dateRange.Min.Format(d0)))
			}
			if dateRange.Max != nil {
				var eDate = dateRange.Max.Add(time.Hour * 24)
				dateRange.Max = &eDate
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %s", columnName, lessThan, eDate.Format(d0)))
			}
		} else {
			key, ok := tf.Tag.Lookup("operator")
			if !ok {
				key = "="
			}
			param, ok := GetDBValue(x, 2, "")
			if ok {
				rawConditions = append(rawConditions, fmt.Sprintf("%s %s %s", columnName, key, param))
			}
		}
	}

	if excluding != nil && len(excluding) > 0 && len(idCol) > 0 {
		var arrValue []string
		l := len(excluding)
		for i := 0; i < l; i++ {
			v, ok := GetDBValue(excluding[i], 2, "")
			if ok {
				arrValue = append(arrValue, v)
			}
		}
		format := fmt.Sprintf("(%s)", strings.Join(arrValue, ","))
		rawConditions = append(rawConditions, fmt.Sprintf("%s NOT IN %s", idCol, format))
	}
	if len(s1) == 0 {
		columns := getColumnsSelect(modelType)
		if len(columns) > 0 {
			s1 = `select  ` + strings.Join(columns, ",") + ` from ` + tableName
		} else {
			s1 = `select * from ` + tableName
		}
	}
	if len(rawJoin) > 0 {
		s1 = s1 + " " + strings.Join(rawJoin, " ")
	}
	if len(qCols) > 0 {
		qConditions := make([]string, 0)
		for i, s := range qCols {
			qConditions = append(qConditions, fmt.Sprintf("%s %s %s", s, like, qQueryValues[i]))
		}
		if len(qConditions) > 0 {
			rawConditions = append(rawConditions, " ("+strings.Join(qConditions, " or ")+") ")
		}
	}
	if len(rawConditions) > 0 {
		s2 := s1 + ` where ` + strings.Join(rawConditions, " AND ") + sortString
		return s2
	}
	s3 := s1 + sortString
	return s3
}
func getFieldByJson(modelType reflect.Type, jsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 && strings.Split(tag1, ",")[0] == jsonName {
			if tag2, ok2 := field.Tag.Lookup("gorm"); ok2 {
				if has := strings.Contains(tag2, "column"); has {
					str1 := strings.Split(tag2, ";")
					num := len(str1)
					for k := 0; k < num; k++ {
						str2 := strings.Split(str1[k], ":")
						for j := 0; j < len(str2); j++ {
							if str2[j] == "column" {
								return i, field.Name, str2[j+1]
							}
						}
					}
				}
			}
			return i, field.Name, ""
		}
	}
	return -1, jsonName, jsonName
}
func getColumn(filterType reflect.Type, i int) string {
	field := filterType.Field(i)
	if tag2, ok := field.Tag.Lookup("gorm"); ok {
		if tag2 == "-" {
			return tag2
		}
		if has := strings.Contains(tag2, "column"); has {
			str1 := strings.Split(tag2, ";")
			num := len(str1)
			for k := 0; k < num; k++ {
				str2 := strings.Split(str1[k], ":")
				for j := 0; j < len(str2); j++ {
					if str2[j] == "column" {
						return str2[j+1]
					}
				}
			}
		}
	}
	return ""
}
func getFieldByBson(modelType reflect.Type, bsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag1, ok1 := field.Tag.Lookup("bson")
		if ok1 && strings.Split(tag1, ",")[0] == bsonName {
			if tag2, ok2 := field.Tag.Lookup("gorm"); ok2 {
				if has := strings.Contains(tag2, "column"); has {
					str1 := strings.Split(tag2, ";")
					num := len(str1)
					for k := 0; k < num; k++ {
						str2 := strings.Split(str1[k], ":")
						for j := 0; j < len(str2); j++ {
							if str2[j] == "column" {
								return i, field.Name, str2[j+1]
							}
						}
					}
				}
			}
			return i, field.Name, ""
		}
	}
	return -1, bsonName, bsonName
}
func getColumnName(modelType reflect.Type, fieldName string) (col string, colExist bool) {
	field, ok := modelType.FieldByName(fieldName)
	if !ok {
		return fieldName, false
	}
	tag2, ok2 := field.Tag.Lookup("gorm")
	if !ok2 {
		return "", true
	}

	if has := strings.Contains(tag2, "column"); has {
		str1 := strings.Split(tag2, ";")
		num := len(str1)
		for i := 0; i < num; i++ {
			str2 := strings.Split(str1[i], ":")
			for j := 0; j < len(str2); j++ {
				if str2[j] == "column" {
					return str2[j+1], true
				}
			}
		}
	}
	//return gorm.ToColumnName(fieldName), false
	return fieldName, false
}
func getColumnsSelect(modelType reflect.Type) []string {
	numField := modelType.NumField()
	columnNameKeys := make([]string, 0)
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		if has := strings.Contains(ormTag, "column"); has {
			str1 := strings.Split(ormTag, ";")
			num := len(str1)
			for i := 0; i < num; i++ {
				str2 := strings.Split(str1[i], ":")
				for j := 0; j < len(str2); j++ {
					if str2[j] == "column" {
						columnName := str2[j+1]
						columnNameTag := getColumnNameFromSqlBuilderTag(field)
						if columnNameTag != nil {
							columnName = *columnNameTag
						}
						columnNameKeys = append(columnNameKeys, columnName)
					}
				}
			}
		}
	}
	return columnNameKeys
}
func buildSort(sortString string, modelType reflect.Type) string {
	var sort = make([]string, 0)
	sorts := strings.Split(sortString, ",")
	for i := 0; i < len(sorts); i++ {
		sortField := strings.TrimSpace(sorts[i])
		fieldName := sortField
		c := sortField[0:1]
		if c == "-" || c == "+" {
			fieldName = sortField[1:]
		}
		columnName := getColumnNameForSearch(modelType, fieldName)
		if len(columnName) > 0 {
			sortType := getSortType(c)
			sort = append(sort, columnName+" "+sortType)
		}
	}
	if len(sort) > 0 {
		return ` order by ` + strings.Join(sort, ",")
	} else {
		return ""
	}
}
func getColumnNameForSearch(modelType reflect.Type, sortField string) string {
	sortField = strings.TrimSpace(sortField)
	i, _, column := getFieldByJson(modelType, sortField)
	if i > -1 {
		return column
	}
	return ""
}
func getSortType(sortType string) string {
	if sortType == "-" {
		return desc
	} else {
		return asc
	}
}

func buildParametersFrom(i int, numCol int, buildParam func(i int) string) string {
	var arrValue []string
	for j := 0; j < numCol; j++ {
		arrValue = append(arrValue, buildParam(i+j+1))
	}
	return strings.Join(arrValue, ",")
}
func buildQ(s string) string {
	if !(strings.HasPrefix(s, "%") && strings.HasSuffix(s, "%")) {
		return "%" + s + "%"
	} else if strings.HasPrefix(s, "%") {
		return s + "%"
	} else if strings.HasSuffix(s, "%") {
		return "%" + s
	}
	return s
}
func prefix(s string) string {
	if strings.HasSuffix(s, "%") {
		return s
	} else {
		return s + "%"
	}
}
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
func PrefixWrapString(v string) string {
	if strings.Index(v, `'`) >= 0 {
		return join(`'`, strings.Replace(v, "'", "''", -1), `'%`)
	}
	return join(`'`, v, `%'`)
}
func AllWrapString(v string) string {
	if strings.Index(v, `'`) >= 0 {
		return join(`'%`, strings.Replace(v, "'", "''", -1), `'%`)
	}
	return join(`'%`, v, `%'`)
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
