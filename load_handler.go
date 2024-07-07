package core

import (
	"context"
	"net/http"
	"reflect"
)

func GetId(w http.ResponseWriter, r *http.Request, modelType reflect.Type, jsonId []string, indexes map[string]int, options ...int) map[string]interface{} {
	id, err := MakeId(r, modelType, jsonId, indexes, options...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return id
}
func Return(w http.ResponseWriter, r *http.Request, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
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
func RespondIfFound(w http.ResponseWriter, r *http.Request, model interface{}, found bool, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
	if err == nil && !found {
		JSON(w, http.StatusNotFound, nil)
	} else {
		Return(w, r, model, err, logError, writeLog, options...)
	}
}
