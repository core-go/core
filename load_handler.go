package service

import (
	"context"
	"net/http"
	"reflect"
)

type LoadHandler struct {
	LoadData  func(ctx context.Context, id interface{}) (interface{}, error)
	Keys      []string
	ModelType reflect.Type
	Indexes   map[string]int
	Error     func(context.Context, string)
	WriteLog  func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource  string
	Action    string
}

func NewLoadHandler(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewLoadHandlerWithLog(load, modelType, logError, writeLog)
}
func NewLoadHandlerWithKeys(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewLoadHandlerWithKeysAndLog(load, keys, modelType, logError, writeLog)
}
func NewLoadHandlerWithLog(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
	return NewLoadHandlerWithKeysAndLog(load, nil, modelType, logError, writeLog, options...)
}
func NewLoadHandlerWithKeysAndLog(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
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
	indexes := GetIndexes(modelType)
	return &LoadHandler{WriteLog: writeLog, LoadData: load, Keys: keys, ModelType: modelType, Indexes: indexes, Error: logError, Resource: resource, Action: action}
}
func (h *LoadHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, er1 := BuildId(r, h.ModelType, h.Keys, h.Indexes)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	} else {
		model, er2 := h.LoadData(r.Context(), id)
		RespondModel(w, r, model, er2, h.Error, h.WriteLog, h.Resource, h.Action)
	}
}
func GetId(w http.ResponseWriter, r *http.Request, modelType reflect.Type, jsonId []string, indexes map[string]int, options... int) (map[string]interface{}, error) {
	id, err := MakeId(r, modelType, jsonId, indexes, options...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return id, err
}
func RespondModel(w http.ResponseWriter, r *http.Request, model interface{}, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, err, logError, writeLog, resource, action)
	} else {
		if model == nil {
			RespondAndLog(w, r, http.StatusNotFound, model, writeLog, false, resource, action, "Not found")
		} else {
			Succeed(w, r, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
