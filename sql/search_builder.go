package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

type SearchBuilder struct {
	Database    *sql.DB
	BuildQuery  func(sm interface{}) (string, []interface{})
	ModelType   reflect.Type
	Map         func(ctx context.Context, model interface{}) (interface{}, error)
	fieldsIndex map[string]int
	ToArray     func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}
func NewSearchBuilder(db *sql.DB, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), options ...func(context.Context, interface{}) (interface{}, error)) (*SearchBuilder, error) {
	return NewSearchBuilderWithArray(db, modelType, buildQuery, nil, options...)
}
func NewSearchBuilderWithArray(db *sql.DB, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*SearchBuilder, error) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	fieldsIndex, err := GetColumnIndexes(modelType)
	if err != nil {
		return nil, err
	}
	builder := &SearchBuilder{Database: db, fieldsIndex: fieldsIndex, BuildQuery: buildQuery, ModelType: modelType, Map: mp, ToArray: toArray}
	return builder, nil
}

func (b *SearchBuilder) Search(ctx context.Context, m interface{}, results interface{}, limit int64, offset int64) (int64, error) {
	sql, params := b.BuildQuery(m)
	total, er2 := BuildFromQuery(ctx, b.Database, b.fieldsIndex, results, sql, params, limit, offset, b.ToArray, b.Map)
	return total, er2
}
