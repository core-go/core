package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func SetField(v interface{}, name string, value string) error {
	// v must be a pointer to a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to struct")
	}

	// Dereference pointer
	rv = rv.Elem()

	// Lookup field by name
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		return fmt.Errorf("not a field name: %s", name)
	}

	// Field must be exported
	if !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", name)
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fv.SetString(value)
	} else {
		fv.SetInt(i)
	}

	return nil
}

func GetJsonName(modelType reflect.Type, fieldName string) (string, bool) {
	field, ok := modelType.FieldByName(fieldName)
	if ok {
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 {
			return strings.Split(tag1, ",")[0], ok1
		}
	}
	return "", false
}

func GetValue(model interface{}, index int) (interface{}, string, error) {
	valueObject := reflect.Indirect(reflect.ValueOf(model))
	return reflect.Indirect(valueObject.Field(index)).Interface(), valueObject.Type().Field(index).Name, nil
}

func GetField(value interface{}, jsonName string) (int, string) {
	val := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 {
			v := strings.Split(tag1, ",")[0]
			if v == jsonName {
				return i, field.Name
			}
		}
	}
	return -1, ""
}

func BuildMapField(modelType reflect.Type) (map[string]int, map[string]int) {
	model := reflect.New(modelType).Interface()
	val := reflect.Indirect(reflect.ValueOf(model))
	m1 := make(map[string]int)
	m2 := make(map[string]int)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 {
			jsonName := strings.Split(tag1, ",")[0]
			m2[jsonName] = i
		} else {
			m2[field.Name] = i
		}
		m1[field.Name] = i
	}
	return m1, m2
}

func ParseIntWithType(value string, idType string) (v interface{}, err error) {
	switch idType {
	case "int64", "*int64":
		return strconv.ParseInt(value, 10, 64)
	case "int", "int32", "*int32":
		return strconv.Atoi(value)
	default:
	}
	return value, nil
}

func FromContext(ctx context.Context, key string, options...string) string {
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