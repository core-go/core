package diff

import (
	"context"
	"net/http"
	"reflect"
)

type ApprHandler struct {
	ApprService ApprService
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

func NewApprHandler(apprService ApprService, modelType reflect.Type, logError func(context.Context, string), option ...int) *ApprHandler {
	offset := 1
	if len(option) > 0 && option[0] >= 0 {
		offset = option[0]
	}
	return NewApprHandlerWithKeysAndLog(apprService, nil, modelType, offset, logError, nil)
}
func NewApprHandlerWithLogs(apprService ApprService, modelType reflect.Type, offset int, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprHandler {
	return NewApprHandlerWithKeysAndLog(apprService, nil, modelType, offset, logError, writeLog, options...)
}
func NewApprHandlerWithKeys(apprService ApprService, modelType reflect.Type, logError func(context.Context, string), idNames []string, option ...int) *ApprHandler {
	offset := 1
	if len(option) > 0 && option[0] >= 0 {
		offset = option[0]
	}
	return NewApprHandlerWithKeysAndLog(apprService, idNames, modelType, offset, logError, nil)
}
func NewApprHandlerWithKeysAndLog(apprService ApprService, keys []string, modelType reflect.Type, offset int, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *ApprHandler {
	if offset < 0 {
		offset = 1
	}
	if keys == nil || len(keys) == 0 {
		keys = GetJsonPrimaryKeys(modelType)
	}
	indexes := GetIndexes(modelType)
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
	return &ApprHandler{Log: writeLog, ApprService: apprService, ModelType: modelType, Keys: keys, Indexes: indexes, Offset: offset, Error: logError, Resource: resource, Action1: action1, Action2: action2}
}

func (c *ApprHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, er1 := BuildId(r, c.ModelType, c.Keys, c.Indexes, c.Offset)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	} else {
		result, er2 := c.ApprService.Approve(r.Context(), id)
		if er2 != nil {
			handleError(w, r, http.StatusOK, internalServerError, c.Error, c.Resource, c.Action1, er2, c.Log)
		} else {
			succeed(w, r, http.StatusOK, result, c.Log, c.Resource, c.Action1)
		}
	}
}

func (c *ApprHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, er1 := BuildId(r, c.ModelType, c.Keys, c.Indexes, c.Offset)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	} else {
		result, er2 := c.ApprService.Reject(r.Context(), id)
		if er2 != nil {
			handleError(w, r, http.StatusOK, internalServerError, c.Error, c.Resource, c.Action2, er2, c.Log)
		} else {
			succeed(w, r, http.StatusOK, result, c.Log, c.Resource, c.Action2)
		}
	}
}
