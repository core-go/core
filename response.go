package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const InternalServerError = "Internal Server Error"

func ErrorWithMessage(w http.ResponseWriter, r *http.Request, code int, err string, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	RespondAndLog(w, r, code, err, writeLog, resource, action, true, err)
}
func Error(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), err error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	Respond(w, r, code, result)
}
func ErrorAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	RespondAndLog(w, r, code, result, writeLog, resource, action, false, err.Error())
}
func Succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	RespondAndLog(w, r, code, result, writeLog, resource, action, true, "")
}
func Respond(w http.ResponseWriter, r *http.Request, code int, result interface{}) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func RespondAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
}

func WriteLogWithGoRoutine(ctx context.Context, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	if writeLog != nil {
		token := ctx.Value("authorization")
		go func() {
			timeOut := 10 * time.Second
			ctxSaveLog, cancel := context.WithTimeout(context.Background(), timeOut)
			defer cancel()

			if authorizationToken, ok := token.(map[string]interface{}); ok {
				ctxSaveLog = context.WithValue(ctxSaveLog, "authorization", authorizationToken)
			}
			err := writeLog(ctxSaveLog, resource, action, success, desc)
			fmt.Printf("saveLogErr: %v\n", err)
		}()
	}
}

func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if result == nil {
		w.Write([]byte("null"))
		return nil
	}
	response, err := Marshal(result)
	if err != nil {
		// log.Println("cannot marshal of result: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	w.Write(response)
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	b, ok1 := v.([]byte)
	if ok1 {
		return b, nil
	}
	s, ok2 := v.(string)
	if ok2 {
		return []byte(s), nil
	}
	return json.Marshal(v)
}
