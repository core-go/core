package search

import "context"

type SearchService interface {
	Search(ctx context.Context, filter interface{}, results interface{}, limit int64, options ...int64) (int64, string, error)
}
