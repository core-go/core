package gin

import (
	"context"
	"github.com/core-go/core"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func CheckId[T any](ctx *gin.Context, body *T, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) error {
	return core.CheckId[T](ctx.Writer, ctx.Request, body, keysJson, mapIndex, opts...)
}
func AfterDeletedWithLog(ctx *gin.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterDeletedWithLog(ctx.Writer, ctx.Request, count, err, logError, writeLog, opts...)
}
func AfterDeleted(ctx *gin.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterDeleted(ctx.Writer, ctx.Request, count, err, logError)
}
func HasError(ctx *gin.Context, errors []core.ErrorMessage, err error, logError func(context.Context, string, ...map[string]interface{}), model interface{}, writeLog func(context.Context, string, string, bool, string) error, opts ...string) bool {
	return core.HasError(ctx.Writer, ctx.Request, errors, err, logError, model, writeLog, opts...)
}
func AfterSavedWithLog(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterSavedWithLog(ctx.Writer, ctx.Request, body, count, err, logError, writeLog, opts...)
}
func AfterSaved(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterSaved(ctx.Writer, ctx.Request, body, count, err, logError)
}
func AfterCreatedWithLog(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterCreatedWithLog(ctx.Writer, ctx.Request, body, count, err, logError, writeLog, opts...)
}
func AfterCreated(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterCreated(ctx.Writer, ctx.Request, body, count, err, logError)
}
func BuildFieldMapAndCheckId[T any](ctx *gin.Context, keysJson []string, mapIndex map[string]int, ignorePatch bool, opts ...func(context.Context, *T) error) (T, map[string]interface{}, error) {
	var obj T
	if ignorePatch == false {
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), core.Method, core.Patch))
	}
	body, er0 := core.BuildMapAndStruct(ctx.Request, &obj)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, er0.Error())
		return obj, body, er0
	}
	er1 := CheckId[T](ctx, &obj, keysJson, mapIndex, opts...)
	return obj, body, er1
}
func BuildMapAndCheckId[T any](ctx *gin.Context, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) (T, map[string]interface{}, error) {
	obj, body, er0 := BuildFieldMapAndCheckId[T](ctx, keysJson, mapIndex, false, opts...)
	if er0 != nil {
		return obj, body, er0
	}
	jsonObj, er1 := core.BodyToJsonMap(ctx.Request, &obj, body, keysJson, mapIndex)
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

func Decode[T any](c *gin.Context, opts ...func(context.Context, *T) error) (T, error) {
	return core.Decode[T](c.Writer, c.Request, opts...)
}
func DecodeAndCheckId[T any](c *gin.Context, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) (T, error) {
	return core.DecodeAndCheckId[T](c.Writer, c.Request, keysJson, mapIndex, opts...)
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
	writeLog func(context.Context, string, string, bool, string) error,
	action *core.ActionConfig,
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

func (h *Handler[T, K]) Load(c *gin.Context) {
	r := c.Request
	id, ok, er1 := core.BuildId[K](r, h.ModelType, h.Keys, h.Indexes, h.IdMap)
	if er1 != nil {
		c.String(http.StatusBadRequest, er1.Error())
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "Id type is not valid (Id type must be K)")
		return
	}
	model, er2 := h.Service.Load(r.Context(), id)
	action := ""
	if h.Action.Load != nil {
		action = *h.Action.Load
	}
	ReturnWithLog(c, model, er2, h.LogError, h.WriteLog, h.Resource, action)
}
func (h *Handler[T, K]) Create(c *gin.Context) {
	var createFn func(context.Context, *T) error
	if h.Builder != nil {
		createFn = h.Builder.Create
	}
	model, er1 := Decode[T](c, createFn)
	r := c.Request
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !HasError(c, errors, er2, h.LogError, model, h.WriteLog, h.Resource, h.Action.Create) {
				res, er3 := h.Service.Create(r.Context(), &model)
				AfterCreatedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
			}
		} else {
			res, er3 := h.Service.Create(r.Context(), &model)
			AfterCreatedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
		}
	}
}
func (h *Handler[T, K]) Update(c *gin.Context) {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	model, er1 := DecodeAndCheckId[T](c, h.Keys, h.Indexes, updateFn)
	if er1 == nil {
		r := c.Request
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !HasError(c, errors, er2, h.LogError, model, h.WriteLog, h.Resource, h.Action.Update) {
				res, er3 := h.Service.Update(r.Context(), &model)
				AfterSavedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
			}
		} else {
			res, er3 := h.Service.Update(r.Context(), &model)
			AfterSavedWithLog(c, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
		}
	}
}
func (h *Handler[T, K]) Patch(c *gin.Context) {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	model, jsonObj, er1 := BuildMapAndCheckId[T](c, h.Keys, h.Indexes, updateFn)
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(c.Request.Context(), &model)
			if !HasError(c, errors, er2, h.LogError, jsonObj, h.WriteLog, h.Resource, h.Action.Patch) {
				res, er3 := h.Service.Patch(c.Request.Context(), jsonObj)
				AfterSavedWithLog(c, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
			}
		} else {
			res, er3 := h.Service.Patch(c.Request.Context(), jsonObj)
			AfterSavedWithLog(c, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
		}
	}
}
func (h *Handler[T, K]) Delete(c *gin.Context) {
	id, ok, er1 := core.BuildId[K](c.Request, h.ModelType, h.Keys, h.Indexes, h.IdMap)
	if er1 != nil {
		c.String(http.StatusBadRequest, er1.Error())
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "Id type is not valid (Id type must be K)")
		return
	}
	res, err := h.Service.Delete(c.Request.Context(), id)
	AfterDeletedWithLog(c, res, err, h.LogError, h.WriteLog, h.Resource, h.Action.Delete)
}
