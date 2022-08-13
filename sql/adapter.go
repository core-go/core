package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

func NewViewAdapter(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	return NewLoaderWithArray(db, tableName, modelType, nil, options...)
}
func NewViewAdapterWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	return NewLoaderWithArray(db, tableName, modelType, toArray, options...)
}
func NewSqlViewAdapter(db *sql.DB, tableName string, modelType reflect.Type, mp func(context.Context, interface{}) (interface{}, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Loader, error) {
	return NewSqlLoader(db, tableName, modelType, mp, toArray, options...)
}

func NewAdapter(db *sql.DB, tableName string, modelType reflect.Type, options ...Mapper) (*Writer, error) {
	return NewWriter(db, tableName, modelType, options...)
}
func NewAdapterWithVersion(db *sql.DB, tableName string, modelType reflect.Type, versionField string, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersion(db, tableName, modelType, versionField, options...)
}
func NewAdapterWithVersionAndArray(db *sql.DB, tableName string, modelType reflect.Type, versionField string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersionAndArray(db, tableName, modelType, versionField, toArray, options...)
}
func NewSqlAdapterWithVersion(db *sql.DB, tableName string, modelType reflect.Type, versionField string, mapper Mapper, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Writer, error) {
	return NewSqlWriterWithVersion(db, tableName, modelType, versionField, mapper, toArray, options...)
}
func NewAdapterWithMap(db *sql.DB, tableName string, modelType reflect.Type, mapper Mapper, options ...func(i int) string) (*Writer, error) {
	return NewWriterWithMap(db, tableName, modelType, mapper, options...)
}
func NewAdapterWithMapAndArray(db *sql.DB, tableName string, modelType reflect.Type, mapper Mapper, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Writer, error) {
	return NewSqlWriterWithVersion(db, tableName, modelType, "", mapper, toArray, options...)
}
func NewAdapterWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...Mapper) (*Writer, error) {
	return NewWriterWithArray(db, tableName, modelType, toArray, options...)
}
