package code

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const internalServerError = "Internal Server Error"

type HandlerConfig struct {
	Master   *bool  `yaml:"master" mapstructure:"master" json:"master,omitempty" gorm:"column:master" bson:"master,omitempty" dynamodbav:"master,omitempty" firestore:"master,omitempty"`
	Id       string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name     string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Resource string `yaml:"resource" mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Action   string `yaml:"action" mapstructure:"action" json:"action,omitempty" gorm:"column:action" bson:"action,omitempty" dynamodbav:"action,omitempty" firestore:"action,omitempty"`
}
type Handler struct {
	Codes          func(ctx context.Context, master string) ([]Model, error)
	RequiredMaster bool
	Error          func(context.Context, string, ...map[string]interface{})
	Log            func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource       string
	Action         string
	Id             string
	Name           string
}

func NewDefaultCodeHandler(load func(ctx context.Context, master string) ([]Model, error), logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *Handler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewCodeHandlerWithLog(load, logError, true, writeLog, "", "")
}
func NewCodeHandlerByConfig(load func(ctx context.Context, master string) ([]Model, error), c HandlerConfig, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *Handler {
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
func NewCodeHandler(load func(ctx context.Context, master string) ([]Model, error), logError func(context.Context, string, ...map[string]interface{}), requiredMaster bool, options ...func(context.Context, string, string, bool, string) error) *Handler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewCodeHandlerWithLog(load, logError, requiredMaster, writeLog, "", "")
}
func NewCodeHandlerWithLog(load func(ctx context.Context, master string) ([]Model, error), logError func(context.Context, string, ...map[string]interface{}), requiredMaster bool, writeLog func(context.Context, string, string, bool, string) error, options ...string) *Handler {
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
			b, er1 := io.ReadAll(r.Body)
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

type QueryHandler struct {
	Get      func(ctx context.Context, key string, max int64) ([]Model, error)
	Select   func(ctx context.Context, key []string) ([]Model, error)
	LogError func(context.Context, string, ...map[string]interface{})
	Keyword  string
	Max      string
	Q        string
}

func NewQueryHandler(load func(ctx context.Context, key string, max int64) ([]Model, error), getData func(ctx context.Context, key []string) ([]Model, error), logError func(context.Context, string, ...map[string]interface{}), opts ...string) *QueryHandler {
	q := "q"
	if len(opts) > 0 && len(opts[0]) > 0 {
		q = opts[0]
	}
	keyword := "q"
	if len(opts) > 1 && len(opts[1]) > 0 {
		keyword = opts[1]
	}
	max := "max"
	if len(opts) > 2 && len(opts[2]) > 0 {
		max = opts[2]
	}
	return &QueryHandler{load, getData, logError, keyword, max, q}
}
func (h *QueryHandler) Query(w http.ResponseWriter, r *http.Request) {
	ps := r.URL.Query()
	keyword := ps.Get(h.Keyword)
	if len(keyword) == 0 {
		vs := make([]string, 0)
		respondModel(w, r, vs, nil, h.LogError, nil)
	} else {
		max := ps.Get(h.Max)
		i, err := strconv.ParseInt(max, 10, 64)
		if err != nil {
			i = 20
		}
		if i < 0 {
			i = 20
		}
		vs, err := h.Get(r.Context(), keyword, i)
		respondModel(w, r, vs, err, h.LogError, nil)
	}
}
func (h *QueryHandler) Load(w http.ResponseWriter, r *http.Request) {
	var req = make([]string, 0)
	method := r.Method
	if method == http.MethodGet {
		q := r.URL.Query().Get(h.Q)
		if len(q) > 0 {
			req = strings.Split(q, ",")
		}
	} else {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if len(req) == 0 {
		respondModel(w, r, req, nil, h.LogError, nil)
	} else {
		models, err := h.Select(r.Context(), req)
		respondModel(w, r, models, err, h.LogError, nil)
	}
}
func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
}
func respondError(w http.ResponseWriter, r *http.Request, code int, result interface{}, logError func(context.Context, string, ...map[string]interface{}), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(r.Context(), err.Error())
	}
	respond(w, r, code, result, writeLog, resource, action, false, err.Error())
}
func succeed(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	respond(w, r, code, result, writeLog, resource, action, true, "")
}

func respondModel(w http.ResponseWriter, r *http.Request, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		respondAndLog(w, r, http.StatusInternalServerError, internalServerError, err, logError, writeLog, resource, action)
	} else {
		if model == nil {
			returnAndLog(w, r, http.StatusNotFound, model, writeLog, false, resource, action, "Not found")
		} else {
			succeed(w, r, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
func respondAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) error {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		if logError != nil {
			logError(r.Context(), err.Error())
			return returnAndLog(w, r, http.StatusInternalServerError, internalServerError, writeLog, false, resource, action, err.Error())
		} else {
			return returnAndLog(w, r, http.StatusInternalServerError, err.Error(), writeLog, false, resource, action, err.Error())
		}
	} else {
		return returnAndLog(w, r, code, result, writeLog, true, resource, action, "")
	}
}
func returnAndLog(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
	return err
}
