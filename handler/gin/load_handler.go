package gin

import (
	"context"
	"github.com/core-go/core"
	"github.com/gin-gonic/gin"
)

func ReturnWithLog(ctx *gin.Context, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, resource string, action string) {
	core.ReturnWithLog(ctx.Writer, ctx.Request, model, err, logError, writeLog, resource, action)
}
func Return(ctx *gin.Context, model interface{}, err error, logError func(context.Context, string, ...map[string]interface{})) {
	core.Return(ctx.Writer, ctx.Request, model, err, logError)
}
