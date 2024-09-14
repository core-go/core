package histories

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/core-go/core/paging"
)

type Handler struct {
	resource      string
	index         int
	load          func(ctx context.Context, resource string, id string, limit int64, nextPageToken string) ([]History, string, error)
	logError      func(context.Context, string, ...map[string]interface{})
	Limit         string
	NextPageToken string
	List          string
	Next          string
}

type GetHistories func(ctx context.Context, resource string, id string, limit int64, offset int64) ([]History, int64, error)

func NewHistoriesHandler(resource string, index int, getHistories func(ctx context.Context, resource string, id string, limit int64, nextPageToken string) ([]History, string, error), logError func(context.Context, string, ...map[string]interface{}), opts ...string) *Handler {
	var list, next, limit, nextPageToken string
	if len(opts) > 0 {
		list = opts[0]
	} else {
		list = "list"
	}
	if len(opts) > 1 {
		next = opts[1]
	} else {
		next = "next"
	}
	if len(opts) > 2 {
		limit = opts[2]
	} else {
		limit = "limit"
	}
	if len(opts) > 3 {
		nextPageToken = opts[3]
	} else {
		nextPageToken = "next"
	}
	return &Handler{resource: resource, index: index, load: getHistories, logError: logError, List: list, Next: next, Limit: limit, NextPageToken: nextPageToken}
}

func (h *Handler) GetHistories(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredString(w, r, h.index)
	if len(id) > 0 {
		limit, nextPageToken, err := paging.GetNext(w, r, 20, h.NextPageToken, h.Limit)
		if err == nil {
			res, next, err := h.load(r.Context(), h.resource, id, limit, nextPageToken)
			if err != nil {
				if h.logError != nil {
					h.logError(r.Context(), err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			m := make(map[string]interface{})
			m[h.List] = res
			m[h.Next] = next
			JSON(w, 200, m)
		}
	}
}
func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	return err
}

func GetString(r *http.Request, options ...int) string {
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
func GetRequiredString(w http.ResponseWriter, r *http.Request, options ...int) string {
	p := GetString(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return ""
	}
	return p
}
