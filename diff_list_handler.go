package service

import (
	"net/http"
	"reflect"
	"strings"
)

type DiffListHandler struct {
	LogWriter       LogWriter
	DiffListService DiffListService
	ModelType       reflect.Type
	modelTypeId     reflect.Type
	IdNames         []string
	Resource        string
}

func NewDiffListHandler(diffListService DiffListService, modelType reflect.Type, idNames []string, resource string, logWriter LogWriter) *DiffListHandler {
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	modelTypeId := newModelTypeID(modelType, idNames)
	return &DiffListHandler{LogWriter: logWriter, DiffListService: diffListService, ModelType: modelType, modelTypeId: modelTypeId, IdNames: idNames, Resource: resource}
}

func (c *DiffListHandler) DiffList(w http.ResponseWriter, r *http.Request) {
	ids, err := BuildIds(r, c.modelTypeId, c.IdNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.DiffListService.Diff(r.Context(), ids)
		if err != nil {
			Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource, "Diff", false, err.Error())
		} else {
			Succeed(w, r, http.StatusOK, result, c.LogWriter, c.Resource, "Diff")
		}
	}
}

func newModelTypeID(modelType reflect.Type, idJsonNames []string) reflect.Type {
	model := reflect.New(modelType).Interface()
	value := reflect.Indirect(reflect.ValueOf(model))
	sf := make([]reflect.StructField, 0)
	for i := 0; i < modelType.NumField(); i++ {
		sf = append(sf, modelType.Field(i))
		field := modelType.Field(i)
		json := field.Tag.Get("json")
		s := strings.Split(json, ",")[0]
		if Find(idJsonNames, s) == false {
			sf[i].Tag = `json:"-"`
		}
	}
	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)
	return reflect.TypeOf(newValue.Interface())
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
