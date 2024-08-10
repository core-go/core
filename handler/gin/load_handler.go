package gin

import (
	"context"
	sv "github.com/core-go/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Return(ctx *gin.Context, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	if err != nil {
		RespondAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, err, logError, writeLog, resource, action)
	} else {
		if sv.IsNil(model) {
			ReturnAndLog(ctx, http.StatusNotFound, nil, writeLog, false, resource, action, "Not found")
		} else {
			Succeed(ctx, http.StatusOK, model, writeLog, resource, action)
		}
	}
}
func RespondAndLog(ctx *gin.Context, code int, result interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	}
	if err != nil {
		if logError != nil {
			logError(ctx.Request.Context(), err.Error())
			ReturnAndLog(ctx, http.StatusInternalServerError, sv.InternalServerError, writeLog, false, resource, action, err.Error())
		} else {
			ReturnAndLog(ctx, http.StatusInternalServerError, err.Error(), writeLog, false, resource, action, err.Error())
		}
	} else {
		ReturnAndLog(ctx, code, result, writeLog, true, resource, action, "")
	}
}
func ReturnAndLog(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, success bool, resource string, action string, desc string) {
	ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request.Context(), resource, action, success, desc)
	}
}
func ErrorAndLog(ctx *gin.Context, code int, result interface{}, logError func(context.Context, string, ...map[string]interface{}), resource string, action string, err error, writeLog func(context.Context, string, string, bool, string) error) {
	if logError != nil {
		logError(ctx.Request.Context(), err.Error())
	}
	ReturnAndLog(ctx, code, result, writeLog, false, resource, action, err.Error())
}
func Succeed(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	ReturnAndLog(ctx, code, result, writeLog, true, resource, action, "")
}
