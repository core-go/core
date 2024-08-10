package gin

import (
	"net/http"
	"strconv"

	"github.com/core-go/core"
	"github.com/gin-gonic/gin"
)

func GetRequiredParam(c *gin.Context, opts ...int) string {
	p := core.GetParam(c.Request, opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return ""
	}
	return p
}
func GetRequiredInt(c *gin.Context, opts ...int) *int {
	p := core.GetParam(c.Request, opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return nil
	}
	i, err := strconv.Atoi(p)
	if err != nil {
		c.String(http.StatusBadRequest, "parameter must be an integer")
		return nil
	}
	return &i
}
func GetRequiredInt64(c *gin.Context, opts ...int) *int64 {
	p := core.GetParam(c.Request, opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return nil
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "parameter must be an integer")
		return nil
	}
	return &i
}
func GetRequiredInt32(c *gin.Context, opts ...int) *int32 {
	p := core.GetParam(c.Request, opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return nil
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "parameter must be an integer")
		return nil
	}
	j := int32(i)
	return &j
}
