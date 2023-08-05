package hive

import (
	"context"
	hv "github.com/beltran/gohive"
	"reflect"
)

type Searcher struct {
	search  func(ctx context.Context, searchModel interface{}, results interface{}, limit int64, offset int64) (int64, error)
}
func NewSearcher(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error)) *Searcher {
	return &Searcher{search: search}
}

func (s *Searcher) Search(ctx context.Context, m interface{}, results interface{}, limit int64, offset int64) (int64, error) {
	return s.search(ctx, m, results, limit, offset)
}
func NewSearcherWithQuery(db *hv.Connection, modelType reflect.Type, buildQuery func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, error) {
	builder, err := NewSearchBuilder(db, modelType, buildQuery, options...)
	if err != nil {
		return nil, err
	}
	return NewSearcher(builder.Search), nil
}
