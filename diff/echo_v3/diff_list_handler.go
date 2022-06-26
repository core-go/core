package echo

import (
	"context"
	d "github.com/core-go/core/diff"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
)

type DiffListHandler struct {
	GetDiff     func(ctx context.Context, ids interface{}) (*[]d.DiffModel, error)
	Keys        []string
	ModelType   reflect.Type
	modelTypeId reflect.Type
	Error       func(context.Context, string)
	Log         func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource    string
	Action      string
	Config      *d.DiffModelConfig
}

func NewDiffListHandler(diff func(context.Context, interface{}) (*[]d.DiffModel, error), modelType reflect.Type, logError func(context.Context, string), config *d.DiffModelConfig, writeLog func(context.Context, string, string, bool, string) error) *DiffListHandler {
	return NewDiffListHandlerWithKeys(diff, nil, modelType, logError, config, writeLog)
}
func NewDiffListHandlerWithKeys(diff func(context.Context, interface{}) (*[]d.DiffModel, error), keys []string, modelType reflect.Type, logError func(context.Context, string), config *d.DiffModelConfig, writeLog func(context.Context, string, string, bool, string) error) *DiffListHandler {
	if keys == nil || len(keys) == 0 {
		keys = d.GetJsonPrimaryKeys(modelType)
	}
	modelTypeId := d.NewModelTypeID(modelType, keys)
	var resource, action string
	if config != nil {
		resource = config.Resource
		action = config.Action
	}
	if len(resource) == 0 {
		resource = d.BuildResourceName(modelType.Name())
	}
	if len(action) == 0 {
		action = "diff"
	}
	return &DiffListHandler{Log: writeLog, GetDiff: diff, ModelType: modelType, modelTypeId: modelTypeId, Keys: keys, Resource: resource, Action: action, Config: config, Error: logError}
}

func (c *DiffListHandler) DiffList(ctx echo.Context) error {
	r := ctx.Request()
	ids, er1 := d.BuildIds(r, c.modelTypeId, c.Keys)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return er1
	} else {
		list, er2 := c.GetDiff(r.Context(), ids)
		if er2 != nil {
			return handleError(ctx, http.StatusInternalServerError, internalServerError, c.Error, c.Resource, c.Action, er2, c.Log)
		} else {
			if c.Config == nil || list == nil || len(*list) == 0 {
				return succeed(ctx, http.StatusOK, list, c.Log, c.Resource, c.Action)
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
				return succeed(ctx, http.StatusOK, l, c.Log, c.Resource, c.Action)
			}
		}
	}
}
