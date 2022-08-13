package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)
func NewSearchLoader(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader, error) {
	return NewSearchLoaderWithArray(db, tableName, modelType, buildQuery, nil, options...)
}
func NewSearchLoaderWithArray(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader, error) {
	build := GetBuild(db)
	return NewSqlSearchLoader(db, tableName, modelType, buildQuery, build, toArray, options...)
}

func NewSqlSearchLoader(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), buildParam func(i int) string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader, error) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	loader, er0 := NewSqlLoader(db, tableName, modelType, mp, toArray, buildParam)
	if er0 != nil {
		return nil, loader, er0
	}
	searcher, er1 := NewSearcherWithArray(db, modelType, buildQuery, toArray, options...)
	return searcher, loader, er1
}
