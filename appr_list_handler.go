package service

import (
	"context"
	"net/http"
	"reflect"
)

type ApprListHandler struct {
	WriteLog        func(ctx context.Context, resource string, action string, success bool, desc string) error
	ApprListService ApprListService
	ModelType       reflect.Type
	IdNames         []string
	LogError        func(context.Context, string)
	Resource        string
}
func NewApprListHandler(apprListService ApprListService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error) *ApprListHandler {
	idNames := GetListFieldsTagJson(modelType)
	return &ApprListHandler{WriteLog: writeLog, ApprListService: apprListService, ModelType: modelType, IdNames: idNames, Resource: resource, LogError: logError}
}

func NewApprListHandlerWithIds(apprListService ApprListService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, idNames []string) *ApprListHandler {
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	return &ApprListHandler{WriteLog: writeLog, ApprListService: apprListService, ModelType: modelType, IdNames: idNames, Resource: resource, LogError: logError}
}

func (c *ApprListHandler) Approve(w http.ResponseWriter, r *http.Request) {
	ids, err := BuildIds(r, c.ModelType, c.IdNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprListService.Approve(r.Context(), ids)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, InternalServerError, c.LogError, c.Resource, "approve", err, c.WriteLog)
		} else {
			Succeed(w, r, http.StatusOK, result, c.WriteLog, c.Resource, "approve")
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
			Error(w, r, http.StatusInternalServerError, InternalServerError, c.LogError, c.Resource, "reject", err, c.WriteLog)
		} else {
			Succeed(w, r, http.StatusOK, result, c.WriteLog, c.Resource, "reject")
		}
	}
}
