package core

import (
	"context"
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
	validate func(context.Context, *T) ([]ErrorMessage, error),
	action *ActionConfig,
	writeLog func(context.Context, string, string, bool, string) error,
	opts ...Builder[T],
) *SearchHandler[T, K] {
	hdl := NewhandlerWithLog[T, K](service, logError, validate, writeLog, action, opts...)
	return &SearchHandler[T, K]{hdl, searchHandler}
}
func NewSearchHandler[T any, K any](
	searchHandler ISearchHandler,
	service Service[T, K],
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, *T) ([]ErrorMessage, error),
	opts ...Builder[T],
) *SearchHandler[T, K] {
	hdl := NewhandlerWithLog[T, K](service, logError, validate, nil, nil, opts...)
	return &SearchHandler[T, K]{hdl, searchHandler}
}
