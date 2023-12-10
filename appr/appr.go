package appr

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/core-go/core/response"
)

type ApprService interface {
	Approve(ctx context.Context, id string, userId string, note string) (int64, error)
	Reject(ctx context.Context, id string, userId string, note string) (int64, error)
}
func NewApprService(service ApprService, resource string, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, index int, opts...string) *ApprHandler {
	var user, action1, action2 string
	if len(opts) > 0 {
		user = opts[0]
	} else {
		user = "userId"
	}
	if len(opts) > 1 {
		action1 = opts[1]
	} else {
		action1 = "approve"
	}
	if len(opts) > 2 {
		action2 = opts[2]
	} else {
		action2 = "reject"
	}
	return &ApprHandler{ApprService: service, Resource: resource, Index: index, User: user, ActionApprove: action1, ActionReject:  action2, LogError: logError, WriteLog: writeLog}
}
var internalServerError = "Internal Server Error"
type ApprHandler struct {
	ApprService   ApprService
	Resource      string
	Index         int
	User          string
	ActionApprove string
	ActionReject  string
	LogError      func(context.Context, string, ...map[string]interface{})
	WriteLog      func(context.Context, string, string, bool, string) error
}

func (h *ApprHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, userId, note, ok := GetParameters(w, r, h.Index, h.User)
	if ok {
		res, err := h.ApprService.Approve(r.Context(), id, userId, note)
		response.HandleResult(w, r, id, res, err, h.Resource, h.ActionApprove, h.LogError, h.WriteLog)
	}
}
func (h *ApprHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, userId, note, ok := GetParameters(w, r, h.Index, h.User)
	if ok {
		res, err := h.ApprService.Reject(r.Context(), id, userId, note)
		response.HandleResult(w, r, id, res, err, h.Resource, h.ActionReject, h.LogError, h.WriteLog)
	}
}
func GetParams(w http.ResponseWriter, r *http.Request, index int, opts...string) (string, string, string, int32, bool) {
	id := GetRequiredParam(w, r, index)
	if len(id) > 0 {
		var userId = UserId
		if len(opts) > 0 {
			userId = opts[0]
		}
		user, ok := RequireUser(r.Context(), w, userId)
		if ok {
			note, version, ok2 := GetBodyWithVersion(w, r)
			return id, user, note, version, ok2
		} else {
			return id, user, "", 0, false
		}
	}
	return id, "", "", 0, false
}
func GetParameters(w http.ResponseWriter, r *http.Request, index int, opts...string) (string, string, string, bool) {
	id := GetRequiredParam(w, r, index)
	if len(id) > 0 {
		var userId = UserId
		if len(opts) > 0 {
			userId = opts[0]
		}
		user, ok := RequireUser(r.Context(), w, userId)
		if ok {
			note, ok2 := GetBodyField(w, r)
			return id, user, note, ok2
		} else {
			return id, user, "", false
		}
	}
	return id, "", "", false
}
func GetBodyField(w http.ResponseWriter, r *http.Request) (string, bool) {
	body := make(map[string]interface{})
	er1 := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return "", false
	} else {
		if len(body) == 0 {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return "", false
		} else if len(body) > 1 {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return "", false
		} else {
			for _, k := range body {
				u2, ok2 := k.(string)
				if ok2 {
					return u2, true
				} else {
					http.Error(w, "invalid body", http.StatusBadRequest)
					return "", false
				}
			}
			http.Error(w, "invalid body", http.StatusBadRequest)
			return "", false
		}
	}
}
func GetBodyWithVersion(w http.ResponseWriter, r *http.Request) (string, int32, bool) {
	body := make(map[string]interface{})
	er1 := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return "", 0, false
	} else {
		if len(body) != 2 {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return "", 0, false
		} else {
			var note string
			var version int32
			count := 0
			for _, k := range body {
				u2, ok2 := k.(string)
				if ok2 {
					note = u2
					count = count + 1
				} else {
					u3, ok3 := k.(float64)
					if ok3 {
						version = int32(u3)
						count = count + 1
					}
				}
			}
			if count != 2 {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return note, version, false
			} else {
				return note, version, true
			}
		}
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
