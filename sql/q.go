package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

func NewQuery(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	return NewLoaderWithArray(db, tableName, modelType, nil, options...)
}
func NewQueryWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	var mp func(ctx context.Context, model interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	return NewSqlLoader(db, tableName, modelType, mp, toArray)
}
func NewSqlQuery(db *sql.DB, tableName string, modelType reflect.Type, mp func(context.Context, interface{}) (interface{}, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Loader, error) {
	return NewSqlLoader(db, tableName, modelType, mp, toArray, options...)
}
