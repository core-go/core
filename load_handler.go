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
func ReturnWithLog(w http.ResponseWriter, r *http.Request, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) error {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		if logError != nil {
			logError(r.Context(), "GET "+r.URL.Path+" with error: "+err.Error())
		}
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "GET "+r.URL.Path+" with error: "+err.Error())
		}
		if logError == nil && writeLog == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return err
	} else {
		if IsNil(model) {
			if writeLog != nil {
				writeLog(r.Context(), resource, action, false, "GET "+r.URL.Path+" not found")
			}
			return JSON(w, http.StatusNotFound, nil)
		} else {
			if writeLog != nil {
				writeLog(r.Context(), resource, action, true, "GET "+r.URL.Path)
			}
			return JSON(w, http.StatusOK, model)
		}
	}
}
func Return(w http.ResponseWriter, r *http.Request, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	if err != nil {
		if logError != nil {
			logError(r.Context(), "GET "+r.URL.Path+" with error: "+err.Error())
			JSON(w, http.StatusInternalServerError, InternalServerError)
		} else {
			JSON(w, http.StatusInternalServerError, err.Error())
		}
		return err
	} else {
		if IsNil(model) {
			return JSON(w, http.StatusNotFound, nil)
		} else {
			return JSON(w, http.StatusOK, model)
		}
	}
}
