package response

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var internalServerError = "Internal Server Error"

func HandleResult(w http.ResponseWriter, r *http.Request, id string, res int64, err error, resource string, action string, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error) {
	if err != nil {
		if logError != nil {
			logError(r.Context(), err.Error())
		}
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, err.Error())
		}
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}
	Return(r.Context(), w, id, res, resource, action, writeLog)
}
func Return(ctx context.Context, w http.ResponseWriter, id string, res int64, resource string, action string, opts...func(context.Context, string, string, bool, string) error) {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(opts) > 0 {
		writeLog = opts[0]
	}
	if res > 0 {
		if writeLog != nil {
			writeLog(ctx, resource, action, true, fmt.Sprintf("%s '%s'", action, id))
		}
		JSON(w, http.StatusOK, res)
	} else if res == 0 {
		if writeLog != nil {
			writeLog(ctx, resource, action, false, fmt.Sprintf("not found '%s'", id))
		}
		JSON(w, http.StatusNotFound, res)
	} else {
		if writeLog != nil {
			writeLog(ctx, resource, action, false, fmt.Sprintf("conflict '%s'", id))
		}
		JSON(w, http.StatusConflict, res)
	}
}
func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	return err
}
