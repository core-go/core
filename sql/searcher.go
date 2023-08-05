package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

type Searcher struct {
	search  func(ctx context.Context, searchModel interface{}, results interface{}, limit int64, offset int64) (int64, error)
	ToArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}
func NewSearcher(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), options... func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) *Searcher {
	var toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	if len(options) > 0 {
		toArray = options[0]
	}
	return &Searcher{search: search, ToArray: toArray}
}

func (s *Searcher) Search(ctx context.Context, m interface{}, results interface{}, limit int64, offset int64) (int64, error) {
	return s.search(ctx, m, results, limit, offset)
}
func NewSearcherWithQuery(db *sql.DB, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, error) {
	return NewSearcherWithArray(db, modelType, buildQuery, nil, options...)
}
func NewSearcherWithArray(db *sql.DB, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, error) {
	builder, err := NewSearchBuilderWithArray(db, modelType, buildQuery, toArray, options...)
	if err != nil {
		return nil, err
	}
	return NewSearcher(builder.Search, toArray), nil
}
