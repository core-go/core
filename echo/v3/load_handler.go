package echo

import (
	"context"
	"github.com/core-go/core"
	"github.com/labstack/echo"
)

func ReturnWithLog(ctx echo.Context, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return core.ReturnWithLog(ctx.Response().Writer, ctx.Request(), model, err, logError, writeLog, resource, action)
}
func Return(ctx echo.Context, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.Return(ctx.Response().Writer, ctx.Request(), model, err, logError)
}
