package diff

import (
	"context"
	"net/http"
	"reflect"
)

type ApprListHandler struct {
	ApprListService ApprListService
	Keys            []string
	ModelType       reflect.Type
	Error           func(context.Context, string)
	Log             func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource        string
	Action1         string
	Action2         string
}

func NewApprListHandler(apprListService ApprListService, modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprListHandler {
	return NewApprListHandlerWithKeys(apprListService, nil, modelType, logError, writeLog, options...)
}

func NewApprListHandlerWithKeys(apprListService ApprListService, keys []string, modelType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprListHandler {
	if keys == nil || len(keys) == 0 {
		keys = GetJsonPrimaryKeys(modelType)
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
		resource = BuildResourceName(modelType.Name())
	}
	return &ApprListHandler{ApprListService: apprListService, ModelType: modelType, Keys: keys, Resource: resource, Error: logError, Log: writeLog, Action1: action1, Action2: action2}
}

func (c *ApprListHandler) Approve(w http.ResponseWriter, r *http.Request) {
	ids, er1 := BuildIds(r, c.ModelType, c.Keys)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	} else {
		result, er2 := c.ApprListService.Approve(r.Context(), ids)
		if er2 != nil {
			handleError(w, r, http.StatusInternalServerError, internalServerError, c.Error, c.Resource, c.Action1, er2, c.Log)
		} else {
			succeed(w, r, http.StatusOK, result, c.Log, c.Resource, c.Action1)
		}
	}
}

func (c *ApprListHandler) Reject(w http.ResponseWriter, r *http.Request) {
	ids, er1 := BuildIds(r, c.ModelType, c.Keys)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	} else {
		result, er2 := c.ApprListService.Reject(r.Context(), ids)
		if er2 != nil {
			handleError(w, r, http.StatusInternalServerError, internalServerError, c.Error, c.Resource, c.Action2, er2, c.Log)
		} else {
			succeed(w, r, http.StatusOK, result, c.Log, c.Resource, c.Action2)
		}
	}
}
