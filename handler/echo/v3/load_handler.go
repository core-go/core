package echo

import (
	"context"
	sv "github.com/core-go/core"
	"github.com/labstack/echo"
	"net/http"
)

func Return(ctx echo.Context, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	if err != nil {
		return RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	} else {
		if sv.IsNil(model) {
			return ReturnAndLog(ctx, http.StatusNotFound, nil, writeLog, false, resource, action, "Not found")
		} else {
			return Succeed(ctx, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
func RespondAndLog(ctx echo.Context, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) error {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		if logError != nil {
			logError(ctx.Request().Context(), err.Error())
			return ReturnAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, writeLog, false, resource, action, err.Error())
		} else {
			return ReturnAndLog(ctx, http.StatusInternalServerError, err.Error(), writeLog, false, resource, action, err.Error())
		}
	} else {
		return ReturnAndLog(ctx, code, result, writeLog, true, resource, action, "")
	}
}
func ReturnAndLog(ctx echo.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) error {
	err := ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request().Context(), resource, action, success, desc)
	}
	return err
}
func Succeed(ctx echo.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return ReturnAndLog(ctx, code, result, writeLog, true, resource, action, "")
}
