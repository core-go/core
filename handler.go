package core

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

const (
	Method = "method"
	Patch  = "patch"
)

type Diff interface {
	Diff(w http.ResponseWriter, r *http.Request)
}
type Approve interface {
	Approve(w http.ResponseWriter, r *http.Request)
	Reject(w http.ResponseWriter, r *http.Request)
}
type Load interface {
	Load(w http.ResponseWriter, r *http.Request)
}
type SearchHandlr interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}
type Query interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}
type QueryHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}
type QueryController interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}
type Handler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
type HttpPort interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
type Controller interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
type ActionConfig struct {
	Load   *string `yaml:"load" mapstructure:"load" json:"load,omitempty" gorm:"column:load" bson:"load,omitempty" dynamodbav:"load,omitempty" firestore:"load,omitempty"`
	Create string  `yaml:"create" mapstructure:"create" json:"create,omitempty" gorm:"column:create" bson:"create,omitempty" dynamodbav:"create,omitempty" firestore:"create,omitempty"`
	Update string  `yaml:"update" mapstructure:"update" json:"update,omitempty" gorm:"column:update" bson:"update,omitempty" dynamodbav:"update,omitempty" firestore:"update,omitempty"`
	Patch  string  `yaml:"patch" mapstructure:"patch" json:"patch,omitempty" gorm:"column:patch" bson:"patch,omitempty" dynamodbav:"patch,omitempty" firestore:"patch,omitempty"`
	Delete string  `yaml:"delete" mapstructure:"delete" json:"delete,omitempty" gorm:"column:delete" bson:"delete,omitempty" dynamodbav:"delete,omitempty" firestore:"delete,omitempty"`
}
type ActionConf struct {
	Search *string `yaml:"search" mapstructure:"search" json:"search,omitempty" gorm:"column:search" bson:"search,omitempty" dynamodbav:"search,omitempty" firestore:"search,omitempty"`
	Load   *string `yaml:"load" mapstructure:"load" json:"load,omitempty" gorm:"column:load" bson:"load,omitempty" dynamodbav:"load,omitempty" firestore:"load,omitempty"`
	Create string  `yaml:"create" mapstructure:"create" json:"create,omitempty" gorm:"column:create" bson:"create,omitempty" dynamodbav:"create,omitempty" firestore:"create,omitempty"`
	Update string  `yaml:"update" mapstructure:"update" json:"update,omitempty" gorm:"column:update" bson:"update,omitempty" dynamodbav:"update,omitempty" firestore:"update,omitempty"`
	Patch  string  `yaml:"patch" mapstructure:"patch" json:"patch,omitempty" gorm:"column:patch" bson:"patch,omitempty" dynamodbav:"patch,omitempty" firestore:"patch,omitempty"`
	Delete string  `yaml:"delete" mapstructure:"delete" json:"delete,omitempty" gorm:"column:delete" bson:"delete,omitempty" dynamodbav:"delete,omitempty" firestore:"delete,omitempty"`
}

func InitAction(conf *ActionConf) ActionConf {
	var c ActionConf
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
func InitializeAction(conf *ActionConfig) ActionConfig {
	var c ActionConfig
	if conf != nil {
		c.Load = conf.Load
		c.Create = conf.Create
		c.Update = conf.Update
		c.Patch = conf.Patch
		c.Delete = conf.Delete
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

type Parameters struct {
	Keys        []string
	Indexes     map[string]int
	ModelType   reflect.Type
	Resource    string
	Action      ActionConf
	Error       func(context.Context, string, ...map[string]interface{})
	Log         func(context.Context, string, string, bool, string) error
	ParamIndex  map[string]int
	FilterIndex int
	CSVIndex    map[string]int
}

func CreateParameters(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), action *ActionConf, paramIndex map[string]int, filterIndex int, csvIndex map[string]int, options ...func(context.Context, string, string, bool, string) error) *Parameters {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	a := InitAction(action)
	resource := BuildResourceName(modelType.Name())
	keys, indexes, _ := BuildMapField(modelType)
	return &Parameters{Keys: keys, Indexes: indexes, ModelType: modelType, Resource: resource, Action: a, Error: logError, Log: writeLog, ParamIndex: paramIndex, FilterIndex: filterIndex, CSVIndex: csvIndex}
}

type Params struct {
	Keys      []string
	Indexes   map[string]int
	ModelType reflect.Type
	Resource  string
	Action    ActionConfig
	Error     func(context.Context, string, ...map[string]interface{})
	Log       func(context.Context, string, string, bool, string) error
}

func MakeParams(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), action *ActionConfig, options ...func(context.Context, string, string, bool, string) error) *Params {
	return CreateParams(modelType, logError, action, options...)
}
func InitParams(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{})) *Params {
	return CreateParams(modelType, logError, nil)
}
func CreateParams(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), action *ActionConfig, options ...func(context.Context, string, string, bool, string) error) *Params {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	a := InitializeAction(action)
	resource := BuildResourceName(modelType.Name())
	keys, indexes, _ := BuildMapField(modelType)
	return &Params{Keys: keys, Indexes: indexes, ModelType: modelType, Resource: resource, Action: a, Error: logError, Log: writeLog}
}

