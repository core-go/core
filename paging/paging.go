package paging

import (
	"fmt"
	"net/http"
	"strconv"
)
func GetInt64(w http.ResponseWriter, r *http.Request, defaultValue int64, name string) (int64, error) {
	ps := r.URL.Query()
	slimit := ps.Get(name)
	var limit int64
	if len(slimit) > 0 {
		l1, err := strconv.ParseInt(slimit, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s must be an integer", name), http.StatusBadRequest)
			return 0, err
		}
		limit = l1
	}
	if limit <= 0 {
		limit = defaultValue
	}
	return limit, nil
}
func GetPage(w http.ResponseWriter, r *http.Request, opts...string) (int64, error) {
	var pageName string
	pageName = "page"
	if len(opts) > 0 && len(opts[0]) > 0 {
		pageName = opts[0]
	}
	return GetInt64(w, r, 1, pageName)
}
func GetOffset(w http.ResponseWriter, r *http.Request, defaultLimit int64, opts...string) (int64, int64, error) {
	var offsetName, limitName string
	offsetName = "offset"
	if len(opts) > 0 && len(opts[0]) > 0 {
		offsetName = opts[0]
	}
	limitName = "limit"
	if len(opts) > 1 && len(opts[1]) > 0 {
		limitName = opts[1]
	}
	limit, err := GetInt64(w, r, defaultLimit, limitName)
	if err != nil {
		return 0, 0, err
	}
	offset, err := GetInt64(w, r, 0, offsetName)
	return limit, offset, err
}
func GetNext(w http.ResponseWriter, r *http.Request, defaultLimit int64, opts...string) (int64, string, error) {
	var nextName, limitName string
	nextName = "next"
	if len(opts) > 0 && len(opts[0]) > 0 {
		nextName = opts[0]
	}
	limitName = "limit"
	if len(opts) > 1 && len(opts[1]) > 0 {
		limitName = opts[1]
	}
	limit, err := GetInt64(w, r, defaultLimit, limitName)
	if err != nil {
		return 0, "", err
	}
	ps := r.URL.Query()
	return limit, ps.Get(nextName), err
}
