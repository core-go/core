package histories

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	Resource string
	Index    int
	Load func(ctx context.Context, resource string, id string, limit int64, offset int64) ([]History, int64, error)
	Limit string
	Offset string
	List string
	Total string
}

type GetHistories func(ctx context.Context, resource string, id string, limit int64, offset int64) ([]History, int64, error)

func NewHistoriesHandler(resource string, getHistories func(ctx context.Context, resource string, id string, limit int64, offset int64) ([]History, int64, error), opts...string) *Handler {
	var list, total, limit, offset string
	if len(opts) > 0 {
		list = opts[0]
	} else {
		list = "list"
	}
	if len(opts) > 1 {
		total = opts[1]
	} else {
		total = "total"
	}
	if len(opts) > 2 {
		limit = opts[2]
	} else {
		limit = "limit"
	}
	if len(opts) > 2 {
		offset = opts[2]
	} else {
		offset = "offset"
	}
	return &Handler{Resource: resource, Load: getHistories, List: list, Total: total, Limit: limit, Offset: offset}
}

func (h *Handler) GetHistories(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredParam(w, r, h.Index)
	if len(id) > 0 {
		ps := r.URL.Query()
		slimit := ps.Get(h.Limit)
		var limit, offset int64
		if len(slimit) > 0 {
			l1, err := strconv.ParseInt(slimit, 10, 64)
			if err != nil {
				http.Error(w, "limit must be an integer", http.StatusBadRequest)
				return
			}
			limit = l1
		}
		if limit <= 0 {
			limit = 20
		}
		soffset := ps.Get(h.Offset)
		if len(soffset) > 0 {
			o1, err := strconv.ParseInt(soffset, 10, 64)
			if err != nil {
				http.Error(w, "offset must be an integer", http.StatusBadRequest)
				return
			}
			offset = o1
		}
		if offset < 0 {
			offset = 0
		}
		res, total, err := h.Load(r.Context(), h.Resource, id, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := make(map[string]interface{})
		m[h.List] = res
		m[h.Total] = total
		JSON(w, 200, m)
	}
}
func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
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
func GetRequiredParam(w http.ResponseWriter,r *http.Request, options ...int) string {
	p := GetParam(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return ""
	}
	return p
}
func GetRequiredInt64(w http.ResponseWriter,r *http.Request, options ...int) *int64 {
	p := GetParam(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return nil
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		http.Error(w, "parameter must be an integer", http.StatusBadRequest)
		return nil
	}
	return &i
}