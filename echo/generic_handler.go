package echo

import (
	"context"
	sv "github.com/core-go/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

const (
	Method = "method"
	Update = "update"
	Patch  = "patch"
)

type GenericHandler struct {
	*LoadHandler
	Status       sv.StatusConfig
	Config       sv.ActionConfig
	service      sv.HGenericService
	modelBuilder sv.ModelBuilder
	Validate     func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	Log          func(ctx context.Context, resource string, action string, success bool, desc string) error
	mapIndex     map[string]int
}

func NewGenericHandler(genericService sv.HGenericService, modelType reflect.Type, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	return NewGenericHandlerWithConfig(genericService, modelType, nil, modelBuilder, logError, validate, options...)
}
func NewGenericHandlerWithKeys(genericService sv.HGenericService, keys []string, modelType reflect.Type, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	return NewGenericHandlerWithKeysAndLog(genericService, keys, modelType, nil, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewGenericHandlerWithConfig(genericService sv.HGenericService, modelType reflect.Type, status *sv.StatusConfig, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 && options[0] != nil {
		writeLog = options[0]
	}
	return NewGenericHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewGenericHandlerWithLog(genericService sv.HGenericService, modelType reflect.Type, status *sv.StatusConfig, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *sv.ActionConfig) *GenericHandler {
	return NewGenericHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, resource, conf)
}
func NewGenericHandlerWithKeysAndLog(genericService sv.HGenericService, keys []string, modelType reflect.Type, status *sv.StatusConfig, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *sv.ActionConfig) *GenericHandler {
	if keys == nil || len(keys) == 0 {
		keys = sv.GetJsonPrimaryKeys(modelType)
	}
	if len(resource) == 0 {
		resource = sv.BuildResourceName(modelType.Name())
	}
	var c sv.ActionConfig
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
	var s sv.StatusConfig
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
	_, jsonMapIndex := sv.BuildMapField(modelType)

	return &GenericHandler{LoadHandler: loadHandler, service: genericService, Status: s, modelBuilder: modelBuilder, Validate: validate, mapIndex: jsonMapIndex, Log: writeLog, Config: c}
}
func (h *GenericHandler) Create(ctx echo.Context) error {
	return h.Insert(ctx)
}
func (h *GenericHandler) Insert(ctx echo.Context) error {
	r := ctx.Request()
	body, er0 := sv.NewModel(h.ModelType, r.Body)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "Invalid Request")
		return er0
	}
	if h.modelBuilder != nil {
		body, er0 = h.modelBuilder.BuildToInsert(r.Context(), body)
		if er0 != nil {
			return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Create, er0, h.Log)
		}
	}
	if h.Validate != nil {
		errors, er1 := h.Validate(r.Context(), body)
		if er1 != nil {
			return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Create, er1, h.Log)
		}
		if len(errors) > 0 {
			//result0 := model.ResultInfo{Status: model.StatusError, Errors: MakeErrors(errors)}
			result0 := sv.ResultInfo{Status: *h.Status.ValidationError, Errors: errors}
			return RespondAndLog(ctx, http.StatusUnprocessableEntity, result0, h.Log, h.Resource, h.Config.Create, false, "Data Validation Failed")
		}
	}
	var count int64
	var er2 error
	count, er2 = h.service.Insert(r.Context(), body)
	if count <= 0 && er2 == nil {
		if h.modelBuilder == nil {
			return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.DuplicateKey), h.Log, h.Resource, h.Config.Create, false, "Duplicate Key")
		}
		i := 0
		for count <= 0 && i <= 5 {
			i++
			body, er2 = h.modelBuilder.BuildToInsert(r.Context(), body)
			if er2 != nil {
				return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Create, er2, h.Log)
			}
			count, er2 = h.service.Insert(r.Context(), body)
			if er2 != nil {
				return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Create, er2, h.Log)
			}
			if count > 0 {
				return Succeed(ctx, http.StatusCreated, sv.SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Create)
			}
			if i == 5 {
				return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.DuplicateKey), h.Log, h.Resource, h.Config.Create, false, "Duplicate Key")
			}
		}
	} else if er2 != nil {
		return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Create, er2, h.Log)
	}
	return Succeed(ctx, http.StatusCreated, sv.SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Create)
}
func (h *GenericHandler) Update(ctx echo.Context) error {
	r := ctx.Request()
	r = r.WithContext(context.WithValue(r.Context(), Method, Update))
	body, er0 := sv.NewModel(h.ModelType, r.Body)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "Invalid Data")
		return er0
	}
	er1 := sv.CheckId(r, body, h.Keys, h.mapIndex)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return er1
	}
	if h.modelBuilder != nil {
		body, er0 = h.modelBuilder.BuildToUpdate(r.Context(), body)
		if er0 != nil {
			return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Update, er0, h.Log)
		}
	}
	if h.Validate != nil {
		errors, er2 := h.Validate(r.Context(), body)
		if er2 != nil {
			return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Update, er2, h.Log)
		}
		if len(errors) > 0 {
			result0 := sv.ResultInfo{Status: *h.Status.ValidationError, Errors: errors}
			return RespondAndLog(ctx, http.StatusUnprocessableEntity, result0, h.Log, h.Resource, h.Config.Update, false, "Data Validation Failed")
		}
	}
	count, er3 := h.service.Update(r.Context(), body)
	if er3 != nil {
		return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Update, er3, h.Log)
	}
	if count == -1 {
		return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.VersionError), h.Log, h.Resource, h.Config.Update, false, "Data Version Error")
	} else if count == 0 {
		return RespondAndLog(ctx, http.StatusNotFound, sv.ReturnStatus(h.Status.NotFound), h.Log, h.Resource, h.Config.Update, false, "Data Not Found")
	} else {
		return Succeed(ctx, http.StatusOK, sv.SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Update)
	}
}

