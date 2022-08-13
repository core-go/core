package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

func NewViewRepository(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	return NewLoaderWithArray(db, tableName, modelType, nil, options...)
}
func NewViewRepositoryWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	return NewLoaderWithArray(db, tableName, modelType, toArray, options...)
}
func NewSqlViewRepository(db *sql.DB, tableName string, modelType reflect.Type, mp func(context.Context, interface{}) (interface{}, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Loader, error) {
	return NewSqlLoader(db, tableName, modelType, mp, toArray, options...)
}

func NewRepository(db *sql.DB, tableName string, modelType reflect.Type, options ...Mapper) (*Writer, error) {
	return NewWriter(db, tableName, modelType, options...)
}
func NewRepositoryWithVersion(db *sql.DB, tableName string, modelType reflect.Type, versionField string, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersion(db, tableName, modelType, versionField, options...)
}
func NewRepositoryWithVersionAndArray(db *sql.DB, tableName string, modelType reflect.Type, versionField string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersionAndArray(db, tableName, modelType, versionField, toArray, options...)
}
func NewSqlRepositoryWithVersion(db *sql.DB, tableName string, modelType reflect.Type, versionField string, mapper Mapper, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Writer, error) {
	return NewSqlWriterWithVersion(db, tableName, modelType, versionField, mapper, toArray, options...)
}
func NewRepositoryWithMap(db *sql.DB, tableName string, modelType reflect.Type, mapper Mapper, options ...func(i int) string) (*Writer, error) {
	return NewWriterWithMap(db, tableName, modelType, mapper, options...)
}
func NewRepositoryWithMapAndArray(db *sql.DB, tableName string, modelType reflect.Type, mapper Mapper, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Writer, error) {
	return NewSqlWriterWithVersion(db, tableName, modelType, "", mapper, toArray, options...)
}
func NewRepositoryWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...Mapper) (*Writer, error) {
	return NewWriterWithArray(db, tableName, modelType, toArray, options...)
}
