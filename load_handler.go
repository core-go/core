package core

import (
	"context"
	"net/http"
	"reflect"
)

type WriteLog func(context.Context, string, string, bool, string) error
type LoadHandler struct {
	LoadData   func(ctx context.Context, id interface{}) (interface{}, error)
	Keys       []string
	ModelType  reflect.Type
	KeyIndexes map[string]int
	Error      func(context.Context, string, ...map[string]interface{})
	WriteLog   WriteLog
	Resource   string
	Activity   string
}

func NewQueryHandler(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewQueryHandlerWithLog(load, modelType, logError, writeLog)
}
func NewQueryHandlerWithKeys(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewQueryHandlerWithKeysAndLog(load, keys, modelType, logError, writeLog)
}
func NewQueryHandlerWithLog(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
	return NewQueryHandlerWithKeysAndLog(load, nil, modelType, logError, writeLog, options...)
}
func NewQueryHandlerWithKeysAndLog(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
	return NewLoadHandlerWithKeysAndLog(load, keys, modelType, logError, writeLog, options...)
}
func NewLoadHandler(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewLoadHandlerWithLog(load, modelType, logError, writeLog)
}
func NewLoadHandlerWithKeys(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewLoadHandlerWithKeysAndLog(load, keys, modelType, logError, writeLog)
}
func NewLoadHandlerWithLog(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
	return NewLoadHandlerWithKeysAndLog(load, nil, modelType, logError, writeLog, options...)
}
func NewLoadHandlerWithKeysAndLog(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
	if keys == nil || len(keys) == 0 {
		keys = GetJsonPrimaryKeys(modelType)
	}
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		action = options[0]
	} else {
		action = "load"
	}
	if len(options) > 1 && len(options[1]) > 0 {
		resource = options[1]
	} else {
		resource = BuildResourceName(modelType.Name())
	}
	indexes := GetKeyIndexes(modelType)
	return &LoadHandler{WriteLog: writeLog, LoadData: load, Keys: keys, ModelType: modelType, KeyIndexes: indexes, Error: logError, Resource: resource, Activity: action}
}
func (h *LoadHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, er1 := BuildId(r, h.ModelType, h.Keys, h.KeyIndexes)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	} else {
		model, er2 := h.LoadData(r.Context(), id)
		Return(w, r, model, er2, h.Error, h.WriteLog, h.Resource, h.Activity)
	}
}
func GetId(w http.ResponseWriter, r *http.Request, modelType reflect.Type, jsonId []string, indexes map[string]int, options... int) map[string]interface{} {
	id, err := MakeId(r, modelType, jsonId, indexes, options...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return id
}
func GetStatus(status int64) int {
	if status <= 0 {
		return http.StatusNotFound
	}
	return http.StatusOK
}
func IsFound(res interface{}) int {
	if IsNil(res) {
		return http.StatusNotFound
	}
	return http.StatusOK
}
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
func Return(w http.ResponseWriter, r *http.Request, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options... string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		RespondAndLog(w, r, http.StatusInternalServerError, InternalServerError, err, logError, writeLog, resource, action)
	} else {
		if IsNil(model) {
			ReturnAndLog(w, r, http.StatusNotFound, nil, writeLog, false, resource, action, "Not found")
		} else {
			Succeed(w, r, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
func RespondIfFound(w http.ResponseWriter, r *http.Request, model interface{}, found bool, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options... string) {
	if err == nil && !found {
		JSON(w, http.StatusNotFound, nil)
	} else {
		Return(w, r, model, err, logError, writeLog, options...)
	}
}
