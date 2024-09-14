package core

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

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
func MatchId(r *http.Request, body interface{}, keysJson []string, mapIndex map[string]int) error {
	var value reflect.Value
	value = reflect.ValueOf(body)
	modelType := value.Type()
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
		modelType = modelType.Elem()
	}
	_, mapParams, err2 := GetParamIds(r, keysJson, 0)
	if err2 != nil {
		return errors.New("Invalid Data " + err2.Error())
	}
	isNil := false
	for _, primaryField := range keysJson {
		indexField, okIndex := mapIndex[primaryField]
		if okIndex {
			if paramId, ok := mapParams[primaryField]; ok {
				indexes := []int{indexField}
				field := value.FieldByIndex(indexes)
				idType := field.Kind().String()
				isPointer := false
				if field.Kind() == reflect.Ptr {
					if !field.IsNil() {
						field = field.Elem()
						idType = field.Kind().String()
					} else {
						isNil = true
						f := modelType.Field(indexField)
						idType = f.Type.String()
					}
					isPointer = true
				}
				if isNil {
					if strings.Index(idType, "string") >= 0 {
						if isPointer {
							field.Set(reflect.ValueOf(&paramId))
						} else {
							field.Set(reflect.ValueOf(paramId))
						}
					} else {
						switch idType {
						case "int64", "*int64":
							i, err := strconv.ParseInt(paramId, 10, 64)
							if err != nil {
								return errors.New("invalid key: " + primaryField)
							}
							if isPointer {
								field.Set(reflect.ValueOf(&i))
							} else {
								field.Set(reflect.ValueOf(i))
							}
						case "int":
							i, err := strconv.Atoi(paramId)
							if err != nil {
								return errors.New("invalid key: " + primaryField)
							}
							if isPointer {
								field.Set(reflect.ValueOf(&i))
							} else {
								field.Set(reflect.ValueOf(i))
							}
						case "int32", "*int32":
							i, err := strconv.Atoi(paramId)
							if err != nil {
								return errors.New("invalid key: " + primaryField)
							}
							i32 := int32(i)
							if isPointer {
								field.Set(reflect.ValueOf(&i32))
							} else {
								field.Set(reflect.ValueOf(i32))
							}
						default:
						}
					}
				} else {
					v := field.Interface()
					sv, ok := v.(string)
					if ok {
						if len(sv) == 0 {
							field.Set(reflect.ValueOf(paramId))
						} else if sv != paramId {
							return errors.New("conflict key in param and body: " + primaryField)
						}
					} else {
						idValue, err := strconv.ParseInt(paramId, 10, 64)
						if err != nil {
							return errors.New("Parameter '" + primaryField + "' must be an integer : " + paramId)
						}
						i := field.Int()
						if i != 0 {
							if i != idValue {
								return errors.New("conflict key in param and body: " + primaryField)
							}
						} else {
							switch idType {
							case "int64":
								field.Set(reflect.ValueOf(idValue))
							case "int":
								i2 := int(idValue)
								field.Set(reflect.ValueOf(i2))
							case "int32":
								i2 := int32(idValue)
								field.Set(reflect.ValueOf(i2))
							default:
							}
						}
					}
				}
			} else {
				return errors.New("Not found param key: " + primaryField)
			}
		} else {
			return errors.New("Not found param key: " + primaryField)
		}
	}
	return nil
}
func MakeId(r *http.Request, modelType reflect.Type, idNames []string, indexes map[string]int, options ...int) (map[string]interface{}, error) {
	modelValue := reflect.New(modelType)
	mapKey := make(map[string]interface{})
	_, mapParams, err2 := GetParamIds(r, idNames, options...)
	if err2 != nil {
		return nil, err2
	}
	for _, idName := range idNames {
		if idValue, ok := mapParams[idName]; ok {
			if len(strings.Trim(idValue, " ")) == 0 {
				return nil, fmt.Errorf("%v is required", idName)
			}
			index, _ := indexes[idName]
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
}
func CreateId(r *http.Request, modelType reflect.Type, idNames []string, indexes map[string]int, options ...int) (interface{}, error) {
	if len(idNames) > 1 {
		return MakeId(r, modelType, idNames, indexes, options...)
	} else if len(idNames) == 1 {
		modelValue := reflect.New(modelType)
		idValue, _, err1 := GetParamIds(r, idNames, options...)
		if err1 != nil {
			return nil, err1
		}
		if idStr, ok := idValue.(string); ok {
			if len(strings.Trim(idStr, " ")) == 0 {
				return nil, fmt.Errorf("%v is required", idNames[0])
			}
			index, _ := indexes[idNames[0]]
			ifField := reflect.Indirect(modelValue).FieldByIndex([]int{index})
			idType := ifField.Type().String()
			switch idType {
			case "int64", "*int64":
				if id, err := strconv.ParseInt(idStr, 10, 64); err != nil {
					return nil, fmt.Errorf("%v is invalid", idNames[0])
				} else {
					return id, nil
				}
			case "int", "int32", "*int32":
				if id, err := strconv.ParseInt(idStr, 10, 32); err != nil {
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
func GetParamIds(r *http.Request, idNames []string, options ...int) (interface{}, map[string]string, error) {
	offset := 0
	if len(options) > 0 && options[0] > 0 {
		offset = options[0]
	}
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
