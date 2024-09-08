package handler

import (
	"context"
	"github.com/core-go/core"
	"net/http"
)

type ISearchHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
}

type SearchHandler[T any, K any] struct {
	*Handler[T, K]
	ISearchHandler
}

func NewSearchHandlerWithLog[T any, K any](
	searchHandler ISearchHandler,
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	action *core.ActionConfig,
	writeLog func(context.Context, string, string, bool, string) error,
	opts ...Builder[T],
) *SearchHandler[T, K] {
	hdl := NewhandlerWithLog[T, K](service, logError, validate, action, writeLog, opts...)
	return &SearchHandler[T, K]{hdl, searchHandler}
}
func NewSearchHandler[T any, K any](
	searchHandler ISearchHandler,
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]core.ErrorMessage, error),
	opts ...Builder[T],
) *SearchHandler[T, K] {
	hdl := NewhandlerWithLog[T, K](service, logError, validate, nil, nil, opts...)
	return &SearchHandler[T, K]{hdl, searchHandler}
}
