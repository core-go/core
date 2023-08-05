package hive

import (
	"context"
	hv "github.com/beltran/gohive"
	"reflect"
)

type Updater struct {
	connection   *hv.Connection
	tableName    string
	Map          func(ctx context.Context, model interface{}) (interface{}, error)
	VersionIndex int
	schema       *Schema
}

func NewUpdater(db *hv.Connection, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) *Updater {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	return NewUpdaterWithVersion(db, tableName, modelType, mp)
}
func NewUpdaterWithVersion(db *hv.Connection, tableName string, modelType reflect.Type, mp func(context.Context, interface{}) (interface{}, error), options ...int) *Updater {
	version := -1
	if len(options) > 0 && options[0] >= 0 {
		version = options[0]
	}
	schema := CreateSchema(modelType)
	return &Updater{connection: db, tableName: tableName, VersionIndex: version, schema: schema, Map: mp}
}

func (w *Updater) Write(ctx context.Context, model interface{}) error {
	if w.Map != nil {
		m2, er0 := w.Map(ctx, model)
		if er0 != nil {
			return er0
		}
		stm := BuildToUpdateWithVersion(w.tableName, m2, w.VersionIndex, w.schema)
		cursor := w.connection.Cursor()
		cursor.Exec(ctx, stm)
		return cursor.Err
	}
	stm := BuildToUpdateWithVersion(w.tableName, model, w.VersionIndex, w.schema)
	cursor := w.connection.Cursor()
	cursor.Exec(ctx, stm)
	return cursor.Err
}
