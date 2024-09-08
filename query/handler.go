package query

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

const internalServerError = "Internal Server Error"

type Load func(ctx context.Context, key string, max int64) ([]string, error)

type QueryHandler struct {
	Load     func(ctx context.Context, key string, max int64) ([]string, error)
	LogError func(context.Context, string, ...map[string]interface{})
	Keyword  string
	Max      string
}

func NewQueryHandler(load func(ctx context.Context, key string, max int64) ([]string, error), logError func(context.Context, string, ...map[string]interface{}), opts ...string) *QueryHandler {
	keyword := "q"
	if len(opts) > 0 && len(opts[0]) > 0 {
		keyword = opts[0]
	}
	max := "max"
	if len(opts) > 1 && len(opts[1]) > 0 {
		max = opts[1]
	}
	return &QueryHandler{load, logError, keyword, max}
}
func (h *QueryHandler) Query(w http.ResponseWriter, r *http.Request) {
	ps := r.URL.Query()
	keyword := ps.Get(h.Keyword)
	if len(keyword) == 0 {
		vs := make([]string, 0)
		respondModel(w, r, vs, nil, h.LogError, nil)
	} else {
		max := ps.Get(h.Max)
		i, err := strconv.ParseInt(max, 10, 64)
		if err != nil {
			i = 20
		}
		if i < 0 {
			i = 20
		}
		vs, err := h.Load(r.Context(), keyword, i)
		respondModel(w, r, vs, err, h.LogError, nil)
	}
}

func respondModel(w http.ResponseWriter, r *http.Request, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		respondAndLog(w, r, http.StatusInternalServerError, internalServerError, err, logError, writeLog, resource, action)
	} else {
		if model == nil {
			returnAndLog(w, r, http.StatusNotFound, model, writeLog, false, resource, action, "Not found")
		} else {
			succeed(w, r, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
func respondAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) error {
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
			return returnAndLog(w, r, http.StatusInternalServerError, internalServerError, writeLog, false, resource, action, err.Error())
		} else {
			return returnAndLog(w, r, http.StatusInternalServerError, err.Error(), writeLog, false, resource, action, err.Error())
		}
	} else {
		return returnAndLog(w, r, code, result, writeLog, true, resource, action, "")
	}
}
func returnAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
	return err
}
func succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, options ...string) error {
	return respondAndLog(w, r, code, result, nil, nil, writeLog, options...)
}
