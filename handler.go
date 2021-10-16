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

type Load interface {
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

func NewGenericHandler(genericService HGenericService, modelType reflect.Type, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	return NewGenericHandlerWithConfig(genericService, modelType, nil, modelBuilder, logError, validate, options...)
}
func NewGenericHandlerWithKeys(genericService HGenericService, keys []string, modelType reflect.Type, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	return NewGenericHandlerWithKeysAndLog(genericService, keys, modelType, nil, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewGenericHandlerWithConfig(genericService HGenericService, modelType reflect.Type, status *StatusConfig, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 && options[0] != nil {
		writeLog = options[0]
	}
	return NewGenericHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewGenericHandlerWithLog(genericService HGenericService, modelType reflect.Type, status *StatusConfig, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *ActionConfig) *GenericHandler {
	return NewGenericHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, resource, conf)
}
func NewGenericHandlerWithKeysAndLog(genericService HGenericService, keys []string, modelType reflect.Type, status *StatusConfig, modelBuilder ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *ActionConfig) *GenericHandler {
	if keys == nil || len(keys) == 0 {
		keys = GetJsonPrimaryKeys(modelType)
	}
	if len(resource) == 0 {
		resource = BuildResourceName(modelType.Name())
	}
	var c ActionConfig
	var writeLog2 func(context.Context, string, string, bool, string) error
	if conf != nil {
		if conf.Load != nil {
			writeLog2 = writeLog
		}
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
			ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Create, er0, h.Log)
		}
	}
	if h.Validate != nil {
		errors, er1 := h.Validate(r.Context(), body)
		if er1 != nil {
			ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Create, er1, h.Log)
			return
		}
		if len(errors) > 0 {
			//result0 := model.ResultInfo{Status: model.StatusError, Errors: MakeErrors(errors)}
			result0 := ResultInfo{Status: *h.Status.ValidationError, Errors: errors}
			RespondAndLog(w, r, http.StatusUnprocessableEntity, result0, h.Log, h.Resource, h.Config.Create, false, "Data Validation Failed")
			return
		}
	}
	var count int64
	var er2 error
	count, er2 = h.service.Insert(r.Context(), body)
	if count <= 0 && er2 == nil {
		if h.modelBuilder == nil {
			RespondAndLog(w, r, http.StatusConflict, ReturnStatus(h.Status.DuplicateKey), h.Log, h.Resource, h.Config.Create, false, "Duplicate Key")
			return
		}
		i := 0
		for count <= 0 && i <= 5 {
			i++
			body, er2 = h.modelBuilder.BuildToInsert(r.Context(), body)
			if er2 != nil {
				ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Create, er2, h.Log)
				return
			}
			count, er2 = h.service.Insert(r.Context(), body)
			if er2 != nil {
				ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Create, er2, h.Log)
				return
			}
			if count > 0 {
				Succeed(w, r, http.StatusCreated, SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Create)
				return
			}
			if i == 5 {
				RespondAndLog(w, r, http.StatusConflict, ReturnStatus(h.Status.DuplicateKey), h.Log, h.Resource, h.Config.Create, false, "Duplicate Key")
				return
			}
		}
	} else if er2 != nil {
		ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Create, er2, h.Log)
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
	er1 := MatchId(r, body, h.Keys, h.mapIndex)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	if h.modelBuilder != nil {
		body, er0 = h.modelBuilder.BuildToUpdate(r.Context(), body)
		if er0 != nil {
			ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Update, er0, h.Log)
			return
		}
	}
	if h.Validate != nil {
		errors, er2 := h.Validate(r.Context(), body)
		if er2 != nil {
			ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Update, er2, h.Log)
			return
		}
		if len(errors) > 0 {
			result0 := ResultInfo{Status: *h.Status.ValidationError, Errors: errors}
			RespondAndLog(w, r, http.StatusUnprocessableEntity, result0, h.Log, h.Resource, h.Config.Update, false, "Data Validation Failed")
			return
		}
	}
	count, er3 := h.service.Update(r.Context(), body)
	if er3 != nil {
		ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Update, er3, h.Log)
		return
	}
	if count == -1 {
		RespondAndLog(w, r, http.StatusConflict, ReturnStatus(h.Status.VersionError), h.Log, h.Resource, h.Config.Update, false, "Data Version Error")
	} else if count == 0 {
		RespondAndLog(w, r, http.StatusNotFound, ReturnStatus(h.Status.NotFound), h.Log, h.Resource, h.Config.Update, false, "Data Not Found")
	} else {
		Succeed(w, r, http.StatusOK, SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Update)
	}
}

func (h *GenericHandler) Patch(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	bodyStruct := reflect.New(h.ModelType).Interface()
	body0, er0 := BuildMapAndStruct(r, bodyStruct)
	if er0 != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}
	er1 := MatchId(r, bodyStruct, h.Keys, h.mapIndex)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	body, er2 := BodyToJson(r, bodyStruct, body0, h.Keys, h.mapIndex, h.modelBuilder)
	if er2 != nil {
		// http.Error(w, "Invalid Data: "+er2.Error(), http.StatusBadRequest)
		ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Patch, er2, h.Log)
		return
	}
	if h.Validate != nil {
		errors, er3 := h.Validate(r.Context(), bodyStruct)
		if er3 != nil {
			ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Patch, er3, h.Log)
			return
		}
		if len(errors) > 0 {
			result0 := ResultInfo{Status: *h.Status.ValidationError, Errors: errors}
			RespondAndLog(w, r, http.StatusUnprocessableEntity, result0, h.Log, h.Resource, h.Config.Patch, false, "Data Validation Failed")
			return
		}
	}
	count, er4 := h.service.Patch(r.Context(), body)
	if er4 != nil {
		ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Patch, er4, h.Log)
		return
	}
	if count == -1 {
		RespondAndLog(w, r, http.StatusConflict, ReturnStatus(h.Status.VersionError), h.Log, h.Resource, h.Config.Patch, false, "Data Version Error")
	} else if count == 0 {
		RespondAndLog(w, r, http.StatusNotFound, ReturnStatus(h.Status.NotFound), h.Log, h.Resource, h.Config.Patch, false, "Data Not Found")
	} else {
		Succeed(w, r, http.StatusOK, SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Patch)
	}
}

func (h *GenericHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, er1 := BuildId(r, h.ModelType, h.Keys, h.Indexes, 0)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	count, er2 := h.service.Delete(r.Context(), id)
	if er2 != nil {
		ErrorAndLog(w, r, http.StatusInternalServerError, InternalServerError, h.Error, h.Resource, h.Config.Delete, er2, h.Log)
		return
	}
	if count > 0 {
		Succeed(w, r, http.StatusOK, count, h.Log, h.Resource, h.Config.Delete)
	} else if count == 0 {
		RespondAndLog(w, r, http.StatusNotFound, count, h.Log, h.Resource, h.Config.Delete, false, "Data Not Found")
	} else {
		RespondAndLog(w, r, http.StatusConflict, count, h.Log, h.Resource, h.Config.Delete, false, "Conflict")
	}
}
