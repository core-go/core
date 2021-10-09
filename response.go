package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const InternalServerError = "Internal Server Error"

func ErrorWithMessage(w http.ResponseWriter, r *http.Request, code int, err string, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	RespondAndLog(w, r, code, err, writeLog, resource, action, true, err)
}
func Error(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), err error) error {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	return JSON(w, code, result)
}
func ErrorAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) error {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	return RespondAndLog(w, r, code, result, writeLog, resource, action, false, err.Error())
}
func Succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return RespondAndLog(w, r, code, result, writeLog, resource, action, true, "")
}
func RespondAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
	return err
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

func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	return err
}
func GetParam(r *http.Request, options ...int) string {
	offset := 0
	if len(options) > 0 && options[0] > 0 {
		offset = options[0]
	}
	s := r.URL.Path
	params := strings.Split(s, "/")
	i := len(params) - 1 - offset
	if i >= 0 {
		return params[i]
	} else {
		return ""
	}
}
func GetInt(r *http.Request, options ...int) (int, bool) {
	s := GetParam(r, options...)
	if len(s) == 0 {
		return 0, false
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return i, true
}
func GetInt64(r *http.Request, options ...int) (int64, bool) {
	s := GetParam(r, options...)
	if len(s) == 0 {
		return 0, false
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}
func GetInt32(r *http.Request, options ...int) (int32, bool) {
	s := GetParam(r, options...)
	if len(s) == 0 {
		return 0, false
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return int32(i), true
}
