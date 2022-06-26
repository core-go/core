package gin

import (
	"context"
	d "github.com/core-go/core/diff"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type ApprListHandler struct {
	ApprListService d.ApprListService
	Keys            []string
	ModelType       reflect.Type
	Error           func(context.Context, string)
	Log             func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource        string
	Action1         string
	Action2         string
}

func NewApprListHandler(apprListService d.ApprListService, modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprListHandler {
	return NewApprListHandlerWithKeys(apprListService, nil, modelType, logError, writeLog, options...)
}

func NewApprListHandlerWithKeys(apprListService d.ApprListService, keys []string, modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprListHandler {
	if keys == nil || len(keys) == 0 {
		keys = d.GetJsonPrimaryKeys(modelType)
	}
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
	return &ApprListHandler{ApprListService: apprListService, ModelType: modelType, Keys: keys, Resource: resource, Error: logError, Log: writeLog, Action1: action1, Action2: action2}
}

func (c *ApprListHandler) Approve(ctx *gin.Context) {
	r := ctx.Request
	ids, er1 := d.BuildIds(r, c.ModelType, c.Keys)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	} else {
		result, er2 := c.ApprListService.Approve(r.Context(), ids)
		if er2 != nil {
			handleError(ctx, http.StatusInternalServerError, internalServerError, c.Error, c.Resource, c.Action1, er2, c.Log)
		} else {
			succeed(ctx, http.StatusOK, result, c.Log, c.Resource, c.Action1)
		}
	}
}

func (c *ApprListHandler) Reject(ctx *gin.Context) {
	r := ctx.Request
	ids, er1 := d.BuildIds(r, c.ModelType, c.Keys)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
	} else {
		result, er2 := c.ApprListService.Reject(r.Context(), ids)
		if er2 != nil {
			handleError(ctx, http.StatusInternalServerError, internalServerError, c.Error, c.Resource, c.Action2, er2, c.Log)
		} else {
			succeed(ctx, http.StatusOK, result, c.Log, c.Resource, c.Action2)
		}
	}
}
