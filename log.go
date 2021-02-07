package service

import (
	"context"
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
}
func Error(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	Respond(w, r, code, result, writeLog, resource, action, false, err.Error())
}
func Succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	Respond(w, r, code, result, writeLog, resource, action, true, "")
}
