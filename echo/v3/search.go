package echo

import (
	"context"
	"github.com/core-go/core"
	"github.com/labstack/echo"
)

type ISearchHandler interface {
	Search(c echo.Context)
}

type SearchHandler[T any, K any] struct {
	*Handler[T, K]
	ISearchHandler
}

func NewSearchHandlerWithLog[T any, K any](
	searchHandler ISearchHandler,
	service core.Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	action *core.ActionConfig,
	writeLog func(context.Context, string, string, bool, string) error,
	opts ...core.Builder[T],
) *SearchHandler[T, K] {
	hdl := NewhandlerWithLog[T, K](service, logError, validate, action, writeLog, opts...)
	return &SearchHandler[T, K]{hdl, searchHandler}
}
func NewSearchHandler[T any, K any](
	searchHandler ISearchHandler,
	service core.Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	opts ...core.Builder[T],
) *SearchHandler[T, K] {
	hdl := NewhandlerWithLog[T, K](service, logError, validate, nil, nil, opts...)
	return &SearchHandler[T, K]{hdl, searchHandler}
}
