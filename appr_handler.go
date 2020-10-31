package service

import (
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
	Indexs      map[string]int
	Offset      int
	Resource    string
}

func NewApprHandler(apprService ApprService, modelType reflect.Type, logWriter LogWriter, idNames []string, resource string, option ...int) *ApprHandler {
	offset := 1
	if len(option) == 1 {
		offset = option[0]
	}
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	indexs := GetIndexes(modelType)
	return &ApprHandler{LogWriter: logWriter, ApprService: apprService, ModelType: modelType, IdNames: idNames, Indexs: indexs, Offset: offset, Resource: resource}
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
	id, err := BuildId(r, c.ModelType, c.IdNames, c.Indexs, c.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprService.Approve(r.Context(), id)
		if err != nil {
			Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource,  "Approve", false, err.Error())
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "Approve")
		}
	}
}

func (c *ApprHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := BuildId(r, c.ModelType, c.IdNames, c.Indexs, c.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprService.Reject(r.Context(), id)
		if err != nil {
			Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource,  "Reject", false, err.Error())
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "Reject")
		}
	}
}
