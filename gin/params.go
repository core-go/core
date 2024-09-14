package gin

import (
	"github.com/core-go/core"
	"github.com/gin-gonic/gin"
)

func GetRequiredString(c *gin.Context, opts ...int) (string, error) {
	return core.GetRequiredString(c.Writer, c.Request, opts...)
}
func GetRequiredInt64(c *gin.Context, opts ...int) (int64, error) {
	return core.GetRequiredInt64(c.Writer, c.Request, opts...)
}
func GetRequiredUint64(c *gin.Context, opts ...int) (uint64, error) {
	return core.GetRequiredUint64(c.Writer, c.Request, opts...)
}
func GetRequiredInt(c *gin.Context, opts ...int) (int, error) {
	return core.GetRequiredInt(c.Writer, c.Request, opts...)
}
func GetRequiredInt32(c *gin.Context, opts ...int) (int32, error) {
	return core.GetRequiredInt32(c.Writer, c.Request, opts...)
}
