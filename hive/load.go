package hive

import (
	"reflect"
	"strings"
)

func BuildFields(modelType reflect.Type) string {
	columns := GetFields(modelType)
	return strings.Join(columns, ",")
}
func GetFields(modelType reflect.Type) []string {
	m := modelType
	if m.Kind() == reflect.Ptr {
		m = m.Elem()
	}
	numField := m.NumField()
	columns := make([]string, 0)
	for idx := 0; idx < numField; idx++ {
		field := m.Field(idx)
		tag, _ := field.Tag.Lookup("gorm")
		if !strings.Contains(tag, IgnoreReadWrite) {
			if has := strings.Contains(tag, "column"); has {
				json := field.Name
				col := json
				str1 := strings.Split(tag, ";")
				num := len(str1)
				for i := 0; i < num; i++ {
					str2 := strings.Split(str1[i], ":")
					for j := 0; j < len(str2); j++ {
						if str2[j] == "column" {
							col = str2[j+1]
							columns = append(columns, col)
						}
					}
				}
			}
		}
	}
	return columns
}
func BuildQuery(table string, modelType reflect.Type) string {
	columns := GetFields(modelType)
	return "select " + strings.Join(columns, ",") + " from " + table + " "
}
