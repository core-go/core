package echo

import (
	"context"
	"github.com/core-go/core"
	"github.com/labstack/echo/v4"
)

func Decode(ctx echo.Context, obj interface{}, options ...func(context.Context, interface{}) error) error {
	return core.Decode(ctx.Response().Writer, ctx.Request(), obj, options...)
}
func DecodeAndCheckId(ctx echo.Context, obj interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) error) error {
	return core.DecodeAndCheckId(ctx.Response().Writer, ctx.Request(), obj, keysJson, mapIndex, options...)
}
func BuildMapAndCheckId(ctx echo.Context, obj interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) error) (map[string]interface{}, error) {
	r, j, err := core.BuildMapAndCheckId(ctx.Response().Writer, ctx.Request(), obj, keysJson, mapIndex, options...)
	ctx.SetRequest(r)
	return j, err
}
func AfterDeletedWithLog(ctx echo.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterDeletedWithLog(ctx.Response().Writer, ctx.Request(), count, err, logError, writeLog, opts...)
}
func AfterDeleted(ctx echo.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterDeleted(ctx.Response().Writer, ctx.Request(), count, err, logError)
}
func HasError(ctx echo.Context, errors []core.ErrorMessage, err error, logError func(context.Context, string, ...map[string]interface{}), model interface{}, writeLog func(context.Context, string, string, bool, string) error, opts ...string) bool {
	return core.HasError(ctx.Response().Writer, ctx.Request(), errors, err, logError, model, writeLog, opts...)
}
func AfterSavedWithLog(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterSavedWithLog(ctx.Response().Writer, ctx.Request(), body, count, err, logError, writeLog, opts...)
}
func AfterSaved(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterSaved(ctx.Response().Writer, ctx.Request(), body, count, err, logError)
}
func AfterCreatedWithLog(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, opts ...string) error {
	return core.AfterCreatedWithLog(ctx.Response().Writer, ctx.Request(), body, count, err, logError, writeLog, opts...)
}
func AfterCreated(ctx echo.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterCreated(ctx.Response().Writer, ctx.Request(), body, count, err, logError)
}
