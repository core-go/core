package hive

import (
	"context"
	hv "github.com/beltran/gohive"
	"reflect"
)

func NewSearchLoader(connection *hv.Connection, tableName string, modelType reflect.Type, buildQuery func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader, error) {
	return NewSqlSearchLoader(connection, tableName, modelType, buildQuery, options...)
}

func NewSqlSearchLoader(connection *hv.Connection, tableName string, modelType reflect.Type, buildQuery func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader, error) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	loader, er0 := NewLoader(connection, tableName, modelType, mp)
	if er0 != nil {
		return nil, loader, er0
	}
	searcher, er1 := NewSearcherWithQuery(connection, modelType, buildQuery, options...)
	return searcher, loader, er1
}
