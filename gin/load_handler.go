package gin

import (
	"context"
	sv "github.com/core-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type LoadHandler struct {
	LoadData   func(ctx context.Context, id interface{}) (interface{}, error)
	Keys       []string
	ModelType  reflect.Type
	KeyIndexes map[string]int
	Error      func(context.Context, string)
	WriteLog   func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource   string
	Activity   string
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
	indexes := sv.GetKeyIndexes(modelType)
	return &LoadHandler{WriteLog: writeLog, LoadData: load, Keys: keys, ModelType: modelType, KeyIndexes: indexes, Error: logError, Resource: resource, Activity: action}
}
func (h *LoadHandler) Load(ctx *gin.Context) {
	r := ctx.Request
	id, er1 := sv.BuildId(r, h.ModelType, h.Keys, h.KeyIndexes, 0)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return
	} else {
		model, er2 := h.LoadData(r.Context(), id)
		RespondModel(ctx, model, er2, h.Error, h.WriteLog, h.Resource, h.Activity)
	}
}

func RespondModel(ctx *gin.Context, model interface{}, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	} else {
		if model == nil {
			ReturnAndLog(ctx, http.StatusNotFound, model, writeLog, false, resource, action, "Not found")
		} else {
			Succeed(ctx, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
func RespondAndLog(ctx *gin.Context, code int, result interface{}, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options... string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		if logError != nil {
			logError(ctx.Request.Context(), err.Error())
			ReturnAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, writeLog, false, resource, action, err.Error())
		} else {
			ReturnAndLog(ctx, http.StatusInternalServerError, err.Error(), writeLog, false, resource, action, err.Error())
		}
	} else {
		ReturnAndLog(ctx, code, result, writeLog, true, resource, action, "")
	}
}
func ReturnAndLog(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) {
	ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request.Context(), resource, action, success, desc)
	}
}
func ErrorAndLog(ctx *gin.Context, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(ctx.Request.Context(), err.Error())
	}
	ReturnAndLog(ctx, code, result, writeLog, false, resource, action, err.Error())
}
func Succeed(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	ReturnAndLog(ctx, code, result, writeLog, true, resource, action, "")
}
