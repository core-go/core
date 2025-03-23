package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

const (
	Method = "method"
	Patch  = "patch"
)

type Load interface {
	Load(w http.ResponseWriter, r *http.Request)
}
type QueryHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}
type Transport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type ActionConfig struct {
	Search *string `yaml:"search" mapstructure:"search" json:"search,omitempty" gorm:"column:search" bson:"search,omitempty" dynamodbav:"search,omitempty" firestore:"search,omitempty"`
	Load   *string `yaml:"load" mapstructure:"load" json:"load,omitempty" gorm:"column:load" bson:"load,omitempty" dynamodbav:"load,omitempty" firestore:"load,omitempty"`
	Create string  `yaml:"create" mapstructure:"create" json:"create,omitempty" gorm:"column:create" bson:"create,omitempty" dynamodbav:"create,omitempty" firestore:"create,omitempty"`
	Update string  `yaml:"update" mapstructure:"update" json:"update,omitempty" gorm:"column:update" bson:"update,omitempty" dynamodbav:"update,omitempty" firestore:"update,omitempty"`
	Patch  string  `yaml:"patch" mapstructure:"patch" json:"patch,omitempty" gorm:"column:patch" bson:"patch,omitempty" dynamodbav:"patch,omitempty" firestore:"patch,omitempty"`
	Delete string  `yaml:"delete" mapstructure:"delete" json:"delete,omitempty" gorm:"column:delete" bson:"delete,omitempty" dynamodbav:"delete,omitempty" firestore:"delete,omitempty"`
}

func InitAction(conf *ActionConfig) ActionConfig {
	var c ActionConfig
	if conf != nil {
		c.Search = conf.Search
		c.Load = conf.Load
		c.Create = conf.Create
		c.Update = conf.Update
		c.Patch = conf.Patch
		c.Delete = conf.Delete
	}
	if c.Search == nil {
		x := "search"
		c.Search = &x
	}
	if c.Load == nil {
		x := "load"
		c.Load = &x
	}
	if len(c.Create) == 0 {
		c.Create = "create"
	}
	if len(c.Update) == 0 {
		c.Update = "update"
	}
	if len(c.Patch) == 0 {
		c.Patch = "patch"
	}
	if len(c.Delete) == 0 {
		c.Delete = "delete"
	}
	return c
}

type Attrs struct {
	Keys    []string
	Indexes map[string]int
	Error   func(context.Context, string, ...map[string]interface{})
}
type Attributes struct {
	Keys     []string
	Indexes  map[string]int
	Resource string
	Action   ActionConfig
	Error    func(context.Context, string, ...map[string]interface{})
	Log      func(context.Context, string, string, bool, string) error
}

func CreateAttributes(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...*ActionConfig) *Attributes {
	var action *ActionConfig
	if len(opts) > 0 && opts[0] != nil {
		action = opts[0]
	}
	a := InitAction(action)
	resource := BuildResourceName(modelType.Name())
	keys, indexes, _ := BuildMapField(modelType)
	return &Attributes{Keys: keys, Indexes: indexes, Resource: resource, Action: a, Error: logError, Log: writeLog}
}
func CreateAttrs(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{})) *Attrs {
	keys, indexes, _ := BuildMapField(modelType)
	return &Attrs{Keys: keys, Indexes: indexes, Error: logError}
}

type Parameters struct {
	Keys        []string
	Indexes     map[string]int
	Resource    string
	Action      ActionConfig
	Error       func(context.Context, string, ...map[string]interface{})
	Log         func(context.Context, string, string, bool, string) error
	ParamIndex  map[string]int
	FilterIndex int
	CSVIndex    map[string]int
}

