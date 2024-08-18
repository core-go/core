package echo

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"

	"github.com/core-go/core"
)

func GetRequiredParam(c echo.Context, opts ...int) string {
	p := core.GetParam(c.Request(), opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return ""
	}
	return p
}
func GetRequiredInt(c echo.Context, opts ...int) (int, error) {
	p := core.GetParam(c.Request(), opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return 0, errors.New("parameter is required")
	}
	i, err := strconv.Atoi(p)
	if err != nil {
		c.String(http.StatusBadRequest, "parameter must be an integer")
	}
	return i, err
}
func GetRequiredInt64(c echo.Context, opts ...int) (int64, error) {
	p := core.GetParam(c.Request(), opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return 0, errors.New("parameter is required")
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "parameter must be an integer")
	}
	return i, err
}
func GetRequiredInt32(c echo.Context, opts ...int) (int32, error) {
	p := core.GetParam(c.Request(), opts...)
	if len(p) == 0 {
		c.String(http.StatusBadRequest, "parameter is required")
		return 0, errors.New("parameter is required")
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "parameter must be an integer")
	}
	return int32(i), err
}
