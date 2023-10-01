package hive

import (
	"errors"
	"reflect"
	"strings"
)

func GetIndexesByTagJson(modelType reflect.Type) (map[string]int, error) {
	ma := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagJson := field.Tag.Get("json")
		if len(tagJson) > 0 {
			ma[tagJson] = i
		}
	}
	return ma, nil
}
func GetColumnsSelect(modelType reflect.Type) []string {
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
						columnName := strings.ToLower(str2[j+1])
						columnNameKeys = append(columnNameKeys, columnName)
					}
				}
			}
		}
	}
	return columnNameKeys
}
