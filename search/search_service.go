package search

import "context"

type SearchService interface {
	Search(ctx context.Context, filter interface{}, results interface{}, limit int64, offset int64) (int64, error)
}
type SearchQuery interface {
	Search(ctx context.Context, filter interface{}, results interface{}, limit int64, offset int64) (int64, error)
}
