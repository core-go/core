package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	Update   func(ctx context.Context, id string, settings Settings) (int64, error)
	LogError func(context.Context, string, ...map[string]interface{})
	WriteLog func(context.Context, string, string, bool, string) error
	User     string
	Resource string
	Action   string
}
func NewSettingsHandler(logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, db *sql.DB, table string, buildParam func(int)string, opts...string) *Handler {
	var user, resource, action string
	var id, language, dateFormat string
	if len(opts) > 0 {
		user = opts[0]
	} else {
		user = "userId"
	}

	if len(opts) > 1 {
		id = opts[1]
	} else {
		id = "id"
	}
	if len(opts) > 2 {
		dateFormat = opts[2]
	} else {
		dateFormat = "dateformat"
	}
	if len(opts) > 3 {
		language = opts[3]
	} else {
		language = "language"
	}

	if len(opts) > 4 {
		resource = opts[4]
	} else {
		resource = "settings"
	}
	if len(opts) > 5 {
		action = opts[5]
	} else {
		action = "save"
	}
	service := NewSettingsService(db, buildParam, table, id, dateFormat, language)
	return NewHandler(service.Save, logError, writeLog, user, resource, action)
}
func NewHandler(save func(context.Context, string, Settings) (int64, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts...string) *Handler {
	var user, resource, action string
	if len(opts) > 0 {
		user = opts[0]
	} else {
		user = "userId"
	}
	if len(opts) > 1 {
		resource = opts[1]
	} else {
		resource = "settings"
	}
	if len(opts) > 2 {
		action = opts[2]
	} else {
		action = "save"
	}
	return &Handler{Update: save, LogError: logError, WriteLog: writeLog, User: user, Resource: resource, Action: action}
}
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	var settings Settings
	er1 := json.NewDecoder(r.Body).Decode(&settings)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	id, ok := RequireUser(r.Context(), w, h.User)
	if ok {
		res, err := h.Update(r.Context(), id, settings)
		if err != nil {
			if h.LogError != nil {
				h.LogError(r.Context(), err.Error(), MakeMap(settings, "settings"))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if h.WriteLog != nil {
				h.WriteLog(r.Context(), h.Resource, h.Action, false, err.Error())
			}
		} else {
			if res > 0 {
				if h.WriteLog != nil {
					h.WriteLog(r.Context(), h.Resource, h.Action, true, fmt.Sprintf("save settings for '%s'", id))
				}
				JSON(w, http.StatusOK, res)
			} else {
				JSON(w, http.StatusNotFound, res)
				if h.WriteLog != nil {
					h.WriteLog(r.Context(), h.Resource, h.Action, false, fmt.Sprintf("not found '%s'", id))
				}
			}
		}
	}
}

func JSON(w http.ResponseWriter, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	return err
}
func GetUser(ctx context.Context, user string) (string, bool) {
	u := ctx.Value(user)
	if u != nil {
		u2, ok2 := u.(string)
		if ok2 {
			return u2, ok2
		}
	}
	return "", false
}
func RequireUser(ctx context.Context, w http.ResponseWriter, user string) (string, bool) {
	userId, ok := GetUser(ctx, user)
	if ok {
		return userId, ok
	} else {
		http.Error(w, "cannot get current user", http.StatusForbidden)
		return "", false
	}
}
func MakeMap(res interface{}, opts ...string) map[string]interface{} {
	key := "request"
	if len(opts) > 0 && len(opts[0]) > 0 {
		key = opts[0]
	}
	m := make(map[string]interface{})
	b, err := json.Marshal(res)
	if err != nil {
		return m
	}
	m[key] = string(b)
	return m
}
