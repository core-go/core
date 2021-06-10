package diff

import (
	"context"
	"net/http"
	"reflect"
)

type DiffListHandler struct {
	GetDiff     func(ctx context.Context, ids interface{}) (*[]DiffModel, error)
	Keys        []string
	ModelType   reflect.Type
	modelTypeId reflect.Type
	Error       func(context.Context, string)
	Log         func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource    string
	Action      string
	Config      *DiffModelConfig
}

func NewDiffListHandler(diff func(context.Context, interface{}) (*[]DiffModel, error), modelType reflect.Type, logError func(context.Context, string), config *DiffModelConfig, writeLog func(context.Context, string, string, bool, string) error) *DiffListHandler {
	return NewDiffListHandlerWithKeys(diff, nil, modelType, logError, config, writeLog)
}
func NewDiffListHandlerWithKeys(diff func(context.Context, interface{}) (*[]DiffModel, error), keys []string, modelType reflect.Type, logError func(context.Context, string), config *DiffModelConfig, writeLog func(context.Context, string, string, bool, string) error) *DiffListHandler {
	if keys == nil || len(keys) == 0 {
		keys = GetJsonPrimaryKeys(modelType)
	}
	modelTypeId := NewModelTypeID(modelType, keys)
	var resource, action string
	if config != nil {
		resource = config.Resource
		action = config.Action
	}
	if len(resource) == 0 {
		resource = BuildResourceName(modelType.Name())
	}
	if len(action) == 0 {
		action = "diff"
	}
	return &DiffListHandler{Log: writeLog, GetDiff: diff, ModelType: modelType, modelTypeId: modelTypeId, Keys: keys, Resource: resource, Action: action, Config: config, Error: logError}
}

func (c *DiffListHandler) DiffList(w http.ResponseWriter, r *http.Request) {
	ids, er1 := BuildIds(r, c.modelTypeId, c.Keys)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
	} else {
		list, er2 := c.GetDiff(r.Context(), ids)
		if er2 != nil {
			handleError(w, r, http.StatusInternalServerError, internalServerError, c.Error, c.Resource, c.Action, er2, c.Log)
		} else {
			if c.Config == nil || list == nil || len(*list) == 0 {
				succeed(w, r, http.StatusOK, list, c.Log, c.Resource, c.Action)
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
				succeed(w, r, http.StatusOK, l, c.Log, c.Resource, c.Action)
			}
		}
	}
}
