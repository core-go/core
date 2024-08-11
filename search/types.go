package search

import "context"

type Search func(ctx context.Context, filter interface{}, results interface{}, limit int64, offset int64) (int64, error)
type SearchFn func(ctx context.Context, filter interface{}, results interface{}, limit int64, nextPageToken string) (string, error)
type WriteLog func(context.Context, string, string, bool, string) error
