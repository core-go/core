package service

import (
	"context"
	"net/http"
	"reflect"
)

const (
	Method = "method"
	Update = "update"
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
type Search interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}
type Handler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
type HGenericService interface {
	Load(ctx context.Context, id interface{}) (interface{}, error)
	Insert(ctx context.Context, model interface{}) (int64, error)
	Update(ctx context.Context, model interface{}) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id interface{}) (int64, error)
}

func ReturnStatus(status int) ResultInfo {
	r := ResultInfo{Status: status}
	return r
}

func SetStatus(obj interface{}, status int) ResultInfo {
	r := ResultInfo{Status: status, Value: obj}
	return r
}

type WriterConfig struct {
	Status *StatusConfig `mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Action *ActionConfig `mapstructure:"action" json:"action,omitempty" gorm:"column:action" bson:"action,omitempty" dynamodbav:"action,omitempty" firestore:"action,omitempty"`
}

type StatusConfig struct {
	DuplicateKey    int  `mapstructure:"duplicate_key" json:"duplicateKey" gorm:"column:duplicatekey" bson:"duplicateKey" dynamodbav:"duplicateKey" firestore:"duplicateKey"`
	NotFound        int  `mapstructure:"not_found" json:"notFound" gorm:"column:notfound" bson:"notFound" dynamodbav:"notFound" firestore:"notFound"`
	Success         int  `mapstructure:"success" json:"success" gorm:"column:success" bson:"success" dynamodbav:"success" firestore:"success"`
	VersionError    int  `mapstructure:"version_error" json:"versionError" gorm:"column:versionerror" bson:"versionError" dynamodbav:"versionError" firestore:"versionError"`
	ValidationError *int `mapstructure:"validation_error" json:"validationError" gorm:"column:validationerror" bson:"validationError" dynamodbav:"validationError" firestore:"validationError"`
	Error           int  `mapstructure:"error" json:"error" gorm:"column:error" bson:"error" dynamodbav:"error" firestore:"error"`
}
type ActionConfig struct {
	Load   *string `mapstructure:"load" json:"load,omitempty" gorm:"column:load" bson:"load,omitempty" dynamodbav:"load,omitempty" firestore:"load,omitempty"`
	Create string  `mapstructure:"create" json:"create,omitempty" gorm:"column:create" bson:"create,omitempty" dynamodbav:"create,omitempty" firestore:"create,omitempty"`
	Update string  `mapstructure:"update" json:"update,omitempty" gorm:"column:update" bson:"update,omitempty" dynamodbav:"update,omitempty" firestore:"update,omitempty"`
	Patch  string  `mapstructure:"patch" json:"patch,omitempty" gorm:"column:patch" bson:"patch,omitempty" dynamodbav:"patch,omitempty" firestore:"patch,omitempty"`
	Delete string  `mapstructure:"delete" json:"delete,omitempty" gorm:"column:delete" bson:"delete,omitempty" dynamodbav:"delete,omitempty" firestore:"delete,omitempty"`
}
func InitializeStatus(status *StatusConfig) StatusConfig {
	var s StatusConfig
	if status != nil {
		s = *status
	} else {
		s.Error = 4
		k := s.Error
		s.DuplicateKey = 0
		s.NotFound = 0
		s.Success = 1
		s.VersionError = 2
		s.ValidationError = &k
	}
	if s.ValidationError == nil {
		k := s.Error
		s.ValidationError = &k
	}
	return s
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
type GenericHandler struct {
	*LoadHandler
	Status       StatusConfig
	Config       ActionConfig
	service      HGenericService
	modelBuilder ModelBuilder
	Validate     func(ctx context.Context, model interface{}) ([]ErrorMessage, error)
	Log          func(ctx context.Context, resource string, action string, success bool, desc string) error
	mapIndex     map[string]int
}

func NewHandler(genericService HGenericService, modelType reflect.Type, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	return NewHandlerWithConfig(genericService, modelType, nil, modelBuilder, logError, validate, options...)
}
func NewHandlerWithKeys(genericService HGenericService, keys []string, modelType reflect.Type, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	return NewHandlerWithKeysAndLog(genericService, keys, modelType, nil, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewHandlerWithConfig(genericService HGenericService, modelType reflect.Type, status *StatusConfig, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 && options[0] != nil {
		writeLog = options[0]
	}
	return NewHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewHandlerWithLog(genericService HGenericService, modelType reflect.Type, status *StatusConfig, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *ActionConfig) *GenericHandler {
	return NewHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, resource, conf)
}
func NewHandlerWithKeysAndLog(genericService HGenericService, keys []string, modelType reflect.Type, status *StatusConfig, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *ActionConfig) *GenericHandler {
	if keys == nil || len(keys) == 0 {
		keys = GetJsonPrimaryKeys(modelType)
	}
	if len(resource) == 0 {
		resource = BuildResourceName(modelType.Name())
	}
	var writeLog2 func(context.Context, string, string, bool, string) error
	if conf != nil && conf.Load != nil {
		writeLog2 = writeLog
	}
	c := InitializeAction(conf)
	s := InitializeStatus(status)
	loadHandler := NewLoadHandlerWithKeysAndLog(genericService.Load, keys, modelType, logError, writeLog2, *c.Load, resource)
	_, jsonMapIndex := BuildMapField(modelType)

	return &GenericHandler{LoadHandler: loadHandler, service: genericService, Status: s, modelBuilder: modelBuilder, Validate: validate, mapIndex: jsonMapIndex, Log: writeLog, Config: c}
}
func (h *GenericHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.Insert(w, r)
}
func (h *GenericHandler) Insert(w http.ResponseWriter, r *http.Request) {
	body, er0 := NewModel(h.ModelType, r.Body)
	if er0 != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	if h.modelBuilder != nil {
		body, er0 = h.modelBuilder.BuildToInsert(r.Context(), body)
		if er0 != nil {
			Respond(w, r, http.StatusInternalServerError, InternalServerError, er0, h.Error, h.Log, h.Resource, h.Config.Create)
		}
	}
	if h.Validate != nil {
		errors, er1 := h.Validate(r.Context(), body)
		if HasError(w, r, errors, er1, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Config.Create) {
			return
		}
	}
	var count int64
	var er2 error
	count, er2 = h.service.Insert(r.Context(), body)
	if count <= 0 && er2 == nil {
		if h.modelBuilder == nil {
			RespondAndLog(w, r, http.StatusConflict, ReturnStatus(h.Status.DuplicateKey), h.Log, false, h.Resource, h.Config.Create, "Duplicate Key")
			return
		}
		i := 0
		for count <= 0 && i <= 5 {
			i++
			body, er2 = h.modelBuilder.BuildToInsert(r.Context(), body)
			if er2 != nil {
				Respond(w, r, http.StatusInternalServerError, InternalServerError, er2, h.Error, h.Log, h.Resource, h.Config.Create)
				return
			}
			count, er2 = h.service.Insert(r.Context(), body)
			if er2 != nil {
				Respond(w, r, http.StatusInternalServerError, InternalServerError, er2, h.Error, h.Log, h.Resource, h.Config.Create)
				return
			}
			if count > 0 {
				Succeed(w, r, http.StatusCreated, SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Create)
				return
			}
			if i == 5 {
				RespondAndLog(w, r, http.StatusConflict, ReturnStatus(h.Status.DuplicateKey), h.Log, false, h.Resource, h.Config.Create, "Duplicate Key")
				return
			}
		}
	} else if er2 != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, er2, h.Error, h.Log, h.Resource, h.Config.Create)
		return
	}
	Succeed(w, r, http.StatusCreated, SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Create)
}
func (h *GenericHandler) Update(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), Method, Update))
	body, er0 := NewModel(h.ModelType, r.Body)
	if er0 != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}
	er1 := CheckId(w, r, body, h.Keys, h.mapIndex)
	if er1 != nil {
		return
	}
	if h.modelBuilder != nil {
		body, er0 = h.modelBuilder.BuildToUpdate(r.Context(), body)
		if er0 != nil {
			Respond(w, r, http.StatusInternalServerError, InternalServerError, er0, h.Error, h.Log, h.Resource, h.Config.Update)
			return
		}
	}
	if h.Validate != nil {
		errors, er2 := h.Validate(r.Context(), body)
		if HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Config.Update) {
			return
		}
	}
	count, er3 := h.service.Update(r.Context(), body)
	HandleResult(w, r, body, count, er3, h.Status, h.Error, h.Log, h.Resource, h.Config.Update)
}
func (h *GenericHandler) Patch(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	bodyStruct := reflect.New(h.ModelType).Interface()
	body0, er0 := BuildMapAndStruct(r, bodyStruct)
	if er0 != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}
	er1 := CheckId(w, r, bodyStruct, h.Keys, h.mapIndex)
	if er1 != nil {
		return
	}
	body, er2 := BodyToJson(w, r, bodyStruct, body0, h.Keys, h.mapIndex, h.modelBuilder.BuildToPatch, h.Error, h.Log, h.Resource, h.Config.Patch)
	if er2 != nil {
		return
	}
	if h.Validate != nil {
		errors, er3 := h.Validate(r.Context(), bodyStruct)
		if HasError(w, r, errors, er3, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Config.Patch) {
			return
		}
	}
	count, er4 := h.service.Patch(r.Context(), body)
	HandleResult(w, r, body, count, er4, h.Status, h.Error, h.Log, h.Resource, h.Config.Patch)
}
func (h *GenericHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, er1 := BuildId(r, h.ModelType, h.Keys, h.Indexes)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	count, er2 := h.service.Delete(r.Context(), id)
	HandleDelete(w, r, count, er2, h.Error, h.Log, h.Resource, h.Config.Delete)
}

func HandleDelete(w http.ResponseWriter, r *http.Request, count int64, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, err, logError, writeLog, resource, action)
		return
	}
	if count > 0 {
		Succeed(w, r, http.StatusOK, count, writeLog, resource, action)
	} else if count == 0 {
		RespondAndLog(w, r, http.StatusNotFound, count, writeLog, false, resource, action, "Data Not Found")
	} else {
		RespondAndLog(w, r, http.StatusConflict, count, writeLog, false, resource, action, "Conflict")
	}
}
func BodyToJson(w http.ResponseWriter, r *http.Request, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int, buildToPatch func(context.Context, interface{}) (interface{}, error), logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options...string) (map[string]interface{}, error) {
	body, err := BodyToJsonMap(r, structBody, body, jsonIds, mapIndex, buildToPatch)
	if err != nil {
		// http.Error(w, "Invalid Data: "+err.Error(), http.StatusBadRequest)
		Respond(w, r, http.StatusInternalServerError, InternalServerError, err, logError, writeLog, options...)
	}
	return body, err
}
func HasError(w http.ResponseWriter, r *http.Request, errors []ErrorMessage, err error, status int, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options...string) bool {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, err, logError, writeLog, resource, action)
		return true
	}
	if len(errors) > 0 {
		result0 := ResultInfo{Status: status, Errors: errors}
		RespondAndLog(w, r, http.StatusUnprocessableEntity, result0, writeLog, false, resource, action, "Data Validation Failed")
		return true
	}
	return false
}
func HandleResult(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, status StatusConfig, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, err, logError, writeLog, resource, action)
		return
	}
	if count == -1 {
		RespondAndLog(w, r, http.StatusConflict, ReturnStatus(status.VersionError), writeLog, false, resource, action, "Data Version Error")
	} else if count == 0 {
		RespondAndLog(w, r, http.StatusNotFound, ReturnStatus(status.NotFound), writeLog, false, resource, action, "Data Not Found")
	} else {
		Succeed(w, r, http.StatusOK, SetStatus(body, status.Success), writeLog, resource, action)
	}
}
func AfterCreated(w http.ResponseWriter, r *http.Request, body interface{}, count int64, err error, status StatusConfig, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, err, logError, writeLog, resource, action)
		return
	}
	if count <= 0 {
		RespondAndLog(w, r, http.StatusConflict, ReturnStatus(status.DuplicateKey), writeLog, false, resource, action, "Duplicate Key")
	} else {
		Succeed(w, r, http.StatusCreated, SetStatus(body, status.Success), writeLog, resource, action)
	}
}
