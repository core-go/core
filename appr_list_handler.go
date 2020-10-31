package service

import (
	"net/http"
	"reflect"
)

type ApprListHandler struct {
	LogWriter       LogWriter
	ApprListService ApprListService
	ModelType       reflect.Type
	IdNames         []string
	Resource        string
}

func NewApprListHandler(apprListService ApprListService, modelType reflect.Type, logWriter LogWriter, idNames []string, resource string) *ApprListHandler {
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	return &ApprListHandler{LogWriter: logWriter, ApprListService: apprListService, ModelType: modelType, IdNames: idNames, Resource: resource}
}

func (c *ApprListHandler) Approve(w http.ResponseWriter, r *http.Request) {
	ids, err := BuildIds(r, c.ModelType, c.IdNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.ApprListService.Approve(r.Context(), ids)
		if err != nil {
			Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource, "Approve", false, err.Error())
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "Approve")
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
			Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource, "Reject", false, err.Error())
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "Reject")
		}
	}
}
