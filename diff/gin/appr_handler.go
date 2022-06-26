package gin

import (
	"context"
	d "github.com/core-go/core/diff"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

const internalServerError = "Internal Server Error"

type ApprHandler struct {
	ApprService d.ApprService
	Keys        []string
	ModelType   reflect.Type
	Error       func(context.Context, string)
	Indexes     map[string]int
	Offset      int
	Log         func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource    string
	Action1     string
	Action2     string
}

func NewApprHandler(apprService d.ApprService, modelType reflect.Type, logError func(context.Context, string), option ...int) *ApprHandler {
	offset := 1
	if len(option) > 0 && option[0] >= 0 {
		offset = option[0]
	}
	return NewApprHandlerWithKeysAndLog(apprService, nil, modelType, offset, logError, nil)
}
func NewApprHandlerWithLogs(apprService d.ApprService, modelType reflect.Type, offset int, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprHandler {
	return NewApprHandlerWithKeysAndLog(apprService, nil, modelType, offset, logError, writeLog, options...)
}
func NewApprHandlerWithKeys(apprService d.ApprService, modelType reflect.Type, logError func(context.Context, string), idNames []string, option ...int) *ApprHandler {
	offset := 1
	if len(option) > 0 && option[0] >= 0 {
		offset = option[0]
	}
	return NewApprHandlerWithKeysAndLog(apprService, idNames, modelType, offset, logError, nil)
}
func NewApprHandlerWithKeysAndLog(apprService d.ApprService, keys []string, modelType reflect.Type, offset int, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprHandler {
	if offset < 0 {
		offset = 1
	}
	if keys == nil || len(keys) == 0 {
		keys = d.GetJsonPrimaryKeys(modelType)
	}
	indexes := d.GetIndexes(modelType)
	var resource, action1, action2 string
	if len(options) > 0 && len(options[0]) > 0 {
		action1 = options[0]
	} else {
		action1 = "approve"
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action2 = options[1]
	} else {
		action2 = "reject"
	}
	if len(options) > 2 && len(options[2]) > 0 {
		resource = options[2]
	} else {
		resource = d.BuildResourceName(modelType.Name())
	}
	return &ApprHandler{Log: writeLog, ApprService: apprService, ModelType: modelType, Keys: keys, Indexes: indexes, Offset: offset, Error: logError, Resource: resource, Action1: action1, Action2: action2}
}

func (c *ApprHandler) Approve(ctx *gin.Context) {
	r := ctx.Request
	id, er1 := d.BuildId(r, c.ModelType, c.Keys, c.Indexes, c.Offset)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	} else {
		result, er2 := c.ApprService.Approve(r.Context(), id)
		if er2 != nil {
			handleError(ctx, http.StatusOK, internalServerError, c.Error, c.Resource, c.Action1, er2, c.Log)
		} else {
			succeed(ctx, http.StatusOK, result, c.Log, c.Resource, c.Action1)
		}
	}
}

func (c *ApprHandler) Reject(ctx *gin.Context) {
	r := ctx.Request
	id, er1 := d.BuildId(r, c.ModelType, c.Keys, c.Indexes, c.Offset)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	} else {
		result, er2 := c.ApprService.Reject(r.Context(), id)
		if er2 != nil {
			handleError(ctx, http.StatusOK, internalServerError, c.Error, c.Resource, c.Action2, er2, c.Log)
		} else {
			succeed(ctx, http.StatusOK, result, c.Log, c.Resource, c.Action2)
		}
	}
}

func respond(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request.Context(), resource, action, success, desc)
	}
}
func handleError(ctx *gin.Context, code int, result interface{}, logError func(context.Context, string), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(ctx.Request.Context(), err.Error())
	}
	respond(ctx, code, result, writeLog, resource, action, false, err.Error())
}
func succeed(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	respond(ctx, code, result, writeLog, resource, action, true, "")
}
