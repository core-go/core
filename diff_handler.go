package service

import (
	"net/http"
	"reflect"
)

type DiffHandler struct {
	LogWriter   LogWriter
	DiffService DiffService
	ModelType   reflect.Type
	IdNames     []string
	Indexes     map[string]int
	Offset      int
	Resource    string
}

func NewDiffHandler(diffService DiffService, modelType reflect.Type, idNames []string, resource string, logWriter LogWriter, option ...int) *DiffHandler {
	offset := 1
	if len(option) == 1 {
		offset = option[0]
	}
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	indexs := GetIndexes(modelType)
	return &DiffHandler{LogWriter: logWriter, DiffService: diffService, ModelType: modelType, IdNames: idNames, Indexes: indexs, Resource: resource, Offset: offset}
}

func (c *DiffHandler) Diff(w http.ResponseWriter, r *http.Request) {
	id, err := BuildId(r, c.ModelType, c.IdNames, c.Indexes, c.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.DiffService.Diff(r.Context(), id)
		if err != nil {
			Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource, "Diff", false, err.Error())
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "Diff")
		}
	}
}