func CreateParameters(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), action *ActionConfig, paramIndex map[string]int, filterIndex int, csvIndex map[string]int, opts ...func(context.Context, string, string, bool, string) error) *Parameters {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(opts) > 0 {
		writeLog = opts[0]
	}
	a := InitAction(action)
	resource := BuildResourceName(modelType.Name())
	keys, indexes, _ := BuildMapField(modelType)
	return &Parameters{Keys: keys, Indexes: indexes, Resource: resource, Action: a, Error: logError, Log: writeLog, ParamIndex: paramIndex, FilterIndex: filterIndex, CSVIndex: csvIndex}
}

func CheckId[T any](w http.ResponseWriter, r *http.Request, body *T, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) error {
	err := MatchId(r, body, keysJson, mapIndex)
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
func Decode[T any](w http.ResponseWriter, r *http.Request, opts ...func(context.Context, *T) error) (T, error) {
	var obj T
	er1 := json.NewDecoder(r.Body).Decode(&obj)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return obj, er1
	}
	if len(opts) > 0 && opts[0] != nil {
		er2 := opts[0](r.Context(), &obj)
		if er2 != nil {
			http.Error(w, er2.Error(), http.StatusInternalServerError)
		}
		return obj, er2
	}
	return obj, nil
}
func DecodeAndCheckId[T any](w http.ResponseWriter, r *http.Request, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) (T, error) {
	obj, er1 := Decode[T](w, r)
	if er1 != nil {
		return obj, er1
	}
	err := CheckId[T](w, r, &obj, keysJson, mapIndex, opts...)
	return obj, err
}
func BuildFieldMapAndCheckId[T any](w http.ResponseWriter, r *http.Request, keysJson []string, mapIndex map[string]int, ignorePatch bool, opts ...func(context.Context, *T) error) (*http.Request, T, map[string]interface{}, error) {
	var obj T
	if ignorePatch == false {
		r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	}
	body, er0 := BuildMapAndStruct(r, &obj, w)
	if er0 != nil {
		return r, obj, body, er0
	}
	er1 := CheckId[T](w, r, &obj, keysJson, mapIndex, opts...)
	return r, obj, body, er1
}
func BuildMapAndCheckId[T any](w http.ResponseWriter, r *http.Request, keysJson []string, mapIndex map[string]int, opts ...func(context.Context, *T) error) (*http.Request, T, map[string]interface{}, error) {
	r2, obj, body, er0 := BuildFieldMapAndCheckId[T](w, r, keysJson, mapIndex, false, opts...)
	if er0 != nil {
		return r2, obj, body, er0
	}
	json, er1 := BodyToJsonMap(r, &obj, body, keysJson, mapIndex)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	}
	return r2, obj, json, er1
}
func HasError(w http.ResponseWriter, r *http.Request, errors []ErrorMessage, err error, logError func(context.Context, string, ...map[string]interface{}), model interface{}, writeLog func(context.Context, string, string, bool, string) error, opts ...string) bool {
	var resource, action string
	if len(opts) > 0 && len(opts[0]) > 0 {
		resource = opts[0]
	}
	if len(opts) > 1 && len(opts[1]) > 0 {
		action = opts[1]
	}
	if err != nil {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, err.Error())
		}
		if logError != nil {
			if IsNil(model) {
				logError(r.Context(), err.Error())
			} else {
				logError(r.Context(), err.Error(), MakeMap(model))
			}
		}
		if logError == nil && writeLog == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}
		return true
	}
	if len(errors) > 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, fmt.Sprintf("Data Validation Failed %+v Error: %+v", model, err))
		}
		if logError != nil {
			logError(r.Context(), fmt.Sprintf("Data Validation Failed %+v Error: %+v", model, err))
		}
		JSON(w, http.StatusUnprocessableEntity, errors)
		return true
	}
	return false
}
func AfterDeletedWithLog(w http.ResponseWriter, r *http.Request, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	var resource, action string
	if len(opts) > 0 && len(opts[0]) > 0 {
		resource = opts[0]
	}
	if len(opts) > 1 && len(opts[1]) > 0 {
		action = opts[1]
	}
	if err != nil {
		if logError != nil {
			logError(r.Context(), "DELETE "+r.URL.Path+" with error "+err.Error())
		}
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "DELETE "+r.URL.Path+" with error "+err.Error())
		}
		if logError == nil && writeLog == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return err
	}
	if count > 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, true, "DELETE "+r.URL.Path)
		}
		return JSON(w, http.StatusOK, count)
	} else if count == 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "Data Not Found "+r.URL.Path)
		}
		return JSON(w, http.StatusNotFound, count)
	} else {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "DELETE "+r.URL.Path+" Conflict")
		}
		return JSON(w, http.StatusConflict, count)
	}
}
func AfterDeleted(w http.ResponseWriter, r *http.Request, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	if err != nil {
		if logError != nil {
			logError(r.Context(), "DELETE "+r.URL.Path+" with error "+err.Error())
		}
		if logError == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return err
	}
	if count > 0 {
		return JSON(w, http.StatusOK, count)
	} else if count == 0 {
		return JSON(w, http.StatusNotFound, count)
	} else {
		return JSON(w, http.StatusConflict, count)
	}
}
func AfterSavedWithLog(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	var resource, action string
	if len(opts) > 0 && len(opts[0]) > 0 {
		resource = opts[0]
	}
	if len(opts) > 1 && len(opts[1]) > 0 {
		action = opts[1]
	}
	if err != nil {
		if logError != nil {
			if IsNil(body) {
				logError(r.Context(), err.Error())
			} else {
				logError(r.Context(), err.Error(), MakeMap(body))
			}
		}
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, err.Error())
		}
		if logError == nil && writeLog == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return err
	}
	if count > 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, true, r.URL.Path)
		}
		if IsNil(body) {
			return JSON(w, http.StatusCreated, count)
		} else {
			return JSON(w, http.StatusCreated, body)
		}
	} else if count == 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "Data Not Found "+r.URL.Path)
		}
		return JSON(w, http.StatusNotFound, count)
	} else {
		if writeLog != nil {
			if IsNil(body) {
				writeLog(r.Context(), resource, action, false, "Conflict Data")
			} else {
				writeLog(r.Context(), resource, action, false, fmt.Sprintf("Conflict Data %+v", body))
			}
		}
		return JSON(w, http.StatusConflict, count)
	}
}
func AfterSaved(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	if err != nil {
		if logError != nil {
			if IsNil(body) {
				logError(r.Context(), err.Error())
			} else {
				logError(r.Context(), err.Error(), MakeMap(body))
			}
		}
		if logError == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return err
	}
	if count > 0 {
		if IsNil(body) {
			return JSON(w, http.StatusCreated, count)
		} else {
			return JSON(w, http.StatusCreated, body)
		}
	} else if count == 0 {
		return JSON(w, http.StatusNotFound, count)
	} else {
		return JSON(w, http.StatusConflict, count)
	}
}
func AfterCreatedWithLog(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	var resource, action string
	if len(opts) > 0 && len(opts[0]) > 0 {
		resource = opts[0]
	}
	if len(opts) > 1 && len(opts[1]) > 0 {
		action = opts[1]
	}
	if err != nil {
		if logError != nil {
			if IsNil(body) {
				logError(r.Context(), err.Error())
			} else {
				logError(r.Context(), err.Error(), MakeMap(body))
			}
		}
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, err.Error())
		}
		if logError == nil && writeLog == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return err
	}
	if count > 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, true, "")
		}
		if IsNil(body) {
			return JSON(w, http.StatusCreated, count)
		} else {
			return JSON(w, http.StatusCreated, body)
		}
	} else {
		if writeLog != nil {
			if IsNil(body) {
				writeLog(r.Context(), resource, action, false, "Duplicate Key")
			} else {
				writeLog(r.Context(), resource, action, false, fmt.Sprintf("Duplicate Key %+v", body))
			}
		}
		return JSON(w, http.StatusConflict, count)
	}
}
func AfterCreated(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	if err != nil {
		if logError != nil {
			if IsNil(body) {
				logError(r.Context(), err.Error())
			} else {
				logError(r.Context(), err.Error(), MakeMap(body))
			}
		}
		if logError == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return err
	}
	if count > 0 {
		if IsNil(body) {
			return JSON(w, http.StatusCreated, count)
		} else {
			return JSON(w, http.StatusCreated, body)
		}
	} else {
		return JSON(w, http.StatusConflict, count)
	}
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
	id, er1 := CreateId(r, modelType, idNames, indexes, opts...)
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
	Validate  func(context.Context, *T) ([]ErrorMessage, error)
	Keys      []string
	Indexes   map[string]int
	Resource  string
	ModelType reflect.Type
	Action    ActionConfig
	WriteLog  func(context.Context, string, string, bool, string) error
	IdMap     bool
	Builder   Builder[T]
}

