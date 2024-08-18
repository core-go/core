package core

import (
	"encoding/json"
	"net/http"
	"reflect"
)

const InternalServerError = "Internal Server Error"

func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	return err
}
func GetStatus(status int64, opts ...int) int {
	if status > 0 {
		if len(opts) > 0 {
			return opts[0]
		}
		return http.StatusOK
	}
	if status == 0 {
		if len(opts) > 1 {
			return opts[1]
		}
		return http.StatusNotFound
	}
	if len(opts) > 2 {
		return opts[2]
	} else if len(opts) > 1 {
		return opts[1]
	}
	return http.StatusConflict
}
func IsFound(res interface{}) int {
	if IsNil(res) {
		return http.StatusNotFound
	}
	return http.StatusOK
}
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
func MakeMap(res interface{}, opts ...string) map[string]interface{} {
	key := "request"
	if len(opts) > 0 && len(opts[0]) > 0 {
		key = opts[0]
	}
	m0, ok0 := res.(map[string]interface{})
	if ok0 {
		return m0
	}
	m := make(map[string]interface{})
	b, err := json.Marshal(res)
	if err != nil {
		return m
	}
	m[key] = string(b)
	return m
}
