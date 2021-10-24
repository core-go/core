package gin

import (
	"context"
	sv "github.com/core-go/service"
	"github.com/gin-gonic/gin"
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
	Action       sv.ActionConfig
	service      sv.HGenericService
	modelBuilder sv.ModelBuilder
	Validate     func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	Log          func(ctx context.Context, resource string, action string, success bool, desc string) error
	Indexes      map[string]int
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

	return &GenericHandler{LoadHandler: loadHandler, service: genericService, Status: s, modelBuilder: modelBuilder, Validate: validate, Indexes: jsonMapIndex, Log: writeLog, Action: c}
}
func (h *GenericHandler) Create(ctx *gin.Context) {
	h.Insert(ctx)
}
func (h *GenericHandler) Insert(ctx *gin.Context) {
	r := ctx.Request
	body, er0 := sv.NewModel(h.ModelType, r.Body)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "Invalid Request")
		return
	}
	if h.modelBuilder != nil {
		body, er0 = h.modelBuilder.BuildToInsert(r.Context(), body)
		if er0 != nil {
			ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Create, er0, h.Log)
		}
	}
	if h.Validate != nil {
		errors, er1 := h.Validate(r.Context(), body)
		if HasError(ctx, errors, er1, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Update) {
			return
		}
	}
	var count int64
	var er2 error
	count, er2 = h.service.Insert(r.Context(), body)
	if count <= 0 && er2 == nil {
		if h.modelBuilder == nil {
			ReturnAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.DuplicateKey), h.Log, false, h.Resource, h.Action.Create, "Duplicate Key")
			return
		}
		i := 0
		for count <= 0 && i <= 5 {
			i++
			body, er2 = h.modelBuilder.BuildToInsert(r.Context(), body)
			if er2 != nil {
				ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Create, er2, h.Log)
				return
			}
			count, er2 = h.service.Insert(r.Context(), body)
			if er2 != nil {
				ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Create, er2, h.Log)
				return
			}
			if count > 0 {
				Succeed(ctx, http.StatusCreated, sv.SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Action.Create)
				return
			}
			if i == 5 {
				ReturnAndLog(ctx, http.StatusConflict, sv.ReturnStatus(h.Status.DuplicateKey), h.Log, false, h.Resource, h.Action.Create, "Duplicate Key")
				return
			}
		}
	} else if er2 != nil {
		ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Create, er2, h.Log)
		return
	}
	Succeed(ctx, http.StatusCreated, sv.SetStatus(body, h.Status.Success), h.Log, h.Resource, h.Action.Create)
}
func (h *GenericHandler) Update(ctx *gin.Context) {
	r := ctx.Request
	r = r.WithContext(context.WithValue(r.Context(), Method, Update))
	body, er0 := sv.NewModel(h.ModelType, r.Body)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "Invalid Data")
		return
	}
	er1 := sv.MatchId(r, body, h.Keys, h.Indexes)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return
	}
	if h.modelBuilder != nil {
		body, er0 = h.modelBuilder.BuildToUpdate(r.Context(), body)
		if er0 != nil {
			ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Update, er0, h.Log)
			return
		}
	}
	if h.Validate != nil {
		errors, er2 := h.Validate(r.Context(), body)
		if HasError(ctx, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Update) {
			return
		}
	}
	count, er3 := h.service.Update(r.Context(), body)
	HandleResult(ctx, body, count, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Update)
}

