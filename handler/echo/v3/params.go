package echo

import (
	"github.com/core-go/core"
	"github.com/labstack/echo"
)

func GetRequiredString(c echo.Context, opts ...int) (string, error) {
	return core.GetRequiredString(c.Response().Writer, c.Request(), opts...)
}
func GetRequiredInt64(c echo.Context, opts ...int) (int64, error) {
	return core.GetRequiredInt64(c.Response().Writer, c.Request(), opts...)
}
func GetRequiredUint64(c echo.Context, opts ...int) (uint64, error) {
	return core.GetRequiredUint64(c.Response().Writer, c.Request(), opts...)
}
func GetRequiredInt(c echo.Context, opts ...int) (int, error) {
	return core.GetRequiredInt(c.Response().Writer, c.Request(), opts...)
}
func GetRequiredInt32(c echo.Context, opts ...int) (int32, error) {
	return core.GetRequiredInt32(c.Response().Writer, c.Request(), opts...)
}
