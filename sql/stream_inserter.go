package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

type StreamInserter struct {
	db           *sql.DB
	tableName    string
	BuildParam   func(i int) string
	Map          func(ctx context.Context, model interface{}) (interface{}, error)
	Driver       string
	schema       *Schema
	batchSize    int
	batch        []interface{}
	ToArray      func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func NewStreamInserter(db *sql.DB, tableName string, modelType reflect.Type, batchSize int, options ...func(context.Context, interface{}) (interface{}, error)) *StreamInserter {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}

	return NewSqlStreamInserter(db, tableName, modelType, batchSize, mp, nil)
}
func NewStreamInserterWithArray(db *sql.DB, tableName string, modelType reflect.Type, batchSize int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) *StreamInserter {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	return NewSqlStreamInserter(db, tableName, modelType, batchSize, mp, toArray)
}
func NewSqlStreamInserter(db *sql.DB, tableName string, modelType reflect.Type, batchSize int,
	mp func(context.Context, interface{}) (interface{}, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) *StreamInserter {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	driver := GetDriver(db)
	schema := CreateSchema(modelType)
	return &StreamInserter{db: db, Driver: driver, schema: schema, tableName: tableName, batchSize: batchSize, BuildParam: buildParam, Map: mp, ToArray: toArray}
}

func (w *StreamInserter) Write(ctx context.Context, model interface{}) error {
	if w.Map != nil {
		m2, er0 := w.Map(ctx, model)
		if er0 != nil {
			return er0
		}
		w.batch = append(w.batch, m2)
	} else {
		w.batch = append(w.batch, model)
	}
	if len(w.batch) >= w.batchSize {
		return w.Flush(ctx)
	}
	return nil
}

func (w *StreamInserter) Flush(ctx context.Context) error {
	// driver := GetDriver(w.db)
	query, args, err := BuildToInsertBatchWithSchema(w.tableName, w.batch, w.Driver, w.ToArray, w.BuildParam, w.schema)
	if err != nil {
		return err
	}
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		w.batch = make([]interface{}, 0)
	}()
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return tx.Commit()
}
