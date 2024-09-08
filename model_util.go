package core

import (
	"reflect"
	"strings"
	"time"
)

func Now() *time.Time {
	n := time.Now()
	return &n
}

func GetValue(model interface{}, index int) (interface{}, string, error) {
	valueObject := reflect.Indirect(reflect.ValueOf(model))
	return reflect.Indirect(valueObject.Field(index)).Interface(), valueObject.Type().Field(index).Name, nil
}

func BuildMapField(modelType reflect.Type) ([]string, map[string]int, map[string]int) {
	model := reflect.New(modelType).Interface()
	val := reflect.Indirect(reflect.ValueOf(model))
	var idFields []string
	m1 := make(map[string]int)
	m2 := make(map[string]int)
	l := val.Type().NumField()
	vt := val.Type()
	for i := 0; i < l; i++ {
		field := vt.Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 {
			jsonName := strings.Split(tag1, ",")[0]
			m1[jsonName] = i
		} else {
			m1[field.Name] = i
		}
		m2[field.Name] = i
		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
				jsonTag := field.Tag.Get("json")
				tags1 := strings.Split(jsonTag, ",")
				if len(tags1) > 0 && tags1[0] != "-" {
					idFields = append(idFields, tags1[0])
				}
			}
		}
	}
	return idFields, m1, m2
}
