package hive

import (
	hv "github.com/beltran/gohive"
	"reflect"
)

func NewAdapter(connection *hv.Connection, tableName string, modelType reflect.Type, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersion(connection, tableName, modelType, "", options...)
}
func NewAdapterWithVersion(connection *hv.Connection, tableName string, modelType reflect.Type, versionField string, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersion(connection, tableName, modelType, versionField, options...)
}
func NewRepository(connection *hv.Connection, tableName string, modelType reflect.Type, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersion(connection, tableName, modelType, "", options...)
}
func NewRepositoryWithVersion(connection *hv.Connection, tableName string, modelType reflect.Type, versionField string, options ...Mapper) (*Writer, error) {
	return NewWriterWithVersion(connection, tableName, modelType, versionField, options...)
}
