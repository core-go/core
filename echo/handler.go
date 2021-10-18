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

func NewHandler(genericService sv.HGenericService, modelType reflect.Type, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	return NewHandlerWithConfig(genericService, modelType, nil, modelBuilder, logError, validate, options...)
}
func NewHandlerWithKeys(genericService sv.HGenericService, keys []string, modelType reflect.Type, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	return NewHandlerWithKeysAndLog(genericService, keys, modelType, nil, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewHandlerWithConfig(genericService sv.HGenericService, modelType reflect.Type, status *sv.StatusConfig, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 && options[0] != nil {
		writeLog = options[0]
	}
	return NewHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewHandlerWithLog(genericService sv.HGenericService, modelType reflect.Type, status *sv.StatusConfig, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *sv.ActionConfig) *GenericHandler {
	return NewHandlerWithKeysAndLog(genericService, nil, modelType, status, modelBuilder, logError, validate, writeLog, resource, conf)
}
func NewHandlerWithKeysAndLog(genericService sv.HGenericService, keys []string, modelType reflect.Type, status *sv.StatusConfig, modelBuilder sv.ModelBuilder, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *sv.ActionConfig) *GenericHandler {
	if keys == nil || len(keys) == 0 {
		keys = sv.GetJsonPrimaryKeys(modelType)
	}
	if len(resource) == 0 {
		resource = sv.BuildResourceName(modelType.Name())
	}
	var writeLog2 func(context.Context, string, string, bool, string) error
	if conf != nil && conf.Load != nil {
		writeLog2 = writeLog
	}
	c := sv.InitializeAction(conf)
	s := sv.InitializeStatus(status)
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
			return RespondAndLog(ctx, http.StatusUnprocessableEntity, result0, h.Log, false, h.Resource, h.Config.Create, "Data Validation Failed")
		}
	}
	var count int64
	var er2 error
	count, er2 = h.service.Insert(r.Context(), body)
	if count <= 0 && er2 == nil {
		if h.modelBuilder == nil {
			return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.DuplicateKey), h.Log, false, h.Resource, h.Config.Create, "Duplicate Key")
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
				return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.DuplicateKey), h.Log, false, h.Resource, h.Config.Create, "Duplicate Key")
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
	er1 := CheckId(ctx, body, h.Keys, h.mapIndex)
	if er1 != nil {
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
		if HasError(ctx, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Config.Update) {
			return er2
		}
	}
	count, er3 := h.service.Update(r.Context(), body)
	return HandleResult(ctx, body, count, er3, h.Status, h.Error, h.Log, h.Resource, h.Config.Patch)
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
	er1 := CheckId(ctx, bodyStruct, h.Keys, h.mapIndex)
	if er1 != nil {
		return er1
	}
	body, er2 := sv.BodyToJsonMap(r, bodyStruct, body0, h.Keys, h.mapIndex, h.modelBuilder.BuildToPatch)
	if er2 != nil {
		return ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Config.Patch, er2, h.Log)
	}
	if h.Validate != nil {
		errors, er3 := h.Validate(r.Context(), &bodyStruct)
		if HasError(ctx, errors, er3, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Config.Patch) {
			return er3
		}
	}
	count, er4 := h.service.Patch(r.Context(), body)
	return HandleResult(ctx, body, count, er4, h.Status, h.Error, h.Log, h.Resource, h.Config.Patch)
}
func (h *GenericHandler) Delete(ctx echo.Context) error {
	r := ctx.Request()
	id, er1 := sv.BuildId(r, h.ModelType, h.Keys, h.Indexes, 0)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, "cannot parse form: "+er1.Error())
		return er1
	}
	count, er2 := h.service.Delete(r.Context(), id)
	return HandleDelete(ctx, count, er2, h.Error, h.Log, h.Resource, h.Config.Delete)
}

func CheckId(ctx echo.Context, body interface{}, keysJson []string, mapIndex map[string]int) error {
	err := sv.MatchId(ctx.Request(), body, keysJson, mapIndex)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	}
	return err
}
func HandleDelete(ctx echo.Context, count int64, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	if err != nil {
		return Respond(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	}
	if count > 0 {
		return Succeed(ctx, http.StatusOK, count, writeLog, resource, action)
	} else if count == 0 {
		return RespondAndLog(ctx, http.StatusNotFound, count, writeLog, false, resource, action, "Data Not Found")
	} else {
		return RespondAndLog(ctx, http.StatusConflict, count, writeLog, false, resource, action, "Conflict")
	}
}
func BodyToJsonWithBuild(ctx echo.Context, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int, buildToPatch func(context.Context, interface{}) (interface{}, error), logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) (map[string]interface{}, error) {
	body, err := sv.BodyToJsonMap(ctx.Request(), structBody, body, jsonIds, mapIndex, buildToPatch)
	if err != nil {
		// http.Error(w, "Invalid Data: "+err.Error(), http.StatusBadRequest)
		Respond(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	}
	return body, err
}
func BodyToJson(ctx echo.Context, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int) (map[string]interface{}, error) {
	body, err := sv.BodyToJsonMap(ctx.Request(), structBody, body, jsonIds, mapIndex)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid Data: "+err.Error())
	}
	return body, err
}
func BuildFieldMapAndCheckId(ctx echo.Context, obj interface{}, keysJson []string, mapIndex map[string]int) (map[string]interface{}, error) {
	body, er0 := sv.BuildMapAndStruct(ctx.Request(), obj)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, er0.Error())
		return body, er0
	}
	er1 := CheckId(ctx, obj, keysJson, mapIndex)
	return body, er1
}
func BuildMapAndCheckId(ctx echo.Context, obj interface{}, keysJson []string, mapIndex map[string]int) (map[string]interface{}, error) {
	body, er0 := BuildFieldMapAndCheckId(ctx, obj, keysJson, mapIndex)
	if er0 != nil {
		return body, er0
	}
	json, er1 := sv.BodyToJsonMap(ctx.Request(), obj, body, keysJson, mapIndex)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	}
	return json, er1
}
func HasError(ctx echo.Context, errors []sv.ErrorMessage, err error, status int, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) bool {
	if err != nil {
		Respond(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		return true
	}
	if len(errors) > 0 {
		result0 := sv.ResultInfo{Status: status, Errors: errors}
		RespondAndLog(ctx, http.StatusUnprocessableEntity, result0, writeLog, false, resource, action, "Data Validation Failed")
		return true
	}
	return false
}
func HandleResult(ctx echo.Context, body interface{}, count int64, err error, status sv.StatusConfig, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	if err != nil {
		return Respond(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	}
	if count == -1 {
		return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(status.VersionError), writeLog, false, resource, action, "Data Version Error")
	} else if count == 0 {
		return RespondAndLog(ctx, http.StatusNotFound, sv.ReturnStatus(status.NotFound), writeLog, false, resource, action, "Data Not Found")
	} else {
		return Succeed(ctx, http.StatusOK, sv.SetStatus(body, status.Success), writeLog, resource, action)
	}
}
func AfterCreated(ctx echo.Context, body interface{}, count int64, err error, status sv.StatusConfig, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	if err != nil {
		return Respond(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	}
	if count <= 0 {
		return RespondAndLog(ctx, http.StatusConflict, sv.ReturnStatus(status.DuplicateKey), writeLog, false, resource, action, "Duplicate Key")
	} else {
		return Succeed(ctx, http.StatusCreated, sv.SetStatus(body, status.Success), writeLog, resource, action)
	}
}
