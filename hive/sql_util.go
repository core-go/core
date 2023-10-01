package hive

import (
	"context"
	"errors"
	hv "github.com/beltran/gohive"
	"reflect"
	"strings"
)

const (
	IgnoreReadWrite = "-"
)

func Query(ctx context.Context, cursor *hv.Cursor, fieldsIndex map[string]int, results interface{}, sql string) error {
	cursor.Exec(ctx, sql)
	if cursor.Err != nil {
		return cursor.Err
	}
	modelType := reflect.TypeOf(results).Elem().Elem()
	tb, er3 := Scan(cursor, modelType, fieldsIndex)
	if er3 != nil {
		return er3
	}
	for _, element := range tb {
		appendToArray(results, element)
	}
	return nil
}

func appendToArray(arr interface{}, item interface{}) interface{} {
	arrValue := reflect.ValueOf(arr)
	elemValue := reflect.Indirect(arrValue)

	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = reflect.Indirect(itemValue)
	}
	elemValue.Set(reflect.Append(elemValue, itemValue))
	return arr
}
func GetColumnIndexes(modelType reflect.Type) (map[string]int, error) {
	ma := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := FindTag(ormTag, "column")
		column = strings.ToLower(column)
		if ok {
			ma[column] = i
		}
	}
	return ma, nil
}

func FindTag(tag string, key string) (string, bool) {
	if has := strings.Contains(tag, key); has {
		str1 := strings.Split(tag, ";")
		num := len(str1)
		for i := 0; i < num; i++ {
			str2 := strings.Split(str1[i], ":")
			for j := 0; j < len(str2); j++ {
				if str2[j] == key {
					return str2[j+1], true
				}
			}
		}
	}
	return "", false
}

func GetColumns(cursors *hv.Cursor) ([]string, map[string]string, error) {
	var mcols = make(map[string]string, 0)
	var columnNames = make([]string, 0)
	m := cursors.Description()
	for _, v := range m {
		k := v[0]
		arr := strings.Split(k, ".")
		col := arr[len(arr)-1]
		columnNames = append(columnNames, col)
		mcols[col] = k
	}
	return columnNames, mcols, nil
}

func Scan(cursors *hv.Cursor, modelType reflect.Type, fieldsIndex map[string]int) (t []interface{}, err error) {
	if fieldsIndex == nil {
		fieldsIndex, err = GetColumnIndexes(modelType)
		if err != nil {
			return
		}
	}
	columns, mcols, er0 := GetColumns(cursors)
	if er0 != nil {
		return nil, er0
	}
	ctx := context.Background()
	for cursors.HasMore(ctx) {
		initModel := reflect.New(modelType).Interface()
		r, _ := StructScan(initModel, columns, fieldsIndex)
		fieldPointers := cursors.RowMap(ctx)
		if cursors.Err != nil {
			return t, cursors.Err
		}
		for _, c := range columns {
			if colm, ok := mcols[c]; ok {
				if v, ok := fieldPointers[colm]; ok {
					if v != nil {
						v = reflect.Indirect(reflect.ValueOf(v)).Interface()
						if fieldValue, ok := r[c]; ok && !IsZeroOfUnderlyingType(v) {
							err1 := ConvertAssign(fieldValue, v)
							if err1 != nil {
								return nil, err1
							}
						}
					}
				}
			}
		}
		t = append(t, initModel)
	}
	return
}
func StructScan(s interface{}, columns []string, fieldsIndex map[string]int) (r map[string]interface{}, swapValues map[int]interface{}) {
	return StructScanAndIgnore(s, columns, fieldsIndex, -1)
}
func StructScanAndIgnore(s interface{}, columns []string, fieldsIndex map[string]int, indexIgnore int) (r map[string]interface{}, swapValues map[int]interface{}) {
	if s != nil {
		modelType := reflect.TypeOf(s).Elem()
		swapValues = make(map[int]interface{}, 0)
		r = make(map[string]interface{}, 0)
		maps := reflect.Indirect(reflect.ValueOf(s))

		if columns == nil {
			for i := 0; i < maps.NumField(); i++ {
				tagBool := modelType.Field(i).Tag.Get("true")
				if tagBool == "" {
					key := modelType.Field(i).Tag.Get("json") //TODO get tag of gorm
					r[key] = maps.Field(i).Addr().Interface()
				} else {
					var str string
					swapValues[i] = reflect.New(reflect.TypeOf(str)).Elem().Addr().Interface()
				}
			}
			return
		}

		for i, columnsName := range columns {
			if i == indexIgnore {
				continue
			}
			var index int
			var ok bool
			var modelField reflect.StructField
			var valueField reflect.Value
			if fieldsIndex == nil {
				if modelField, ok = modelType.FieldByName(columnsName); !ok {
					var t interface{}
					r[columnsName] = &t
					continue
				}
				valueField = maps.FieldByName(columnsName)
			} else {
				if index, ok = fieldsIndex[columnsName]; !ok {
					var t interface{}
					r[columnsName] = &t
					continue
				}
				modelField = modelType.Field(index)
				valueField = maps.Field(index)
			}

			x := valueField.Addr().Interface()
			tagBool := modelField.Tag.Get("true")
			if tagBool == "" {
				r[columnsName] = x
			} else {
				var str string
				y := reflect.New(reflect.TypeOf(str))
				swapValues[index] = y.Elem().Addr().Interface()
				r[columnsName] = swapValues[index]
			}
		}
	}
	return
}
func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
