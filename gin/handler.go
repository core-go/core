package gin

import (
	"context"
	"encoding/json"
	sv "github.com/core-go/core"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

const (
	Method = "method"
	Update = "update"
	Patch  = "patch"
)

func CreatePatchAndParams(modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), patch func(context.Context, map[string]interface{}) (int64, error), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), build func(context.Context, interface{}) (interface{}, error), action *sv.ActionConfig, options ...func(context.Context, string, string, bool, string) error) (*PatchHandler, *sv.Params) {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	a := sv.InitializeAction(action)
	resource := sv.BuildResourceName(modelType.Name())
	keys, indexes, _ := sv.BuildMapField(modelType)
	patchHandler := &PatchHandler{PrimaryKeys: keys, FieldIndexes: indexes, ObjectType: modelType, Save: patch, ValidateData: validate, Build: build, LogError: logError, WriteLog: writeLog, ResourceType: resource, Activity: a.Patch}
	params := &sv.Params{Keys: keys, Indexes: indexes, ModelType: modelType, Resource: resource, Action: a, Error: logError, Log: writeLog, Validate: validate}
	return patchHandler, params
}
func NewPatchHandler(patch func(context.Context, map[string]interface{}) (int64, error), modelType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), build func(context.Context, interface{}) (interface{}, error), action string, options ...func(context.Context, string, string, bool, string) error) *PatchHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	resource := sv.BuildResourceName(modelType.Name())
	keys, indexes, _ := sv.BuildMapField(modelType)
	return &PatchHandler{PrimaryKeys: keys, FieldIndexes: indexes, ObjectType: modelType, Save: patch, ValidateData: validate, Build: build, LogError: logError, WriteLog: writeLog, ResourceType: resource, Activity: action}
}

type PatchHandler struct {
	PrimaryKeys  []string
	FieldIndexes map[string]int
	ObjectType   reflect.Type
	Save         func(ctx context.Context, user map[string]interface{}) (int64, error)
	ValidateData func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	Build        func(ctx context.Context, model interface{}) (interface{}, error)
	LogError     func(context.Context, string, ...map[string]interface{})
	WriteLog     func(context.Context, string, string, bool, string) error
	ResourceType string
	Activity     string
}
func (h *PatchHandler) Patch(ctx *gin.Context) {
	r := ctx.Request
	r = r.WithContext(context.WithValue(r.Context(), Method, Patch))
	bodyStruct := reflect.New(h.ObjectType).Interface()
	body0, er0 := sv.BuildMapAndStruct(r, bodyStruct)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "Invalid Data")
		return
	}
	er1 := CheckId(ctx, bodyStruct, h.PrimaryKeys, h.FieldIndexes)
	if er1 != nil {
		return
	}
	body, er2 := sv.BodyToJsonMap(r, bodyStruct, body0, h.PrimaryKeys, h.FieldIndexes, h.Build)
	if er2 != nil {
		ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.LogError, h.ResourceType, h.Activity, er2, h.WriteLog)
		return
	}
	if h.ValidateData != nil {
		errors, er3 := h.ValidateData(r.Context(), &bodyStruct)
		if HasError(ctx, errors, er3, h.LogError, h.WriteLog, h.ResourceType, h.Activity) {
			return
		}
	}
	count, er4 := h.Save(r.Context(), body)
	HandleResult(ctx, body, count, er4, h.LogError, h.WriteLog, h.ResourceType, h.Activity)
}
type GenericHandler struct {
	*LoadHandler
	Action   sv.ActionConfig
	service  sv.SimpleService
	builder  sv.Builder
	Validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	Log      func(ctx context.Context, resource string, action string, success bool, desc string) error
	Indexes  map[string]int
}

