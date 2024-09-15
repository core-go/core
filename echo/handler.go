package echo

import (
	"context"
	"github.com/core-go/core"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

func CheckId[T any](ctx echo.Context, body *T, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) error {
	return core.CheckId[T](ctx.Response().Writer, ctx.Request(), body, keysJson, mapIndex, opts...)
}
func AfterDeletedWithLog(ctx echo.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterDeletedWithLog(ctx.Response().Writer, ctx.Request(), count, err, logError, writeLog, opts...)
}
func AfterDeleted(ctx echo.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterDeleted(ctx.Response().Writer, ctx.Request(), count, err, logError)
}
func HasError(ctx echo.Context, errors []core.ErrorMessage, err error, logError func(context.Context, string, ...map[string]interface{}), model interface{}, writeLog func(context.Context, string, string, bool, string) error, opts ...string) bool {
	return core.HasError(ctx.Response().Writer, ctx.Request(), errors, err, logError, model, writeLog, opts...)
}
func AfterSavedWithLog(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterSavedWithLog(ctx.Response().Writer, ctx.Request(), body, count, err, logError, writeLog, opts...)
}
func AfterSaved(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterSaved(ctx.Response().Writer, ctx.Request(), body, count, err, logError)
}
func AfterCreatedWithLog(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterCreatedWithLog(ctx.Response().Writer, ctx.Request(), body, count, err, logError, writeLog, opts...)
}
func AfterCreated(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterCreated(ctx.Response().Writer, ctx.Request(), body, count, err, logError)
}
func BuildFieldMapAndCheckId[T any](ctx echo.Context, keysJson []string, mapIndex map[string]int, ignorePatch bool, opts ...func(context.Context, *T) error) (T, map[string]interface{}, error) {
	var obj T
	if ignorePatch == false {
		r := ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), core.Method, core.Patch))
		ctx.SetRequest(r)
	}
	body, er0 := core.BuildMapAndStruct(ctx.Request(), &obj)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, er0.Error())
		return obj, body, er0
	}
	er1 := CheckId[T](ctx, &obj, keysJson, mapIndex, opts...)
	return obj, body, er1
}
func BuildMapAndCheckId[T any](ctx echo.Context, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) (T, map[string]interface{}, error) {
	obj, body, er0 := BuildFieldMapAndCheckId[T](ctx, keysJson, mapIndex, false, opts...)
	if er0 != nil {
		return obj, body, er0
	}
	jsonObj, er1 := core.BodyToJsonMap(ctx.Request(), &obj, body, keysJson, mapIndex)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	}
	return obj, jsonObj, er1
}

type Handler[T any, K any] struct {
	Service   core.Service[T, K]
	LogError  func(context.Context, string, ...map[string]interface{})
	Validate  func(context.Context, *T) ([]core.ErrorMessage, error)
	Keys      []string
	Indexes   map[string]int
	Resource  string
	ModelType reflect.Type
	Action    core.ActionConfig
	WriteLog  func(context.Context, string, string, bool, string) error
	IdMap     bool
	Builder   core.Builder[T]
}

func Decode[T any](c echo.Context, opts ...func(context.Context, *T) error) (T, error) {
	return core.Decode[T](c.Response().Writer, c.Request(), opts...)
}
func DecodeAndCheckId[T any](c echo.Context, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) (T, error) {
	return core.DecodeAndCheckId[T](c.Response().Writer, c.Request(), keysJson, mapIndex, opts...)
}
func Newhandler[T any, K any](
	service core.Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	opts ...core.Builder[T],
) *Handler[T, K] {
	return NewhandlerWithLog[T, K](service, logError, validate, nil, nil, opts...)
}
func NewhandlerWithLog[T any, K any](
	service core.Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	action *core.ActionConfig,
	writeLog func(context.Context, string, string, bool, string) error,
	opts ...core.Builder[T],
) *Handler[T, K] {
	var b core.Builder[T]
	if len(opts) > 0 && opts[0] != nil {
		b = opts[0]
	}
	var t T
	modelType := reflect.TypeOf(t)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	var k K
	kType := reflect.TypeOf(k)
	idMap := false
	if kType.Kind() == reflect.Map {
		idMap = true
	}
	resource := core.BuildResourceName(modelType.Name())
	keys, indexes, _ := core.BuildMapField(modelType)
	a := core.InitAction(action)
	return &Handler[T, K]{Service: service, LogError: logError, Validate: validate, Keys: keys, Indexes: indexes, Resource: resource, ModelType: modelType, Builder: b, Action: a, WriteLog: writeLog, IdMap: idMap}
}

