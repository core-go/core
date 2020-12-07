package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

type ApprHandler struct {
	LogWriter   LogWriter
	ApprService ApprService
	ModelType   reflect.Type
	IdNames     []string
	Indexes     map[string]int
	LogError    func(context.Context, string)
	Offset      int
	Resource    string
}

func NewApprHandler(apprService ApprService, modelType reflect.Type, logWriter LogWriter, idNames []string, resource string, logError func(context.Context, string), option ...int) *ApprHandler {
	offset := 1
	if len(option) == 1 {
		offset = option[0]
	}
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	indexs := GetIndexes(modelType)
	return &ApprHandler{LogWriter: logWriter, ApprService: apprService, ModelType: modelType, IdNames: idNames, Indexes: indexs, Offset: offset, Resource: resource, LogError: logError}
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
			Error(w, r, http.StatusOK, StatusError, c.LogError, c.Resource, "approve", err, c.LogWriter)
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "approve")
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
			Error(w, r, http.StatusOK, StatusError, c.LogError, c.Resource, "reject", err, c.LogWriter)
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "reject")
		}
	}
}
