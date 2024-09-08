package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	t1 = "2006-01-02T15:04:05Z"
	t2 = "2006-01-02T15:04:05-0700"
	t3 = "2006-01-02T15:04:05.0000000-0700"

	l1 = len(t1)
	l2 = len(t2)
	l3 = len(t3)
)

func Decode(w http.ResponseWriter, r *http.Request, obj interface{}, options ...func(context.Context, interface{}) error) error {
	er1 := json.NewDecoder(r.Body).Decode(obj)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return er1
	}
	if len(options) > 0 && options[0] != nil {
		er2 := options[0](r.Context(), obj)
		if er2 != nil {
			http.Error(w, er2.Error(), http.StatusInternalServerError)
		}
		return er2
	}
	return nil
}
func GetParam(r *http.Request, opts ...int) string {
	offset := 0
	if len(opts) > 0 && opts[0] > 0 {
		offset = opts[0]
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
func GetRequiredParam(w http.ResponseWriter, r *http.Request, opts ...int) (string, error) {
	p := GetParam(r, opts...)
	if len(p) == 0 {
		se := "parameter is required"
		http.Error(w, se, http.StatusBadRequest)
		return p, errors.New(se)
	}
	return p, nil
}
func GetRequiredInt(w http.ResponseWriter, r *http.Request, opts ...int) (int, error) {
	p := GetParam(r, opts...)
	if len(p) == 0 {
		se := "parameter is required"
		http.Error(w, se, http.StatusBadRequest)
		return 0, errors.New(se)
	}
	i, err := strconv.Atoi(p)
	if err != nil {
		http.Error(w, "parameter must be an integer", http.StatusBadRequest)
		return 0, err
	}
	return i, nil
}
func GetRequiredInt64(w http.ResponseWriter, r *http.Request, opts ...int) (int64, error) {
	p := GetParam(r, opts...)
	if len(p) == 0 {
		se := "parameter is required"
		http.Error(w, se, http.StatusBadRequest)
		return 0, errors.New(se)
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		http.Error(w, "parameter must be an integer", http.StatusBadRequest)
		return 0, err
	}
	return i, nil
}
func GetRequiredUint64(w http.ResponseWriter, r *http.Request, opts ...int) (uint64, error) {
	p := GetParam(r, opts...)
	if len(p) == 0 {
		se := "parameter is required"
		http.Error(w, se, http.StatusBadRequest)
		return 0, errors.New(se)
	}
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		http.Error(w, "parameter must be an unsigned integer", http.StatusBadRequest)
		return 0, err
	}
	return i, nil
}
func GetRequiredInt32(w http.ResponseWriter, r *http.Request, opts ...int) (int32, error) {
	p := GetParam(r, opts...)
	if len(p) == 0 {
		se := "parameter is required"
		http.Error(w, se, http.StatusBadRequest)
		return 0, errors.New(se)
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		http.Error(w, "parameter must be an integer", http.StatusBadRequest)
		return 0, err
	}
	return int32(i), nil
}
func GetRequiredParams(w http.ResponseWriter, r *http.Request, opts ...int) []string {
	p := GetParam(r, opts...)
	if len(p) == 0 {
		http.Error(w, "parameters are required", http.StatusBadRequest)
		return nil
	}
	return strings.Split(p, ",")
}
func GetParams(r *http.Request, opts ...int) []string {
	p := GetParam(r, opts...)
	return strings.Split(p, ",")
}
func GetInt(r *http.Request, opts ...int) (int, bool) {
	s := GetParam(r, opts...)
	if len(s) == 0 {
		return 0, false
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return i, true
}
func GetInt64(r *http.Request, opts ...int) (int64, bool) {
	s := GetParam(r, opts...)
	if len(s) == 0 {
		return 0, false
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}
func GetInt32(r *http.Request, opts ...int) (int32, bool) {
	s := GetParam(r, opts...)
	if len(s) == 0 {
		return 0, false
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return int32(i), true
}
func GetTime(r *http.Request, opts ...int) *time.Time {
	s := GetParam(r, opts...)
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
func QueryString(v url.Values, name string, opts ...string) string {
	s := v.Get(name)
	if len(s) > 0 {
		return s
	}
	if len(opts) > 0 {
		return opts[0]
	}
	return ""
}
func QueryTime(v url.Values, name string, opts ...time.Time) *time.Time {
	s := QueryString(v, name)
	if len(s) > 0 {
		t := CreateTime(s)
		if t != nil {
			return t
		}
	}
	if len(opts) > 0 {
		return &opts[0]
	}
	return nil
}
func QueryRequiredTime(w http.ResponseWriter, s url.Values, name string) *time.Time {
	v := QueryTime(s, name)
	if v == nil {
		http.Error(w, fmt.Sprintf("%s is a required time", name), http.StatusBadRequest)
		return nil
	}
	return v
}
func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
