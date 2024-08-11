package search

import (
	"context"
	"encoding/json"
	"net/http"
)

func RespondError(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string, ...map[string]interface{}), resource string, action string, err error, writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	Respond(w, r, code, result, writeLog, resource, action, false, err.Error())
}
func Respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error, resource string, action string, success bool, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
}

func succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error, resource string, action string) {
	Respond(w, r, code, result, writeLog, resource, action, true, "")
}
