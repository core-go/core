package code

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const internalServerError = "Internal Server Error"

type HandlerConfig struct {
	Master   *bool  `mapstructure:"master" json:"master,omitempty" gorm:"column:master" bson:"master,omitempty" dynamodbav:"master,omitempty" firestore:"master,omitempty"`
	Id       string `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name     string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Resource string `mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Action   string `mapstructure:"action" json:"action,omitempty" gorm:"column:action" bson:"action,omitempty" dynamodbav:"action,omitempty" firestore:"action,omitempty"`
}
type Handler struct {
	Codes          func(ctx context.Context, master string) ([]Model, error)
	RequiredMaster bool
	Error          func(context.Context, string)
	Log            func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource       string
	Action         string
	Id             string
	Name           string
}

func NewDefaultCodeHandler(load func(ctx context.Context, master string) ([]Model, error), logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *Handler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewCodeHandlerWithLog(load, logError, true, writeLog, "", "")
}
func NewCodeHandlerByConfig(load func(ctx context.Context, master string) ([]Model, error), c HandlerConfig, logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *Handler {
	var requireMaster bool
	if c.Master != nil {
		requireMaster = *c.Master
	} else {
		requireMaster = true
	}
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	h := NewCodeHandlerWithLog(load, logError, requireMaster, writeLog, c.Resource, c.Action)
	h.Id = c.Id
	h.Name = c.Name
	return h
}
func NewCodeHandler(load func(ctx context.Context, master string) ([]Model, error), logError func(context.Context, string), requiredMaster bool, options ...func(context.Context, string, string, bool, string) error) *Handler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewCodeHandlerWithLog(load, logError, requiredMaster, writeLog, "", "")
}
func NewCodeHandlerWithLog(load func(ctx context.Context, master string) ([]Model, error), logError func(context.Context, string), requiredMaster bool, writeLog func(context.Context, string, string, bool, string) error, options ...string) *Handler {
	var resource, action string
	if len(options) >= 1 && len(options[0]) > 0 {
		resource = options[0]
	} else {
		resource = "code"
	}
	if len(options) >= 2 && len(options[1]) > 0 {
		action = options[1]
	} else {
		action = "load"
	}
	h := Handler{Codes: load, Resource: resource, Action: action, RequiredMaster: requiredMaster, Log: writeLog, Error: logError}
	return &h
}
func (h *Handler) Load(w http.ResponseWriter, r *http.Request) {
	code := ""
	if h.RequiredMaster {
		if r.Method == "GET" {
			i := strings.LastIndex(r.RequestURI, "/")
			if i >= 0 {
				code = r.RequestURI[i+1:]
			}
		} else {
			b, er1 := ioutil.ReadAll(r.Body)
			if er1 != nil {
				http.Error(w, "Body cannot is empty", http.StatusBadRequest)
				return
			}
			code = strings.Trim(string(b), " ")
		}
	}
	result, er4 := h.Codes(r.Context(), code)
	if er4 != nil {
		respondError(w, r, http.StatusInternalServerError, internalServerError, h.Error, h.Resource, h.Action, er4, h.Log)
	} else {
		if len(h.Id) == 0 && len(h.Name) == 0 {
			succeed(w, r, http.StatusOK, result, h.Log, h.Resource, h.Action)
		} else {
			rs := make([]map[string]string, 0)
			for _, r := range result {
				m := make(map[string]string)
				m[h.Id] = r.Id
				m[h.Name] = r.Name
				rs = append(rs, m)
			}
			succeed(w, r, http.StatusOK, rs, h.Log, h.Resource, h.Action)
		}
	}
}
func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
}
func respondError(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	respond(w, r, code, result, writeLog, resource, action, false, err.Error())
}
func succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	respond(w, r, code, result, writeLog, resource, action, true, "")
}
