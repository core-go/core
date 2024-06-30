package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/core-go/core"
)

type Service[T any, K any] interface {
	Load(ctx context.Context, id K) (*T, error)
	Create(ctx context.Context, model *T) (int64, error)
	Update(ctx context.Context, model *T) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id K) (int64, error)
}

type Handler[T any, K any] struct {
	Service   Service[T, K]
	LogError  func(context.Context, string, ...map[string]interface{})
	Validate  func(context.Context, *T) ([]core.ErrorMessage, error)
	Keys      []string
	Indexes   map[string]int
	Resource  string
	ModelType reflect.Type
	Action    core.ActionConf
	WriteLog  func(context.Context, string, string, bool, string) error
	IdMap     bool
}

func Newhandler[T any, K any](
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
) *Handler[T, K] {
	return NewhandlerWithLog[T, K](service, logError, validate, nil, nil)
}
func NewhandlerWithLog[T any, K any](
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	action *core.ActionConf,
	writeLog func(context.Context, string, string, bool, string) error,
) *Handler[T, K] {
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
	return &Handler[T, K]{Service: service, LogError: logError, Validate: validate, Keys: keys, Indexes: indexes, Resource: resource, ModelType: modelType, Action: a, WriteLog: writeLog, IdMap: idMap}
}

func mapToStruct(obj interface{}, des interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &des)
}

func BuildId[K any](r *http.Request, modelType reflect.Type, idNames []string, indexes map[string]int, isMap bool, opts ...int) (K, bool, error) {
	var k K
	id, er1 := core.BuildId(r, modelType, idNames, indexes, opts...)
	if er1 != nil {
		return k, false, er1
	}
	var ok bool
	if len(idNames) <= 1 {
		k, ok = id.(K)
		return k, ok, nil
	} else {
		if isMap {
			k, ok = id.(K)
			return k, ok, nil
		} else {
			err := mapToStruct(id, &k)
			if err != nil {
				return k, true, err
			}
			return k, true, nil
		}
	}
}
func (h *Handler[T, K]) Load(w http.ResponseWriter, r *http.Request) {
	id, ok, er1 := BuildId[K](r, h.ModelType, h.Keys, h.Indexes, h.IdMap)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	if !ok {
		http.Error(w, "Id type is not valid (Id type must be K)", http.StatusBadRequest)
		return
	}
	model, er2 := h.Service.Load(r.Context(), id)
	action := ""
	if h.Action.Load != nil {
		action = *h.Action.Load
	}
	core.Return(w, r, model, er2, h.LogError, h.WriteLog, h.Resource, action)
}
func (h *Handler[T, K]) Create(w http.ResponseWriter, r *http.Request) {
	var model T
	er1 := core.Decode(w, r, &model)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &model)
		if !core.HasError(w, r, errors, er2, h.LogError, h.WriteLog, h.Resource, h.Action.Create) {
			res, er3 := h.Service.Create(r.Context(), &model)
			core.AfterCreated(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
		}
	}
}
func (h *Handler[T, K]) Update(w http.ResponseWriter, r *http.Request) {
	var model T
	er1 := core.DecodeAndCheckId(w, r, &model, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &model)
		if !core.HasError(w, r, errors, er2, h.LogError, h.WriteLog, h.Resource, h.Action.Update) {
			res, er3 := h.Service.Update(r.Context(), &model)
			core.HandleResult(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
		}
	}
}
func (h *Handler[T, K]) Patch(w http.ResponseWriter, r *http.Request) {
	var model T
	r, jsonObj, er1 := core.BuildMapAndCheckId(w, r, &model, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &model)
		if !core.HasError(w, r, errors, er2, h.LogError, h.WriteLog, h.Resource, h.Action.Patch) {
			res, er3 := h.Service.Patch(r.Context(), jsonObj)
			core.HandleResult(w, r, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
		}
	}
}
func (h *Handler[T, K]) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok, er1 := BuildId[K](r, h.ModelType, h.Keys, h.Indexes, h.IdMap)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	if !ok {
		http.Error(w, "Id type is not valid (Id type must be K)", http.StatusBadRequest)
		return
	}
	res, err := h.Service.Delete(r.Context(), id)
	core.HandleDelete(w, r, res, err, h.LogError, h.WriteLog, h.Resource, h.Action.Delete)
}
