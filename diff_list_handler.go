package service

import (
	"context"
	"net/http"
	"reflect"
	"strings"
)

type DiffListHandler struct {
	WriteLog        func(ctx context.Context, resource string, action string, success bool, desc string) error
	DiffListService DiffListService
	ModelType       reflect.Type
	modelTypeId     reflect.Type
	IdNames         []string
	Resource        string
	LogError        func(context.Context, string)
	Config          *DiffModelConfig
}
func NewDiffListHandler(diffListService DiffListService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, config *DiffModelConfig) *DiffListHandler {
	idNames := GetListFieldsTagJson(modelType)
	modelTypeId := newModelTypeID(modelType, idNames)
	return &DiffListHandler{WriteLog: writeLog, DiffListService: diffListService, ModelType: modelType, modelTypeId: modelTypeId, IdNames: idNames, Resource: resource, Config: config, LogError: logError}
}
func NewDiffListHandlerWithIds(diffListService DiffListService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, idNames []string, config *DiffModelConfig) *DiffListHandler {
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	modelTypeId := newModelTypeID(modelType, idNames)
	return &DiffListHandler{WriteLog: writeLog, DiffListService: diffListService, ModelType: modelType, modelTypeId: modelTypeId, IdNames: idNames, Resource: resource, Config: config, LogError: logError}
}

func (c *DiffListHandler) DiffList(w http.ResponseWriter, r *http.Request) {
	ids, err := BuildIds(r, c.modelTypeId, c.IdNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		list, err := c.DiffListService.Diff(r.Context(), ids)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, InternalServerError, c.LogError, c.Resource, "diff", err, c.WriteLog)
		} else {
			if c.Config == nil || list == nil || len(*list) == 0 {
				Succeed(w, r, http.StatusOK, list, c.WriteLog, c.Resource, "diff")
			} else {
				l := make([]map[string]interface{}, 0)
				for _, result := range *list {
					m := make(map[string]interface{})
					if result.Id != nil {
						m[c.Config.Id] = result.Id
					}
					if result.Origin != nil {
						m[c.Config.Origin] = result.Origin
					}
					if result.Value != nil {
						m[c.Config.Value] = result.Value
					}
					if len(result.By) > 0 {
						m[c.Config.By] = result.By
					}
					l = append(l, m)
				}
				Succeed(w, r, http.StatusOK, l, c.WriteLog, c.Resource, "diff")
			}
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
