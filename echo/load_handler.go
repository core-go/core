package echo

import (
	"context"
	sv "github.com/core-go/service"
	"github.com/labstack/echo/v4"
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
func (h *LoadHandler) Load(ctx echo.Context) error {
	r := ctx.Request()
	id, er1 := sv.BuildId(r, h.ModelType, h.Keys, h.Indexes, 0)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return er1
	} else {
		model, er2 := h.LoadData(r.Context(), id)
		return RespondModel(ctx, model, er2, h.Error, h.WriteLog, h.Resource, h.Action)
	}
}

func RespondModel(ctx echo.Context, model interface{}, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	if err != nil {
		return Respond(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	} else {
		if model == nil {
			return RespondAndLog(ctx, http.StatusNotFound, model, writeLog, false, resource, action, "Not found")
		} else {
			return Succeed(ctx, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
func Respond(ctx echo.Context, code int, result interface{}, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options... string) error {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		if logError != nil {
			logError(ctx.Request().Context(), err.Error())
			return RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, writeLog, false, resource, action, err.Error())
		} else {
			return RespondAndLog(ctx, http.StatusInternalServerError, err.Error(), writeLog, false, resource, action, err.Error())
		}
	} else {
		return RespondAndLog(ctx, code, result, writeLog, true, resource, action, "")
	}
}
func RespondAndLog(ctx echo.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) error {
	err := ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request().Context(), resource, action, success, desc)
	}
	return err
}

func ErrorAndLog(ctx echo.Context, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) error {
	if logError != nil {
		logError(ctx.Request().Context(), err.Error())
	}
	RespondAndLog(ctx, code, result, writeLog, false, resource, action, err.Error())
	return err
}
func Succeed(ctx echo.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return RespondAndLog(ctx, code, result, writeLog, true, resource, action, "")
}
