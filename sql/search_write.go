package sql

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
)

func NewSearchWriterWithVersionAndMap(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), versionField string, mapper Mapper, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Searcher, *Writer, error) {
	if mapper == nil {
		searcher, er0 := NewSearcherWithArray(db, modelType, buildQuery, toArray)
		if er0 != nil {
			return searcher, nil, er0
		}
		writer, er1 := NewSqlWriterWithVersion(db, tableName, modelType, versionField, mapper, toArray, options...)
		return searcher, writer, er1
	} else {
		searcher, er0 := NewSearcherWithArray(db, modelType, buildQuery, toArray, mapper.DbToModel)
		if er0 != nil {
			return searcher, nil, er0
		}
		writer, er1 := NewSqlWriterWithVersion(db, tableName, modelType, versionField, mapper, toArray, options...)
		return searcher, writer, er1
	}
}
func NewSearchWriterWithVersion(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), versionField string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options...Mapper) (*Searcher, *Writer, error) {
	var mapper Mapper
	if len(options) > 0 {
		mapper = options[0]
	}
	return NewSearchWriterWithVersionAndMap(db, tableName, modelType, buildQuery, versionField, mapper, toArray)
}
func NewSearchWriterWithMap(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), mapper Mapper, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options...string) (*Searcher, *Writer, error) {
	var versionField string
	if len(options) > 0 {
		versionField = options[0]
	}
	return NewSearchWriterWithVersionAndMap(db, tableName, modelType, buildQuery, versionField, mapper, toArray)
}
func NewSearchWriter(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), options...Mapper) (*Searcher, *Writer, error) {
	return NewSearchWriterWithArray(db, tableName, modelType, buildQuery, nil, options...)
}
func NewSearchWriterWithArray(db *sql.DB, tableName string, modelType reflect.Type, buildQuery func(interface{}) (string, []interface{}), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options...Mapper) (*Searcher, *Writer, error) {
	build := GetBuild(db)
	var mapper Mapper
	if len(options) > 0 {
		mapper = options[0]
	}
	return NewSearchWriterWithVersionAndMap(db, tableName, modelType, buildQuery, "", mapper, toArray, build)
}
