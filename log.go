package service

import (
	"context"
	"encoding/json"
	"net/http"
)

type LogWriter interface {
	Write(ctx context.Context, resource string, action string, success bool, desc string) error
}

func RespondString(w http.ResponseWriter, r *http.Request, code int, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(result))
}
func Respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, logWriter LogWriter, resource string, action string, success bool, desc string) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	if logWriter != nil {
		logWriter.Write(r.Context(), resource, action, success, desc)
	}
}

func Succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, logWriter LogWriter, resource string, action string) {
	Respond(w, r, code, result, logWriter, resource, action, true, "")
}