func (h *GenericHandler) Patch(ctx *gin.Context) {
	r := ctx.Request
	r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	bodyStruct := reflect.New(h.ModelType).Interface()
	body0, er0 := sv.BuildMapAndStruct(r, bodyStruct)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "Invalid Data")
		return
	}
	er1 := CheckId(ctx, bodyStruct, h.Keys, h.Indexes)
	if er1 != nil {
		return
	}
	body, er2 := sv.BodyToJsonMap(r, bodyStruct, body0, h.Keys, h.Indexes, h.modelBuilder.BuildToPatch)
	if er2 != nil {
		ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Patch, er2, h.Log)
		return
	}
	if h.Validate != nil {
		errors, er3 := h.Validate(r.Context(), &bodyStruct)
		if HasError(ctx, errors, er3, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Patch) {
			return
		}
	}
	count, er4 := h.service.Patch(r.Context(), body)
	HandleResult(ctx, body, count, er4, h.Status, h.Error, h.Log, h.Resource, h.Action.Patch)
}
func (h *GenericHandler) Delete(ctx *gin.Context) {
	r := ctx.Request
	id, er1 := sv.BuildId(r, h.ModelType, h.Keys, h.KeyIndexes)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, "cannot parse form: "+er1.Error())
		return
	}
	count, er2 := h.service.Delete(r.Context(), id)
	HandleDelete(ctx, count, er2, h.Error, h.Log, h.Resource, h.Action.Delete)
}
func CheckId(ctx *gin.Context, body interface{}, keysJson []string, mapIndex map[string]int) error {
	err := sv.MatchId(ctx.Request, body, keysJson, mapIndex)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	}
	return err
}
func HandleDelete(ctx *gin.Context, count int64, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		return
	}
	if count > 0 {
		Succeed(ctx, http.StatusOK, count, writeLog, resource, action)
	} else if count == 0 {
		ReturnAndLog(ctx, http.StatusNotFound, count, writeLog, false, resource, action, "Data Not Found")
	} else {
		ReturnAndLog(ctx, http.StatusConflict, count, writeLog, false, resource, action, "Conflict")
	}
}
func BodyToJsonWithBuild(ctx *gin.Context, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int, buildToPatch func(context.Context, interface{}) (interface{}, error), logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) (map[string]interface{}, error) {
	body, err := sv.BodyToJsonMap(ctx.Request, structBody, body, jsonIds, mapIndex, buildToPatch)
	if err != nil {
		// http.Error(w, "Invalid Data: "+err.Error(), http.StatusBadRequest)
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	}
	return body, err
}
func BodyToJson(ctx *gin.Context, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int) (map[string]interface{}, error) {
	body, err := sv.BodyToJsonMap(ctx.Request, structBody, body, jsonIds, mapIndex)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid Data: "+err.Error())
	}
	return body, err
}
func BuildFieldMapAndCheckId(ctx *gin.Context, obj interface{}, keysJson []string, mapIndex map[string]int) (map[string]interface{}, error) {
	body, er0 := sv.BuildMapAndStruct(ctx.Request, obj)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, er0.Error())
		return body, er0
	}
	er1 := CheckId(ctx, obj, keysJson, mapIndex)
	return body, er1
}
func BuildMapAndCheckId(ctx *gin.Context, obj interface{}, keysJson []string, mapIndex map[string]int) (map[string]interface{}, error) {
	body, er0 := BuildFieldMapAndCheckId(ctx, obj, keysJson, mapIndex)
	if er0 != nil {
		return body, er0
	}
	json, er1 := sv.BodyToJsonMap(ctx.Request, obj, body, keysJson, mapIndex)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	}
	return json, er1
}
func HasError(ctx *gin.Context, errors []sv.ErrorMessage, err error, status int, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) bool {
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		return true
	}
	if len(errors) > 0 {
		result0 := sv.ResultInfo{Status: status, Errors: errors}
		ReturnAndLog(ctx, http.StatusUnprocessableEntity, result0, writeLog, false, resource, action, "Data Validation Failed")
		return true
	}
	return false
}
func HandleResult(ctx *gin.Context, body interface{}, count int64, err error, status sv.StatusConfig, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		return
	}
	if count == -1 {
		ReturnAndLog(ctx, http.StatusConflict, sv.ReturnStatus(status.VersionError), writeLog, false, resource, action, "Data Version Error")
	} else if count == 0 {
		ReturnAndLog(ctx, http.StatusNotFound, sv.ReturnStatus(status.NotFound), writeLog, false, resource, action, "Data Not Found")
	} else {
		Succeed(ctx, http.StatusOK, sv.SetStatus(body, status.Success), writeLog, resource, action)
	}
}
func AfterCreated(ctx *gin.Context, body interface{}, count int64, err error, status sv.StatusConfig, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		return
	}
	if count <= 0 {
		ReturnAndLog(ctx, http.StatusConflict, sv.ReturnStatus(status.DuplicateKey), writeLog, false, resource, action, "Duplicate Key")
	} else {
		Succeed(ctx, http.StatusCreated, sv.SetStatus(body, status.Success), writeLog, resource, action)
	}
}
func Respond(ctx *gin.Context, code int, result interface{}, err error, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	} else {
		Succeed(ctx, code, result, writeLog, resource, action)
	}
}
