package gin

import (
	"context"
	"github.com/core-go/core"
	"github.com/gin-gonic/gin"
)

func Decode(ctx *gin.Context, obj interface{}, options ...func(context.Context, interface{}) error) error {
	return core.Decode(ctx.Writer, ctx.Request, obj, options...)
}
func DecodeAndCheckId(ctx *gin.Context, obj interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) error) error {
	return core.DecodeAndCheckId(ctx.Writer, ctx.Request, obj, keysJson, mapIndex, options...)
}
func BuildMapAndCheckId(ctx *gin.Context, obj interface{}, keysJson []string, mapIndex map[string]int, options ...func(context.Context, interface{}) error) (map[string]interface{}, error) {
	r, j, err := core.BuildMapAndCheckId(ctx.Writer, ctx.Request, obj, keysJson, mapIndex, options...)
	ctx.Request = r
	return j, err
}
func AfterDeletedWithLog(ctx *gin.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return core.AfterDeletedWithLog(ctx.Writer, ctx.Request, count, err, logError, writeLog, resource, action)
}
func AfterDeleted(ctx *gin.Context, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterDeleted(ctx.Writer, ctx.Request, count, err, logError)
}
func HasError(ctx *gin.Context, errors []core.ErrorMessage, err error, logError func(context.Context, string, ...map[string]interface{}), model interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string) bool {
	return core.HasError(ctx.Writer, ctx.Request, errors, err, logError, model, writeLog, resource, action)
}
func AfterSavedWithLog(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return core.AfterSavedWithLog(ctx.Writer, ctx.Request, body, count, err, logError, writeLog, resource, action)
}
func AfterSaved(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterSaved(ctx.Writer, ctx.Request, body, count, err, logError)
}
func AfterCreatedWithLog(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) error {
	return core.AfterCreatedWithLog(ctx.Writer, ctx.Request, body, count, err, logError, writeLog, resource, action)
}
func AfterCreated(ctx *gin.Context, body interface{}, count int64, err error, logError func(context.Context, string, ...map[string]interface{})) error {
	return core.AfterCreated(ctx.Writer, ctx.Request, body, count, err, logError)
}
