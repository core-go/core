package echo

import (
	"context"
	d "github.com/core-go/core/diff"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
)

type DiffHandler struct {
	GetDiff   func(ctx context.Context, id interface{}) (*d.DiffModel, error)
	Keys      []string
	ModelType reflect.Type
	Error     func(context.Context, string)
	Indexes   map[string]int
	Offset    int
	Log       func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource  string
	Action    string
	Config    *d.DiffModelConfig
}

func NewDiffHandler(diff func(context.Context, interface{}) (*d.DiffModel, error), modelType reflect.Type, logError func(context.Context, string), config *d.DiffModelConfig, writeLog func(context.Context, string, string, bool, string) error, options ...int) *DiffHandler {
	return NewDiffHandlerWithKeys(diff, nil, modelType, logError, config, writeLog, options...)
}
func NewDiffHandlerWithKeys(diff func(context.Context, interface{}) (*d.DiffModel, error), keys []string, modelType reflect.Type, logError func(context.Context, string), config *d.DiffModelConfig, writeLog func(context.Context, string, string, bool, string) error, options ...int) *DiffHandler {
	offset := 1
	if len(options) > 0 {
		offset = options[0]
	}
	if keys == nil || len(keys) == 0 {
		keys = d.GetJsonPrimaryKeys(modelType)
	}
	indexes := d.GetIndexes(modelType)
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
	return &DiffHandler{Log: writeLog, GetDiff: diff, ModelType: modelType, Keys: keys, Indexes: indexes, Resource: resource, Offset: offset, Config: config, Error: logError}
}

func (c *DiffHandler) Diff(ctx echo.Context) error {
	r := ctx.Request()
	id, er1 := d.BuildId(r, c.ModelType, c.Keys, c.Indexes, c.Offset)
	if er1 != nil {
		ctx.String(http.StatusBadRequest, er1.Error())
		return er1
	} else {
		result, er2 := c.GetDiff(r.Context(), id)
		if er2 != nil {
			return handleError(ctx, http.StatusInternalServerError, internalServerError, c.Error, c.Resource, c.Action, er2, c.Log)
		} else {
			if c.Config == nil {
				return succeed(ctx, http.StatusOK, result, c.Log, c.Resource, c.Action)
			} else {
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
				return succeed(ctx, http.StatusOK, m, c.Log, c.Resource, c.Action)
			}
		}
	}
}
