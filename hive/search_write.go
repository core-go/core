package hive

import (
	hv "github.com/beltran/gohive"
	"reflect"
)

func NewSearchWriterWithVersionAndMap(connection *hv.Connection, tableName string, modelType reflect.Type, buildQuery func(interface{}) string, versionField string, mapper Mapper) (*Searcher, *Writer, error) {
	if mapper == nil {
		searcher, er0 := NewSearcherWithQuery(connection, modelType, buildQuery)
		if er0 != nil {
			return searcher, nil, er0
		}
		writer, er1 := NewWriterWithVersion(connection, tableName, modelType, versionField, mapper)
		return searcher, writer, er1
	} else {
		searcher, er0 := NewSearcherWithQuery(connection, modelType, buildQuery, mapper.DbToModel)
		if er0 != nil {
			return searcher, nil, er0
		}
		writer, er1 := NewWriterWithVersion(connection, tableName, modelType, versionField, mapper)
		return searcher, writer, er1
	}
}
func NewSearchWriterWithVersion(connection *hv.Connection, tableName string, modelType reflect.Type, buildQuery func(interface{}) string, versionField string, options...Mapper) (*Searcher, *Writer, error) {
	var mapper Mapper
	if len(options) > 0 {
		mapper = options[0]
	}
	return NewSearchWriterWithVersionAndMap(connection, tableName, modelType, buildQuery, versionField, mapper)
}
func NewSearchWriterWithMap(connection *hv.Connection, tableName string, modelType reflect.Type, buildQuery func(interface{}) string, mapper Mapper, options...string) (*Searcher, *Writer, error) {
	var versionField string
	if len(options) > 0 {
		versionField = options[0]
	}
	return NewSearchWriterWithVersionAndMap(connection, tableName, modelType, buildQuery, versionField, mapper)
}
func NewSearchWriter(connection *hv.Connection, tableName string, modelType reflect.Type, buildQuery func(interface{}) string, options...Mapper) (*Searcher, *Writer, error) {
	var mapper Mapper
	if len(options) > 0 {
		mapper = options[0]
	}
	return NewSearchWriterWithVersionAndMap(connection, tableName, modelType, buildQuery, "", mapper)
}
