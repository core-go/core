package core

import (
	"context"
	"net/http"
)

type WriteLog func(context.Context, string, string, bool, string) error
type BuildParam func(int) string
type Log func(context.Context, string, ...map[string]interface{})
type Search func(ctx context.Context, filter interface{}, results interface{}, limit int64, offset int64) (int64, error)
type SearchFn func(ctx context.Context, filter interface{}, results interface{}, limit int64, nextPageToken string) (string, error)
type Generate func(context.Context) (string, error)
type Sequence func(context.Context, string) (int64, error)
type HandleFn func(w http.ResponseWriter, r *http.Request)