func (h *Handler[T, K]) Load(c echo.Context) error {
	r := c.Request()
	id, ok, er1 := core.BuildId[K](r, h.ModelType, h.Keys, h.Indexes, h.IdMap)
	if er1 != nil {
		return c.String(http.StatusBadRequest, er1.Error())
	}
	if !ok {
		return c.String(http.StatusBadRequest, "Id type is not valid (Id type must be K)")
	}
	model, er2 := h.Service.Load(r.Context(), id)
	action := ""
	if h.Action.Load != nil {
		action = *h.Action.Load
	}
	return ReturnWithLog(c, model, er2, h.LogError, h.WriteLog, h.Resource, action)
}
func (h *Handler[T, K]) Create(c echo.Context) error {
	var createFn func(context.Context, *T) error
	if h.Builder != nil {
		createFn = h.Builder.Create
	}
	model, er1 := Decode[T](c, createFn)
	if er1 != nil {
		return er1
	}
	r := c.Request()
	if h.Validate != nil {
		errors, er2 := h.Validate(r.Context(), &model)
		if HasError(c, errors, er2, h.LogError, model, h.WriteLog, h.Resource, h.Action.Create) {
			return er2
		}
		res, er3 := h.Service.Create(r.Context(), &model)
		return AfterCreatedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
	} else {
		res, er3 := h.Service.Create(r.Context(), &model)
		return AfterCreatedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
	}
}
func (h *Handler[T, K]) Update(c echo.Context) error {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	model, er1 := DecodeAndCheckId[T](c, h.Keys, h.Indexes, updateFn)
	if er1 != nil {
		return er1
	}
	r := c.Request()
	if h.Validate != nil {
		errors, er2 := h.Validate(r.Context(), &model)
		if HasError(c, errors, er2, h.LogError, model, h.WriteLog, h.Resource, h.Action.Update) {
			return er2
		}
		res, er3 := h.Service.Update(r.Context(), &model)
		return AfterSavedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
	} else {
		res, er3 := h.Service.Update(r.Context(), &model)
		return AfterSavedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
	}
}
func (h *Handler[T, K]) Patch(c echo.Context) error {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	model, jsonObj, er1 := BuildMapAndCheckId[T](c, h.Keys, h.Indexes, updateFn)
	if er1 != nil {
		return er1
	}
	if h.Validate != nil {
		errors, er2 := h.Validate(c.Request().Context(), &model)
		if HasError(c, errors, er2, h.LogError, jsonObj, h.WriteLog, h.Resource, h.Action.Patch) {
			return er2
		}
		res, er3 := h.Service.Patch(c.Request().Context(), jsonObj)
		return AfterSavedWithLog(c, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
	} else {
		res, er3 := h.Service.Patch(c.Request().Context(), jsonObj)
		return AfterSavedWithLog(c, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
	}
}
func (h *Handler[T, K]) Delete(c echo.Context) error {
	id, ok, er1 := core.BuildId[K](c.Request(), h.ModelType, h.Keys, h.Indexes, h.IdMap)
	if er1 != nil {
		return c.String(http.StatusBadRequest, er1.Error())
	}
	if !ok {
		return c.String(http.StatusBadRequest, "Id type is not valid (Id type must be K)")
	}
	res, err := h.Service.Delete(c.Request().Context(), id)
	return AfterDeletedWithLog(c, res, err, h.LogError, h.WriteLog, h.Resource, h.Action.Delete)
}
