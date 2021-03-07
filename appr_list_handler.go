package service

import (
	"context"
	"net/http"
	"reflect"
)

type ApprListHandler struct {
	Error           func(context.Context, string)
	Log             func(ctx context.Context, resource string, action string, success bool, desc string) error
	ApprListService ApprListService
	ModelType       reflect.Type
	IdNames         []string
	Resource        string
}

func NewApprListHandler(apprListService ApprListService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error) *ApprListHandler {
	idNames := GetListFieldsTagJson(modelType)
	return &ApprListHandler{Log: writeLog, ApprListService: apprListService, ModelType: modelType, IdNames: idNames, Resource: resource, Error: logError}
}

func NewApprListHandlerWithIds(apprListService ApprListService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, idNames []string) *ApprListHandler {
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	return &ApprListHandler{Log: writeLog, ApprListService: apprListService, ModelType: modelType, IdNames: idNames, Resource: resource, Error: logError}
}

func (c *ApprListHandler) Approve(w http.ResponseWriter, r *http.Request) {
	ids, err := BuildIds(r, c.ModelType, c.IdNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprListService.Approve(r.Context(), ids)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, InternalServerError, c.Error, c.Resource, "approve", err, c.Log)
		} else {
			Succeed(w, r, http.StatusOK, result, c.Log, c.Resource, "approve")
		}
	}
}

func (c *ApprListHandler) Reject(w http.ResponseWriter, r *http.Request) {
	ids, err := BuildIds(r, c.ModelType, c.IdNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprListService.Reject(r.Context(), ids)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, InternalServerError, c.Error, c.Resource, "reject", err, c.Log)
		} else {
			Succeed(w, r, http.StatusOK, result, c.Log, c.Resource, "reject")
		}
	}
}
