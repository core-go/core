package echo

import (
	"context"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

const internalServerError = "Internal Server Error"

type QueryHandler struct {
	Load     func(ctx context.Context, key string, max int64) ([]string, error)
	LogError func(context.Context, string, ...map[string]interface{})
	Keyword  string
	Max      string
}

func NewQueryHandler(load func(ctx context.Context, key string, max int64) ([]string, error), logError func(context.Context, string, ...map[string]interface{}), opts ...string) *QueryHandler {
	keyword := "keyword"
	if len(opts) > 0 && len(opts[0]) > 0 {
		keyword = opts[0]
	}
	max := "max"
	if len(opts) > 1 && len(opts[1]) > 0 {
		max = opts[1]
	}
	return &QueryHandler{load, logError, keyword, max}
}
func (h *QueryHandler) Query(ctx echo.Context) error {
	ps := ctx.Request().URL.Query()
	keyword := ps.Get(h.Keyword)
	if len(keyword) == 0 {
		vs := make([]string, 0)
		return ctx.JSON(http.StatusOK, vs)
	} else {
		max := ps.Get(h.Max)
		i, err := strconv.ParseInt(max, 10, 64)
		if err != nil {
			i = 20
		}
		if i < 0 {
			i = 20
		}
		vs, err := h.Load(ctx.Request().Context(), keyword, i)
		if err != nil {
			h.LogError(ctx.Request().Context(), err.Error())
			return ctx.String(http.StatusInternalServerError, internalServerError)
		} else {
			return ctx.JSON(http.StatusOK, vs)
		}
	}
}
