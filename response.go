package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	InternalServerError = "Internal Server Error"

	t1 = "2006-01-02T15:04:05Z"
	t2 = "2006-01-02T15:04:05-0700"
	t3 = "2006-01-02T15:04:05.0000000-0700"

	l1 = len(t1)
	l2 = len(t2)
	l3 = len(t3)
)

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

func GetParam(r *http.Request, options... int) string {
	offset := 0
	if len(options) > 0 && options[0] > 0 {
		offset = options[0]
	}
	s := r.URL.Path
	params := strings.Split(s, "/")
	i := len(params)-1-offset
	if i >= 0 {
		return params[i]
	} else {
		return ""
	}
}

func GetParams(r *http.Request, options ...int) []string {
	p := GetParam(r, options...)
	return strings.Split(p, ",")
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
func GetTime(r *http.Request, options ...int) *time.Time {
	s := GetParam(r, options...)
	return CreateTime(s)
}
func CreateTime(s string) *time.Time {
	l := len(s)
	p := ""
	switch l {
	case l1:
		p = t1
	case l2:
		p = t2
	case l3:
		p = t3
	default:
		p = ""
	}
	if len(p) == 0 {
		return nil
	}
	t, err := time.Parse(p, s)
	if err != nil {
		return nil
	}
	return &t
}
func QueryString(v url.Values, name string, options... string) string {
	s := v.Get(name)
	if len(s) > 0 {
		return s
	}
	if len(options) > 0 {
		return options[0]
	}
	return ""
}
func QueryStrings(v url.Values, name string, options...[]string) []string {
	s, ok := v[name]
	if ok {
		return s
	}
	if len(options) > 0 {
		return options[0]
	}
	return nil
}
func QueryTime(v url.Values, name string, options...time.Time) *time.Time {
	s := QueryString(v, name)
	if len(s) == 0 {
		t := CreateTime(s)
		if t != nil {
			return t
		}
	}
	if len(options) > 0 {
		return &options[0]
	}
	return nil
}
func QueryInt64(v url.Values, name string, options...int64) *int64 {
	s := QueryString(v, name)
	if len(s) == 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil
		}
		return &i
	}
	if len(options) > 0 {
		return &options[0]
	}
	return nil
}
func QueryInt32(v url.Values, name string, options...int64) *int32 {
	i := QueryInt64(v, name, options...)
	if i != nil {
		j := int32(*i)
		return &j
	}
	return nil
}
func QueryInt(v url.Values, name string, options...int64) *int {
	i := QueryInt64(v, name, options...)
	if i != nil {
		j := int(*i)
		return &j
	}
	return nil
}
func QueryRequiredString(w http.ResponseWriter, v url.Values, name string) string {
	s := QueryString(v, name)
	if len(s) == 0 {
		http.Error(w, fmt.Sprintf("%s is required", name), http.StatusBadRequest)
	}
	return s
}
func QueryRequiredStrings(w http.ResponseWriter, v url.Values, name string, options...string) []string {
	s := QueryString(v, name)
	if len(s) == 0 {
		http.Error(w, fmt.Sprintf("%s is required", name), http.StatusBadRequest)
		return nil
	} else {
		if len(options) > 0 && len(options[0]) > 0 {
			return strings.Split(s, options[0])
		} else {
			return strings.Split(s, ",")
		}
	}
}
func QueryRequiredTime(w http.ResponseWriter, s url.Values, name string) *time.Time {
	v := QueryTime(s, name)
	if v == nil {
		http.Error(w, fmt.Sprintf("%s is a required time", name), http.StatusBadRequest)
		return nil
	}
	return v
}
func QueryRequiredInt64(w http.ResponseWriter, s url.Values, name string) *int64 {
	v := QueryInt64(s, name)
	if v == nil {
		http.Error(w, fmt.Sprintf("%s is a required integer", name), http.StatusBadRequest)
		return nil
	}
	return v
}
func QueryRequiredInt32(w http.ResponseWriter, s url.Values, name string) *int32 {
	v := QueryInt32(s, name)
	if v == nil {
		http.Error(w, fmt.Sprintf("%s is a required integer", name), http.StatusBadRequest)
		return nil
	}
	return v
}
func QueryRequiredInt(w http.ResponseWriter, s url.Values, name string) *int {
	v := QueryInt(s, name)
	if v == nil {
		http.Error(w, fmt.Sprintf("%s is a required integer", name), http.StatusBadRequest)
		return nil
	}
	return v
}
