package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

type ApprHandler struct {
	WriteLog    func(ctx context.Context, resource string, action string, success bool, desc string) error
	ApprService ApprService
	ModelType   reflect.Type
	IdNames     []string
	Indexes     map[string]int
	LogError    func(context.Context, string)
	Offset      int
	Resource    string
}
func NewApprHandler(apprService ApprService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, option ...int) *ApprHandler {
	offset := 1
	if len(option) == 1 {
		offset = option[0]
	}
	idNames := GetListFieldsTagJson(modelType)
	indexs := GetIndexes(modelType)
	return &ApprHandler{WriteLog: writeLog, ApprService: apprService, ModelType: modelType, IdNames: idNames, Indexes: indexs, Offset: offset, Resource: resource, LogError: logError}
}
func NewApprHandlerWithIds(apprService ApprService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, idNames []string, option ...int) *ApprHandler {
	offset := 1
	if len(option) == 1 {
		offset = option[0]
	}
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	indexs := GetIndexes(modelType)
	return &ApprHandler{WriteLog: writeLog, ApprService: apprService, ModelType: modelType, IdNames: idNames, Indexes: indexs, Offset: offset, Resource: resource, LogError: logError}
}

func (c *ApprHandler) newModel(body interface{}) (out interface{}) {
	req := reflect.New(c.ModelType).Interface()
	if body != nil {
		switch s := body.(type) {
		case io.Reader:
			err := json.NewDecoder(s).Decode(&req)
			if err != nil {
				return err
			}
			return req
		}
	}
	return req
}

func (c *ApprHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := BuildId(r, c.ModelType, c.IdNames, c.Indexes, c.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprService.Approve(r.Context(), id)
		if err != nil {
			Error(w, r, http.StatusOK, StatusError, c.LogError, c.Resource, "approve", err, c.WriteLog)
		} else {
			Succeed(w, r, http.StatusOK, result, c.WriteLog, c.Resource, "approve")
		}
	}
}

func (c *ApprHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := BuildId(r, c.ModelType, c.IdNames, c.Indexes, c.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprService.Reject(r.Context(), id)
		if err != nil {
			Error(w, r, http.StatusOK, StatusError, c.LogError, c.Resource, "reject", err, c.WriteLog)
		} else {
			Succeed(w, r, http.StatusOK, result, c.WriteLog, c.Resource, "reject")
		}
	}
}
