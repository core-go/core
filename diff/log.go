package diff

import (
	"context"
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
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
func handleError(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	respond(w, r, code, result, writeLog, resource, action, false, err.Error())
}
func succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	respond(w, r, code, result, writeLog, resource, action, true, "")
}
