package gin

import (
	"context"
	sv "github.com/core-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type LoadHandler struct {
	LoadData  func(ctx context.Context, id interface{}) (interface{}, error)
	Keys      []string
	ModelType reflect.Type
	Indexes   map[string]int
	Error     func(context.Context, string)
	WriteLog  func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource  string
	Action    string
}

func NewLoadHandler(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewLoadHandlerWithLog(load, modelType, logError, writeLog)
}
func NewLoadHandlerWithKeys(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *LoadHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewLoadHandlerWithKeysAndLog(load, keys, modelType, logError, writeLog)
}
func NewLoadHandlerWithLog(load func(context.Context, interface{}) (interface{}, error), modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
	return NewLoadHandlerWithKeysAndLog(load, nil, modelType, logError, writeLog, options...)
}
func NewLoadHandlerWithKeysAndLog(load func(context.Context, interface{}) (interface{}, error), keys []string, modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LoadHandler {
	if keys == nil || len(keys) == 0 {
		keys = sv.GetJsonPrimaryKeys(modelType)
	}
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		action = options[0]
	} else {
		action = "load"
	}
	if len(options) > 1 && len(options[1]) > 0 {
		resource = options[1]
	} else {
		resource = sv.BuildResourceName(modelType.Name())
	}
	indexes := sv.GetIndexes(modelType)
	return &LoadHandler{WriteLog: writeLog, LoadData: load, Keys: keys, ModelType: modelType, Indexes: indexes, Error: logError, Resource: resource, Action: action}
}
func (h *LoadHandler) Load(ctx *gin.Context) {
	r := ctx.Request
	id, er1 := sv.BuildId(r, h.ModelType, h.Keys, h.Indexes, 0)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return
	} else {
		model, er2 := h.LoadData(r.Context(), id)
		if er2 != nil {
			ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action, er2, h.WriteLog)
		} else {
			if model == nil {
				Respond(ctx, http.StatusNotFound, model, h.WriteLog, h.Resource, h.Action, false, "Not found")
			} else {
				Succeed(ctx, http.StatusOK, model, h.WriteLog, h.Resource, h.Action)
			}
		}
	}
}

func Respond(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request.Context(), resource, action, success, desc)
	}
}

func RespondAndLog(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) {
	ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request.Context(), resource, action, success, desc)
	}
}

func ErrorAndLog(ctx *gin.Context, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(ctx.Request.Context(), err.Error())
	}
	RespondAndLog(ctx, code, result, writeLog, resource, action, false, err.Error())
}
func Succeed(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	RespondAndLog(ctx, code, result, writeLog, resource, action, true, "")
}
