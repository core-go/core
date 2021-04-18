package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func GetJsonPrimaryKeys(modelType reflect.Type) []string {
	numField := modelType.NumField()
	var idFields []string
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
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
	return idFields
}
func IsTheSameType(modelType reflect.Type, obj interface{}) bool {
	typeOf := reflect.TypeOf(obj)
	if typeOf == modelType || (typeOf.Kind() == reflect.Ptr && reflect.Indirect(reflect.ValueOf(obj)).Type() == modelType) {
		return true
	}
	return false
}
// Return the struct, not pointer
func NewModel(modelType reflect.Type, body io.ReadCloser) (out interface{}, err error) {
	if body != nil {
		req := reflect.New(modelType).Interface()
		err := json.NewDecoder(body).Decode(&req)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
	return nil, nil
}
func BuildResourceName(s string) string {
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if s2[i] != s[i] {
			s3 += "-" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}
	if string(s3[0]) == "-" || string(s3[0]) == "_" {
		return s3[1:]
	}
	return s3
}
func CheckId(r *http.Request, body interface{}, keysJson []string, mapIndex map[string]int) error {
	var value reflect.Value
	value = reflect.Indirect(reflect.ValueOf(body))//value must be struct
	_, mapParams, err2 := GetParamIds(r, keysJson, 0)
	if err2 != nil {
		return errors.New("Invalid Data "+err2.Error())
	}
	for _, primaryField := range keysJson {
		indexField, okIndex := mapIndex[primaryField]
		if okIndex {
			if paramId, ok := mapParams[primaryField]; ok {
				indexes := []int{indexField}
				field := value.FieldByIndex(indexes)
				idType := field.Kind().String()
				if field.IsZero() {
					if strings.Index(idType, "int") >= 0 {
						valueField, err := ParseIntWithType(paramId, idType)
						if err != nil {
							return errors.New("invalid key: "+primaryField)
						}
						field.Set(reflect.ValueOf(valueField))
					} else {
						field.Set(reflect.ValueOf(paramId))
					}
				} else {
					if strings.Index(idType, "int") >= 0 {
						idValue, err := strconv.ParseInt(paramId, 10, 64)
						if err != nil || field.Int() != idValue {
							return errors.New("conflict key in param and body: "+primaryField)
						}
					} else {
						if !reflect.DeepEqual(field.Interface(), paramId) {
							return errors.New("conflict key in param and body: "+primaryField)
						}
					}
				}
			} else {
				return errors.New("Not found param key: "+primaryField)
			}
		} else {
			return errors.New("Not found param key: "+primaryField)
		}
	}
	return nil
}
func BuildId(r *http.Request, modelType reflect.Type, idNames []string, indexs map[string]int, offset int) (interface{}, error) {
	modelValue := reflect.New(modelType)
	if len(idNames) > 1 {
		mapKey := make(map[string]interface{})
		_, mapParams, err2 := GetParamIds(r, idNames, offset)
		if err2 != nil {
			return nil, err2
		}
		for _, idName := range idNames {
			if idValue, ok := mapParams[idName]; ok {
				if len(strings.Trim(idValue, " ")) == 0 {
					return nil, fmt.Errorf("%v is required", idName)
				}
				index, _ := indexs[idName]
				ifField := reflect.Indirect(modelValue).FieldByIndex([]int{index})
				idType := ifField.Type().String()
				switch idType {
				case "int64", "*int64":
					if id, err := strconv.ParseInt(idValue, 10, 64); err != nil {
						return nil, fmt.Errorf("%v is invalid", idName)
					} else {
						mapKey[idName] = id
					}
				case "int", "int32", "*int32":
					if id, err := strconv.ParseInt(idValue, 10, 32); err != nil {
						return nil, fmt.Errorf("%v is invalid", idName)
					} else {
						mapKey[idName] = id
					}
				default:
					mapKey[idName] = idValue
				}
			} else {
				return nil, fmt.Errorf("%v is required", idName)
			}
		}
		return mapKey, nil
	} else if len(idNames) == 1 {
		idValue, _, err1 := GetParamIds(r, idNames, offset)
		if err1 != nil {
			return nil, err1
		}
		if idstr, ok := idValue.(string); ok {
			if len(strings.Trim(idstr, " ")) == 0 {
				return nil, fmt.Errorf("%v is required", idNames[0])
			}
			index, _ := indexs[idNames[0]]
			ifField := reflect.Indirect(modelValue).FieldByIndex([]int{index})
			idType := ifField.Type().String()
			switch idType {
			case "int64", "*int64":
				if id, err := strconv.ParseInt(idstr, 10, 64); err != nil {
					return nil, fmt.Errorf("%v is invalid", idNames[0])
				} else {
					return id, nil
				}
			case "int", "int32", "*int32":
				if id, err := strconv.ParseInt(idstr, 10, 32); err != nil {
					return nil, fmt.Errorf("%v is invalid", idNames[0])
				} else {
					return id, nil
				}
			default:
				return idValue, nil
			}
		} else {
			return nil, errors.New("error parser string get id by uri")
		}
	} else {
		return nil, errors.New("invalid model type: no id of this model type")
	}
}
func BuildIds(r *http.Request, modelType reflect.Type, idNames []string) (interface{}, error) {
	if len(idNames) > 1 {
		return NewModels(r.Body, modelType)
	} else if len(idNames) == 1 {
		modelTypeKey := GetFieldType(modelType, idNames[0])
		if modelTypeKey != nil {
			return NewModels(r.Body, modelTypeKey)
		}
	}
	return nil, errors.New("invalid model type: no id of this model type")
}
func GetFieldType(modelType reflect.Type, jsonName string) reflect.Type {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if tag, ok := field.Tag.Lookup("json"); ok {
			if strings.Split(tag, ",")[0] == jsonName {
				return field.Type
			}
		}
	}
	return nil
}
func NewModels(body interface{}, modelType reflect.Type) (out interface{}, err error) {
	req := reflect.New(reflect.SliceOf(modelType)).Interface()
	if body != nil {
		switch dec := body.(type) {
		case io.Reader:
			err := json.NewDecoder(dec).Decode(&req)
			if err != nil {
				return nil, err
			}
			return req, nil
		}
	}
	return nil, nil
}
func GetIndexes(modelType reflect.Type) map[string]int {
	numField := modelType.NumField()
	mapJsonNameIndex := make(map[string]int, 0)
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)

		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
				jsonTag := field.Tag.Get("json")
				tags1 := strings.Split(jsonTag, ",")
				if len(tags1) > 0 && tags1[0] != "-" {
					mapJsonNameIndex[tags1[0]] = i
				}
			}
		}
	}
	return mapJsonNameIndex
}
func GetParamIds(r *http.Request, idNames []string, offset int) (interface{}, map[string]string, error) {
	sizeName := len(idNames)
	params := strings.Split(r.RequestURI, "/")
	// remove some item last array
	params = params[:len(params)-offset]
	sizeParam := len(params)
	start := sizeParam - sizeName
	if sizeParam >= start {
		// get params
		params = params[start:sizeParam]
		mapParams := make(map[string]string)
		if sizeName == 1 {
			if len(params) != 1 {
				return nil, nil, errors.New("bad request")
			}
			// convert map param
			mapParams[idNames[0]] = params[0]
			return params[0], mapParams, nil
		}
		// convert map param
		for i, v := range params {
			mapParams[idNames[i]] = v
		}
		return params, mapParams, nil
	}
	return nil, nil, errors.New("bad request")
}
