package core

import (
	"context"
	"reflect"
	"strings"
)

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

func FromContext(ctx context.Context, key string, options ...string) string {
	var authorization string
	if len(options) > 0 {
		authorization = options[0]
	}
	if len(authorization) > 0 {
		token := ctx.Value(authorization)
		if token != nil {
			if authorizationToken, exist := token.(map[string]interface{}); exist {
				return FromMap(key, authorizationToken)
			}
		}
		return ""
	} else {
		u := ctx.Value(key)
		if u != nil {
			v, ok := u.(string)
			if ok {
				return v
			}
		}
		return ""
	}
}
func FromMap(key string, data map[string]interface{}) string {
	if data == nil {
		return ""
	}
	u := data[key]
	if u != nil {
		v, ok := u.(string)
		if ok {
			return v
		}
	}
	return ""
}
