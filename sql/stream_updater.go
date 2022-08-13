package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

type StreamUpdater struct {
	db           *sql.DB
	tableName    string
	BuildParam   func(i int) string
	Map          func(ctx context.Context, model interface{}) (interface{}, error)
	BoolSupport  bool
	VersionIndex int
	schema       *Schema
	batchSize    int
	batch        []interface{}
	ToArray      func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func NewStreamUpdater(db *sql.DB, tableName string, modelType reflect.Type, batchSize int, options ...func(context.Context, interface{}) (interface{}, error)) *StreamUpdater {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}

	return NewSqlStreamUpdater(db, tableName, modelType, batchSize, mp, nil)
}
func NewStreamUpdaterWithArray(db *sql.DB, tableName string, modelType reflect.Type, batchSize int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) *StreamUpdater {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	return NewSqlStreamUpdater(db, tableName, modelType, batchSize, mp, toArray)
}
func NewSqlStreamUpdater(db *sql.DB, tableName string, modelType reflect.Type, batchSize int,
	mp func(context.Context, interface{}) (interface{}, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) *StreamUpdater {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	driver := GetDriver(db)
	boolSupport := driver == DriverPostgres
	schema := CreateSchema(modelType)
	return &StreamUpdater{db: db, BoolSupport: boolSupport, VersionIndex: -1, schema: schema, tableName: tableName, batchSize: batchSize, BuildParam: buildParam, Map: mp, ToArray: toArray}
}

func (w *StreamUpdater) Write(ctx context.Context, model interface{}) error {
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

func (w *StreamUpdater) Flush(ctx context.Context) error {
	var queryArgsArray []Statement
	for _, v := range w.batch {
		query, args := BuildToUpdateWithArray(w.tableName, v, w.BuildParam, w.BoolSupport, w.ToArray, w.schema)
		queryArgs := Statement{
			Query: query,
			Params: args,
		}
		queryArgsArray = append(queryArgsArray, queryArgs)
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		w.batch = make([]interface{}, 0)
	}()

	for _, v := range queryArgsArray {
		_, err = tx.Exec(v.Query, v.Params...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
