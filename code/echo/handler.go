package echo

import (
	"context"
	"encoding/json"
	co "github.com/core-go/core/code"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const internalServerError = "Internal Server Error"

type Handler struct {
	Codes          func(ctx context.Context, master string) ([]co.Model, error)
	RequiredMaster bool
	Error          func(context.Context, string, ...map[string]interface{})
	Log            func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource       string
	Action         string
	Id             string
	Name           string
}

func NewDefaultCodeHandler(load func(ctx context.Context, master string) ([]co.Model, error), logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *Handler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewCodeHandlerWithLog(load, logError, true, writeLog, "", "")
}
func NewCodeHandlerByConfig(load func(ctx context.Context, master string) ([]co.Model, error), c co.HandlerConfig, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *Handler {
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
func NewCodeHandler(load func(ctx context.Context, master string) ([]co.Model, error), logError func(context.Context, string, ...map[string]interface{}), requiredMaster bool, options ...func(context.Context, string, string, bool, string) error) *Handler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewCodeHandlerWithLog(load, logError, requiredMaster, writeLog, "", "")
}
func NewCodeHandlerWithLog(load func(ctx context.Context, master string) ([]co.Model, error), logError func(context.Context, string, ...map[string]interface{}), requiredMaster bool, writeLog func(context.Context, string, string, bool, string) error, options ...string) *Handler {
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
func (h *Handler) Load(ctx echo.Context) error {
	r := ctx.Request()
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
				ctx.String(http.StatusBadRequest, "Body cannot is empty")
				return er1
			}
			code = strings.Trim(string(b), " ")
		}
	}
	result, er4 := h.Codes(r.Context(), code)
	if er4 != nil {
		return respondError(ctx, http.StatusInternalServerError, internalServerError, h.Error, h.Resource, h.Action, er4, h.Log)
	} else {
		if len(h.Id) == 0 && len(h.Name) == 0 {
			return succeed(ctx, http.StatusOK, result, h.Log, h.Resource, h.Action)
		} else {
			rs := make([]map[string]string, 0)
			for _, r := range result {
				m := make(map[string]string)
				m[h.Id] = r.Id
				m[h.Name] = r.Name
				rs = append(rs, m)
			}
			return succeed(ctx, http.StatusOK, rs, h.Log, h.Resource, h.Action)
		}
	}
}

type QueryHandler struct {
	Get      func(ctx context.Context, key string, max int64) ([]co.Model, error)
	Select   func(ctx context.Context, key []string) ([]co.Model, error)
	LogError func(context.Context, string, ...map[string]interface{})
	Keyword  string
	Max      string
	Q        string
}

func NewQueryHandler(load func(ctx context.Context, key string, max int64) ([]co.Model, error), getData func(ctx context.Context, key []string) ([]co.Model, error), logError func(context.Context, string, ...map[string]interface{}), opts ...string) *QueryHandler {
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

func (h *QueryHandler) Query(ctx echo.Context) error {
	ps := ctx.Request().URL.Query()
	keyword := ps.Get(h.Keyword)
	if len(keyword) == 0 {
		vs := make([]string, 0)
		return ctx.JSON(http.StatusOK, vs)
	} else {
		max := ps.Get(h.Max)
		i, err := strconv.ParseInt(max, 10, 64)
		if err != nil {
			i = 20
		}
		if i < 0 {
			i = 20
		}
		vs, err := h.Get(ctx.Request().Context(), keyword, i)
		if err != nil {
			h.LogError(ctx.Request().Context(), err.Error())
			return ctx.String(http.StatusInternalServerError, internalServerError)
		} else {
			return ctx.JSON(http.StatusOK, vs)
		}
	}
}
func (h *QueryHandler) Load(ctx echo.Context) error {
	r := ctx.Request()
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
			return ctx.String(http.StatusBadRequest, err.Error())
		}
	}
	if len(req) == 0 {
		return ctx.JSON(http.StatusOK, req)
	}
	models, err := h.Select(r.Context(), req)
	if err != nil {
		h.LogError(r.Context(), err.Error())
		return ctx.String(http.StatusInternalServerError, internalServerError)
	} else {
		return ctx.JSON(http.StatusOK, models)
	}
}
func respond(ctx echo.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) error {
	err := ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request().Context(), resource, action, success, desc)
	}
	return err
}
func respondError(ctx echo.Context, code int, result interface{}, logError func(context.Context, string, ...map[string]interface{}), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) error {
	if logError != nil {
		logError(ctx.Request().Context(), err.Error())
	}
	respond(ctx, code, result, writeLog, resource, action, false, err.Error())
	return err
}
func succeed(ctx echo.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return respond(ctx, code, result, writeLog, resource, action, true, "")
}
