package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/core-go/core"
)

const (
	Method = "method"
	Patch  = "patch"
)

type Validate[T any] func(ctx context.Context, model T) ([]core.ErrorMessage, error)

type Builder[T any] interface {
	Create(context.Context, *T) error
	Update(context.Context, *T) error
}

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
	Builder   Builder[T]
}

func Newhandler[T any, K any](
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	opts ...Builder[T],
) *Handler[T, K] {
	return NewhandlerWithLog[T, K](service, logError, validate, nil, nil, opts...)
}
func NewhandlerWithLog[T any, K any](
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	action *core.ActionConf,
	writeLog func(context.Context, string, string, bool, string) error,
	opts ...Builder[T],
) *Handler[T, K] {
	var b Builder[T]
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
	var createFn func(context.Context, *T) error
	if h.Builder != nil {
		createFn = h.Builder.Create
	}
	var model T
	er1 := Decode[T](w, r, &model, createFn)
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !core.HasError(w, r, errors, er2, h.LogError, h.WriteLog, h.Resource, h.Action.Create) {
				res, er3 := h.Service.Create(r.Context(), &model)
				core.AfterCreated(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
			}
		} else {
			res, er3 := h.Service.Create(r.Context(), &model)
			core.AfterCreated(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
		}
	}
}
func (h *Handler[T, K]) Update(w http.ResponseWriter, r *http.Request) {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	model, er1 := DecodeAndCheckId(w, r, h.Keys, h.Indexes, updateFn)
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !core.HasError(w, r, errors, er2, h.LogError, h.WriteLog, h.Resource, h.Action.Update) {
				res, er3 := h.Service.Update(r.Context(), &model)
				core.HandleResult(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
			}
		} else {
			res, er3 := h.Service.Update(r.Context(), &model)
			core.HandleResult(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
		}
	}
}
func (h *Handler[T, K]) Patch(w http.ResponseWriter, r *http.Request) {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	r, model, jsonObj, er1 := BuildMapAndCheckId(w, r, h.Keys, h.Indexes, updateFn)
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !core.HasError(w, r, errors, er2, h.LogError, h.WriteLog, h.Resource, h.Action.Patch) {
				res, er3 := h.Service.Patch(r.Context(), jsonObj)
				core.HandleResult(w, r, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
			}
		} else {
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

func Decode[T any](w http.ResponseWriter, r *http.Request, obj *T, opts ...func(context.Context, *T) error) error {
	er1 := json.NewDecoder(r.Body).Decode(obj)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return er1
	}
	if len(opts) > 0 && opts[0] != nil {
		er2 := opts[0](r.Context(), obj)
		if er2 != nil {
			http.Error(w, er2.Error(), http.StatusInternalServerError)
		}
		return er2
	}
	return nil
}
func DecodeAndCheckId[T any](w http.ResponseWriter, r *http.Request, keysJson []string, mapIndex map[string]int, options ...func(context.Context, *T) error) (T, error) {
	var obj T
	er1 := core.Decode(w, r, &obj)
	if er1 != nil {
		return obj, er1
	}
	err := CheckId[T](w, r, &obj, keysJson, mapIndex, options...)
	return obj, err
}
func CheckId[T any](w http.ResponseWriter, r *http.Request, body *T, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) error {
	err := core.MatchId(r, body, keysJson, mapIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if len(opts) > 0 && opts[0] != nil {
		er2 := opts[0](r.Context(), body)
		if er2 != nil {
			http.Error(w, er2.Error(), http.StatusInternalServerError)
		}
		return er2
	}
	return nil
}
func BuildMapAndCheckId[T any](w http.ResponseWriter, r *http.Request, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) (*http.Request, T, map[string]interface{}, error) {
	r2, obj, body, er0 := BuildFieldMapAndCheckId[T](w, r, keysJson, mapIndex, false, opts...)
	if er0 != nil {
		return r2, obj, body, er0
	}
	json, er1 := core.BodyToJsonMap(r, &obj, body, keysJson, mapIndex)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	}
	return r2, obj, json, er1
}

func BuildFieldMapAndCheckId[T any](w http.ResponseWriter, r *http.Request, keysJson []string, mapIndex map[string]int, ignorePatch bool, opts ...func(context.Context, *T) error) (*http.Request, T, map[string]interface{}, error) {
	var obj T
	if ignorePatch == false {
		r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	}
	body, er0 := core.BuildMapAndStruct(r, &obj, w)
	if er0 != nil {
		return r, obj, body, er0
	}
	er1 := CheckId[T](w, r, &obj, keysJson, mapIndex, opts...)
	return r, obj, body, er1
}