func NewHandler(genericService sv.SimpleService, modelType reflect.Type, modelBuilder sv.Builder, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	return NewHandlerWithConfig(genericService, modelType, modelBuilder, logError, validate, options...)
}
func NewHandlerWithKeys(genericService sv.SimpleService, keys []string, modelType reflect.Type, modelBuilder sv.Builder, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 {
		writeLog = options[0]
	}
	return NewHandlerWithKeysAndLog(genericService, keys, modelType, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewHandlerWithConfig(genericService sv.SimpleService, modelType reflect.Type, modelBuilder sv.Builder, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), options ...func(context.Context, string, string, bool, string) error) *GenericHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 && options[0] != nil {
		writeLog = options[0]
	}
	return NewHandlerWithKeysAndLog(genericService, nil, modelType, modelBuilder, logError, validate, writeLog, "", nil)
}
func NewHandlerWithLog(genericService sv.SimpleService, modelType reflect.Type, modelBuilder sv.Builder, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *sv.ActionConfig) *GenericHandler {
	return NewHandlerWithKeysAndLog(genericService, nil, modelType, modelBuilder, logError, validate, writeLog, resource, conf)
}
func NewHandlerWithKeysAndLog(genericService sv.SimpleService, keys []string, modelType reflect.Type, modelBuilder sv.Builder, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), writeLog func(context.Context, string, string, bool, string) error, resource string, conf *sv.ActionConfig) *GenericHandler {
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
	loadHandler := NewLoadHandlerWithKeysAndLog(genericService.Load, keys, modelType, logError, writeLog2, *c.Load, resource)
	_, jsonMapIndex, _ := sv.BuildMapField(modelType)

	return &GenericHandler{LoadHandler: loadHandler, service: genericService, builder: modelBuilder, Validate: validate, Indexes: jsonMapIndex, Log: writeLog, Action: c}
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
	if h.builder != nil {
		body, er0 = h.builder.Create(r.Context(), body)
		if er0 != nil {
			ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Create, er0, h.Log)
		}
	}
	if h.Validate != nil {
		errors, er1 := h.Validate(r.Context(), body)
		if HasError(ctx, errors, er1, h.Error, h.Log, h.Resource, h.Action.Update) {
			return
		}
	}
	var count int64
	var er2 error
	count, er2 = h.service.Insert(r.Context(), body)
	if count <= 0 && er2 == nil {
		if h.builder == nil {
			ReturnAndLog(ctx, http.StatusConflict, -1, h.Log, false, h.Resource, h.Action.Create, "Duplicate Key")
			return
		}
		i := 0
		for count <= 0 && i <= 5 {
			i++
			body, er2 = h.builder.Create(r.Context(), body)
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
				Succeed(ctx, http.StatusCreated, body, h.Log, h.Resource, h.Action.Create)
				return
			}
			if i == 5 {
				ReturnAndLog(ctx, http.StatusConflict, count, h.Log, false, h.Resource, h.Action.Create, "Duplicate Key")
				return
			}
		}
	} else if er2 != nil {
		ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Create, er2, h.Log)
		return
	}
	Succeed(ctx, http.StatusCreated, body, h.Log, h.Resource, h.Action.Create)
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
	if h.builder != nil {
		body, er0 = h.builder.Update(r.Context(), body)
		if er0 != nil {
			ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Update, er0, h.Log)
			return
		}
	}
	if h.Validate != nil {
		errors, er2 := h.Validate(r.Context(), body)
		if HasError(ctx, errors, er2, h.Error, h.Log, h.Resource, h.Action.Update) {
			return
		}
	}
	count, er3 := h.service.Update(r.Context(), body)
	HandleResult(ctx, body, count, er3, h.Error, h.Log, h.Resource, h.Action.Update)
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
	body, er2 := sv.BodyToJsonMap(r, bodyStruct, body0, h.Keys, h.Indexes, h.builder.Patch)
	if er2 != nil {
		ErrorAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, h.Error, h.Resource, h.Action.Patch, er2, h.Log)
		return
	}
	if h.Validate != nil {
		errors, er3 := h.Validate(r.Context(), &bodyStruct)
		if HasError(ctx, errors, er3, h.Error, h.Log, h.Resource, h.Action.Patch) {
			return
		}
	}
	count, er4 := h.service.Patch(r.Context(), body)
	HandleResult(ctx, body, count, er4, h.Error, h.Log, h.Resource, h.Action.Patch)
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
func HandleDelete(ctx *gin.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
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
func BodyToJsonWithBuild(ctx *gin.Context, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int, buildToPatch func(context.Context, interface{}) (interface{}, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) (map[string]interface{}, error) {
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
func BuildMapAndCheckId(ctx *gin.Context, obj interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) (interface{}, error)) (map[string]interface{}, error) {
	body, er0 := BuildFieldMapAndCheckId(ctx, obj, keysJson, mapIndex)
	if er0 != nil {
		return body, er0
	}
	json, er1 := sv.BodyToJsonMap(ctx.Request, obj, body, keysJson, mapIndex, options...)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	}
	return json, er1
}
func HasError(ctx *gin.Context, errors []sv.ErrorMessage, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) bool {
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		return true
	}
	if len(errors) > 0 {
		ReturnAndLog(ctx, http.StatusUnprocessableEntity, errors, writeLog, false, resource, action, "Data Validation Failed")
		return true
	}
	return false
}
func HandleResult(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		if sv.IsNil(body) {
			RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		} else {
			logError(ctx.Request.Context(), err.Error(), sv.MakeMap(body))
			RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, nil, writeLog, resource, action)
		}
		return
	}
	if count == -1 {
		ReturnAndLog(ctx, http.StatusConflict, count, writeLog, false, resource, action, "Data Version Error")
	} else if count == 0 {
		ReturnAndLog(ctx, http.StatusNotFound, 0, writeLog, false, resource, action, "Data Not Found")
	} else {
		if sv.IsNil(body) {
			Succeed(ctx, http.StatusOK, count, writeLog, resource, action)
		} else {
			Succeed(ctx, http.StatusOK, body, writeLog, resource, action)
		}
	}
}
func AfterCreated(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		if sv.IsNil(body) {
			RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		} else {
			logError(ctx.Request.Context(), err.Error(), sv.MakeMap(body))
			RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
		}
		return
	}
	if count <= 0 {
		ReturnAndLog(ctx, http.StatusConflict, count, writeLog, false, resource, action, "Duplicate Key")
	} else {
		Succeed(ctx, http.StatusCreated, body, writeLog, resource, action)
	}
}
func Respond(ctx *gin.Context, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
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
/*
func Return222(ctx *gin.Context, code int, result sv.ResultInfo, status sv.StatusConfig, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
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
		if code == http.StatusCreated {
			if result.Status == status.DuplicateKey {
				Succeed(ctx, http.StatusConflict, result, writeLog, resource, action)
			} else {
				Succeed(ctx, code, result, writeLog, resource, action)
			}
		} else {
			if result.Status == status.NotFound {
				Succeed(ctx, http.StatusNotFound, result, writeLog, resource, action)
			} else if result.Status == status.VersionError {
				Succeed(ctx, http.StatusConflict, result, writeLog, resource, action)
			} else {
				Succeed(ctx, code, result, writeLog, resource, action)
			}
		}
	}
}
 */
func Result(ctx *gin.Context, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), opts...interface{}) {
	if err != nil {
		if len(opts) > 0 && opts[0] != nil {
			b, er2 := json.Marshal(opts[0])
			if er2 != nil {
				m := make(map[string]interface{})
				m["request"] = string(b)
				logError(ctx.Request.Context(), err.Error(), m)
			} else {
				logError(ctx.Request.Context(), err.Error())
			}
			ctx.String(http.StatusInternalServerError, sv.InternalServerError)
		} else {
			logError(ctx.Request.Context(), err.Error(), nil)
			ctx.String(http.StatusInternalServerError, sv.InternalServerError)
		}
	} else {
		ctx.JSON(code, result)
	}
}