func CheckId(w http.ResponseWriter, r *http.Request, body interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) error) error {
	err := MatchId(r, body, keysJson, mapIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if len(options) > 0 && options[0] != nil {
		er2 := options[0](r.Context(), body)
		if er2 != nil {
			http.Error(w, er2.Error(), http.StatusInternalServerError)
		}
		return er2
	}
	return nil
}
func DecodeAndCheckId(w http.ResponseWriter, r *http.Request, obj interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) error) error {
	er1 := Decode(w, r, obj)
	if er1 != nil {
		return er1
	}
	return CheckId(w, r, obj, keysJson, mapIndex, options...)
}
func BuildFieldMapAndCheckId(w http.ResponseWriter, r *http.Request, obj interface{}, keysJson []string, mapIndex map[string]int, ignorePatch bool) (*http.Request, map[string]interface{}, error) {
	if ignorePatch == false {
		r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	}
	body, er0 := BuildMapAndStruct(r, obj, w)
	if er0 != nil {
		return r, body, er0
	}
	er1 := CheckId(w, r, obj, keysJson, mapIndex)
	return r, body, er1
}
func BuildMapAndCheckId(w http.ResponseWriter, r *http.Request, obj interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) error) (*http.Request, map[string]interface{}, error) {
	r2, body, er0 := BuildFieldMapAndCheckId(w, r, obj, keysJson, mapIndex, false)
	if er0 != nil {
		return r2, body, er0
	}
	json, er1 := BodyToJsonMap(r, obj, body, keysJson, mapIndex, options...)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	}
	return r2, json, er1
}
func HasError(w http.ResponseWriter, r *http.Request, errors []ErrorMessage, err error, logError func(context.Context, string, ...map[string]interface{}), model interface{}, writeLog func(context.Context, string, string, bool, string) error, options ...string) bool {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
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
		writeLog(r.Context(), resource, action, false, fmt.Sprintf("Data Validation Failed %+v Error: %+v", model, err))
		JSON(w, http.StatusUnprocessableEntity, errors)
		return true
	}
	return false
}
func AfterDeletedWithLog(w http.ResponseWriter, r *http.Request, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
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
		return
	}
	if count > 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, true, "DELETE "+r.URL.Path)
		}
		JSON(w, http.StatusOK, count)
	} else if count == 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "Data Not Found "+r.URL.Path)
		}
		JSON(w, http.StatusNotFound, count)
	} else {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "DELETE "+r.URL.Path+" Conflict")
		}
		JSON(w, http.StatusConflict, count)
	}
}
func AfterDeleted(w http.ResponseWriter, r *http.Request, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) {
	if err != nil {
		if logError != nil {
			logError(r.Context(), "DELETE "+r.URL.Path+" with error "+err.Error())
		}
		if logError == nil {
			JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			JSON(w, http.StatusInternalServerError, InternalServerError)
		}
		return
	}
	if count > 0 {
		JSON(w, http.StatusOK, count)
	} else if count == 0 {
		JSON(w, http.StatusNotFound, count)
	} else {
		JSON(w, http.StatusConflict, count)
	}
}
func AfterSavedWithLog(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
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
		return
	}
	if count > 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, true, r.URL.Path)
		}
		if IsNil(body) {
			JSON(w, http.StatusCreated, count)
		} else {
			JSON(w, http.StatusCreated, body)
		}
	} else if count == 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, false, "Data Not Found "+r.URL.Path)
		}
		JSON(w, http.StatusNotFound, count)
	} else {
		if writeLog != nil {
			if IsNil(body) {
				writeLog(r.Context(), resource, action, false, "Conflict Data")
			} else {
				writeLog(r.Context(), resource, action, false, fmt.Sprintf("Conflict Data %+v", body))
			}
		}
		JSON(w, http.StatusConflict, count)
	}
}
func AfterSaved(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) {
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
		return
	}
	if count > 0 {
		if IsNil(body) {
			JSON(w, http.StatusCreated, count)
		} else {
			JSON(w, http.StatusCreated, body)
		}
	} else if count == 0 {
		JSON(w, http.StatusNotFound, count)
	} else {
		JSON(w, http.StatusConflict, count)
	}
}
func AfterCreatedWithLog(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
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
		return
	}
	if count > 0 {
		if writeLog != nil {
			writeLog(r.Context(), resource, action, true, "")
		}
		if IsNil(body) {
			JSON(w, http.StatusCreated, count)
		} else {
			JSON(w, http.StatusCreated, body)
		}
	} else {
		if writeLog != nil {
			if IsNil(body) {
				writeLog(r.Context(), resource, action, false, "Duplicate Key")
			} else {
				writeLog(r.Context(), resource, action, false, fmt.Sprintf("Duplicate Key %+v", body))
			}
		}
		JSON(w, http.StatusConflict, count)
	}
}
func AfterCreated(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) {
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
		return
	}
	if count > 0 {
		if IsNil(body) {
			JSON(w, http.StatusCreated, count)
		} else {
			JSON(w, http.StatusCreated, body)
		}
	} else {
		JSON(w, http.StatusConflict, count)
	}
}