func (h *GenericHandler) Patch(ctx echo.Context) error {
	r := ctx.Request()
	r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	bodyStruct := reflect.New(h.ModelType).Interface()
	body0, er0 := sv.BuildMapAndStruct(r, bodyStruct)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "Invalid Data")
		return er0
	}
	er1 := sv.CheckId(r, bodyStruct, h.Keys, h.mapIndex)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return er1
	}
	body, er2 := sv.BodyToJson(r, bodyStruct, body0, h.Keys, h.mapIndex, h.modelBuilder)
	if er2 != nil {
		return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Patch, er2, h.Log)
	}
	if h.Validate != nil {
		errors, er3 := h.Validate(r.Context(), &bodyStruct)
		if er3 != nil {
			return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Patch, er3, h.Log)
		}
		if len(errors) > 0 {
			result0 := sv.ResultInfo{Status: *h.Status.ValidationError, Errors: errors}
			return RespondAndLog(ctx, http.StatusUnprocessableEntity, result0, h.Log, h.Resource, h.Config.Patch, false, "Data Validation Failed")
		}
	}
	count, er4 := h.service.Patch(r.Context(), body)
	if er4 != nil {
		return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Patch, er4, h.Log)
	}
	if count == -1 {
		return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.VersionError), h.Log, h.Resource, h.Config.Patch, false, "Data Version Error")
	} else if count == 0 {
		return RespondAndLog(ctx, http.StatusNotFound, sv.ReturnStatus(h.Status.NotFound), h.Log, h.Resource, h.Config.Patch, false, "Data Not Found")
	} else {
		return Succeed(ctx, http.StatusOK, sv.SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Config.Patch)
	}
}

func (h *GenericHandler) Delete(ctx echo.Context) error {
	r := ctx.Request()
	id, er1 := sv.BuildId(r, h.ModelType, h.Keys, h.Indexes, 0)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, "cannot parse form: "+er1.Error())
		return er1
	}
	count, er2 := h.service.Delete(r.Context(), id)
	if er2 != nil {
		return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Delete, er2, h.Log)
	}
	if count > 0 {
		return Succeed(ctx, http.StatusOK, count, h.Log, h.Resource, h.Config.Delete)
	} else if count == 0 {
		return RespondAndLog(ctx, http.StatusNotFound, count, h.Log, h.Resource, h.Config.Delete, false, "Data Not Found")
	} else {
		return RespondAndLog(ctx, http.StatusConflict, count, h.Log, h.Resource, h.Config.Delete, false, "Conflict")
	}
}