func Newhandler[T any, K any](
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]ErrorMessage, error),
	opts ...Builder[T],
) *Handler[T, K] {
	return NewhandlerWithLog[T, K](service, logError, validate, nil, nil, opts...)
}
func NewhandlerWithLog[T any, K any](
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]ErrorMessage, error),
	writeLog func(context.Context, string, string, bool, string) error,
	action *ActionConfig,
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
	resource := BuildResourceName(modelType.Name())
	keys, indexes, _ := BuildMapField(modelType)
	a := InitAction(action)
	return &Handler[T, K]{Service: service, LogError: logError, Validate: validate, Keys: keys, Indexes: indexes, Resource: resource, ModelType: modelType, Builder: b, Action: a, WriteLog: writeLog, IdMap: idMap}
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
	ReturnWithLog(w, r, model, er2, h.LogError, h.WriteLog, h.Resource, action)
}
func (h *Handler[T, K]) Create(w http.ResponseWriter, r *http.Request) {
	var createFn func(context.Context, *T) error
	if h.Builder != nil {
		createFn = h.Builder.Create
	}
	model, er1 := Decode[T](w, r, createFn)
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !HasError(w, r, errors, er2, h.LogError, model, h.WriteLog, h.Resource, h.Action.Create) {
				res, er3 := h.Service.Create(r.Context(), &model)
				AfterCreatedWithLog(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
			}
		} else {
			res, er3 := h.Service.Create(r.Context(), &model)
			AfterCreatedWithLog(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Create)
		}
	}
}
func (h *Handler[T, K]) Update(w http.ResponseWriter, r *http.Request) {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	model, er1 := DecodeAndCheckId[T](w, r, h.Keys, h.Indexes, updateFn)
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !HasError(w, r, errors, er2, h.LogError, model, h.WriteLog, h.Resource, h.Action.Update) {
				res, er3 := h.Service.Update(r.Context(), &model)
				AfterSavedWithLog(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
			}
		} else {
			res, er3 := h.Service.Update(r.Context(), &model)
			AfterSavedWithLog(w, r, &model, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Update)
		}
	}
}
func (h *Handler[T, K]) Patch(w http.ResponseWriter, r *http.Request) {
	var updateFn func(context.Context, *T) error
	if h.Builder != nil {
		updateFn = h.Builder.Update
	}
	r, model, jsonObj, er1 := BuildMapAndCheckId[T](w, r, h.Keys, h.Indexes, updateFn)
	if er1 == nil {
		if h.Validate != nil {
			errors, er2 := h.Validate(r.Context(), &model)
			if !HasError(w, r, errors, er2, h.LogError, jsonObj, h.WriteLog, h.Resource, h.Action.Patch) {
				res, er3 := h.Service.Patch(r.Context(), jsonObj)
				AfterSavedWithLog(w, r, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
			}
		} else {
			res, er3 := h.Service.Patch(r.Context(), jsonObj)
			AfterSavedWithLog(w, r, jsonObj, res, er3, h.LogError, h.WriteLog, h.Resource, h.Action.Patch)
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
	AfterDeletedWithLog(w, r, res, err, h.LogError, h.WriteLog, h.Resource, h.Action.Delete)
}
