package core

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type Builder interface {
	Create(ctx context.Context, model interface{}) (interface{}, error)
	Update(ctx context.Context, model interface{}) (interface{}, error)
	Patch(ctx context.Context, model interface{}) (interface{}, error)
}
func BuildMapAndStruct(r *http.Request, interfaceBody interface{}, options...http.ResponseWriter) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()
	body := make(map[string]interface{})
	er1 := json.NewDecoder(strings.NewReader(s)).Decode(&body)
	if er1 != nil {
		return nil, er1
	}
	er2 := json.NewDecoder(strings.NewReader(s)).Decode(interfaceBody)
	if er2 != nil {
		if len(options) > 0 && options[0] != nil {
			http.Error(options[0], er2.Error(), http.StatusBadRequest)
		}
		return nil, er2
	}
	return body, nil
}
func BodyToJsonMap(r *http.Request, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int, options...func(context.Context, interface{}) (interface{}, error)) (map[string]interface{}, error) {
	var controlModel interface{}
	var buildToPatch func(context.Context, interface{}) (interface{}, error)
	if len(options) > 0 && options[0] != nil {
		buildToPatch = options[0]
	}
	if buildToPatch != nil {
		var er0 error
		controlModel, er0 = buildToPatch(r.Context(), structBody)
		if er0 != nil {
			return nil, er0
		}
		inRec, er1 := json.Marshal(controlModel)
		if er1 != nil {
			return nil, er1
		}
		var model map[string]interface{}
		json.Unmarshal(inRec, &model)
		for k, v := range model {
			stringKind := reflect.TypeOf(v).String()
			if (v != nil && stringKind == "float64" && v.(float64) != 0) || (v != nil && stringKind != "float64" && v != "") {
				body[k] = v
			}
		}
	}
	valueOfReq := reflect.ValueOf(structBody)
	if valueOfReq.Kind() == reflect.Ptr {
		valueOfReq = reflect.Indirect(valueOfReq)
	}
	for _, jsonName := range jsonIds {
		if i, ok := mapIndex[jsonName]; ok && i >= 0 {
			v, _, er4 := GetValue(structBody, i)
			if er4 == nil {
				body[jsonName] = v
			}
		}
	}
	result := make(map[string]interface{})
	for keyJsonName, _ := range body {
		v2 := body[keyJsonName]
		if v2 == nil {
			result[keyJsonName] = v2
		} else if i, ok := mapIndex[keyJsonName]; ok && i >= 0 {
			v, _, er4 := GetValue(structBody, i)
			if er4 == nil {
				result[keyJsonName] = v
			}
		}
	}
	return result, nil
}
