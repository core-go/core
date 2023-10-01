package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Result(w http.ResponseWriter, r *http.Request, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), opts...interface{}) error {
	if err != nil {
		if len(opts) > 0 && opts[0] != nil {
			b, er2 := json.Marshal(opts[0])
			if er2 == nil {
				m := make(map[string]interface{})
				m["request"] = string(b)
				logError(r.Context(), err.Error(), m)
			} else {
				logError(r.Context(), err.Error())
			}
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return err
		} else {
			logError(r.Context(), err.Error(), nil)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return err
		}
	} else {
		return JSON(w, code, result)
	}
}
func ErrorWithMessage(w http.ResponseWriter, r *http.Request, code int, err string, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	ReturnAndLog(w, r, code, err, writeLog, true, resource, action, err)
}
func Error(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), err error) error {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	return JSON(w, code, result)
}
func RespondAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options... string) error {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		if logError != nil {
			logError(r.Context(), err.Error())
			return ReturnAndLog(w, r, http.StatusInternalServerError, InternalServerError, writeLog, false, resource, action, err.Error())
		} else {
			return ReturnAndLog(w, r, http.StatusInternalServerError, err.Error(), writeLog, false, resource, action, err.Error())
		}
	} else {
		return ReturnAndLog(w, r, code, result, writeLog, true, resource, action, "")
	}
}
func Succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, options... string) error {
	return RespondAndLog(w, r, code, result, nil, nil, writeLog, options...)
}
func ReturnAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
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
