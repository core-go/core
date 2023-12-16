package notifications

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/core-go/core/paging"
)

type Handler struct {
	Service  NotificationsPort
	LogError func(context.Context, string, ...map[string]interface{})
	UserId   string
	Read     string
	Limit    string
	NextPageToken string
	List     string
	Next     string
}

func NewNotificationsHandler(service NotificationsPort, logError func(context.Context, string, ...map[string]interface{}), opts...string) *Handler {
	var userId, read, list, next, limit, nextPageToken string
	if len(opts) > 1 {
		userId = opts[1]
	} else {
		userId = "userId"
	}
	if len(opts) > 1 {
		read = opts[1]
	} else {
		read = "read"
	}
	if len(opts) > 2 {
		list = opts[2]
	} else {
		list = "list"
	}
	if len(opts) > 3 {
		next = opts[3]
	} else {
		next = "next"
	}
	if len(opts) > 4 {
		limit = opts[4]
	} else {
		limit = "limit"
	}
	if len(opts) > 5 {
		nextPageToken = opts[5]
	} else {
		nextPageToken = "next"
	}
	return &Handler{Service: service, LogError: logError, UserId: userId, Read: read, List: list, Next: next, Limit: limit, NextPageToken: nextPageToken}
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	id, ok := RequireUser(r.Context(), w, h.UserId)
	if ok {
		ps := r.URL.Query()
		sread := ps.Get(h.Read)
		var read *bool
		if sread == "true" {
			b := true
			read = &b
		}
		limit, nextPageToken, err := paging.GetNext(w, r, 20, h.NextPageToken, h.Limit)
		res, next, err := h.Service.GetNotifications(r.Context(), id, read, limit, nextPageToken)
		if err != nil {
			if h.LogError != nil {
				h.LogError(r.Context(), err.Error())
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
func (h *Handler) SetRead(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.Service.SetRead(r.Context(), id, true)
		if err != nil {
			if h.LogError != nil {
				h.LogError(r.Context(), err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if res > 0 {
			JSON(w, http.StatusOK, res)
		} else {
			JSON(w, http.StatusNotFound, res)
		}
	}
}
func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	return err
}
var UserId = "userId"
func ApplyUserId(str string) {
	UserId = str
}
func GetUser(ctx context.Context, opt...string) (string, bool) {
	user := UserId
	if len(opt) > 0 && len(opt[0]) > 0 {
		user = opt[0]
	}
	u := ctx.Value(user)
	if u != nil {
		u2, ok2 := u.(string)
		if ok2 {
			return u2, ok2
		}
	}
	return "", false
}
func RequireUser(ctx context.Context, w http.ResponseWriter, opt...string) (string, bool) {
	userId, ok := GetUser(ctx, opt...)
	if ok {
		return userId, ok
	} else {
		http.Error(w, "cannot get current user", http.StatusForbidden)
		return "", false
	}
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
