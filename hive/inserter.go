package hive

import (
	"context"
	hv "github.com/beltran/gohive"
	"reflect"
)

type Inserter struct {
	connection   *hv.Connection
	tableName    string
	Map          func(ctx context.Context, model interface{}) (interface{}, error)
	schema       *Schema
	VersionIndex int
}

func NewInserterWithMap(connection *hv.Connection, tableName string, modelType reflect.Type, mp func(context.Context, interface{}) (interface{}, error), options ...int) *Inserter {
	versionIndex := -1
	if len(options) > 0 && options[0] >= 0 {
		versionIndex = options[0]
	}
	schema := CreateSchema(modelType)
	return &Inserter{connection: connection, tableName: tableName, Map: mp, schema: schema, VersionIndex: versionIndex}
}

func NewInserter(db *hv.Connection, tableName string, modelType reflect.Type, options ...func(ctx context.Context, model interface{}) (interface{}, error)) *Inserter {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	return NewInserterWithMap(db, tableName, modelType, mp)
}

func (w *Inserter) Write(ctx context.Context, model interface{}) error {
	if w.Map != nil {
		m2, er0 := w.Map(ctx, model)
		if er0 != nil {
			return er0
		}
		stm := BuildToInsertWithVersion(w.tableName, m2, w.VersionIndex, false, w.schema)
		cursor := w.connection.Cursor()
		cursor.Exec(ctx, stm)
		return cursor.Err
	}
	stm := BuildToInsertWithVersion(w.tableName, model, w.VersionIndex, false, w.schema)
	cursor := w.connection.Cursor()
	cursor.Exec(ctx, stm)
	return cursor.Err
}
