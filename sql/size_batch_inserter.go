package sql

import (
	"context"
	"database/sql"
	"reflect"
)

type SizeBatchInserter struct {
	db         *sql.DB
	tableName  string
	BuildParam func(i int) string
	Map        func(ctx context.Context, model interface{}) (interface{}, error)
}
func NewSizeBatchInserter(db *sql.DB, tableName string, options...func(context.Context, interface{}) (interface{}, error)) *SizeBatchInserter {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) > 0 && options[0] != nil {
		mp = options[0]
	}
	return NewSizeSqlBatchInserter(db, tableName, mp)
}
func NewSizeSqlBatchInserter(db *sql.DB, tableName string, mp func(context.Context, interface{}) (interface{}, error), options...func(i int) string) *SizeBatchInserter {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	return &SizeBatchInserter{db: db, tableName: tableName, BuildParam: buildParam, Map: mp}
}

func (w *SizeBatchInserter) Write(ctx context.Context, models interface{}) ([]int, []int, error) {
	successIndices := make([]int, 0)
	failIndices := make([]int, 0)
	var models2 interface{}
	var er0 error
	if w.Map != nil {
		models2, er0 = MapModels(ctx, models, w.Map)
		if er0 != nil {
			s0 := reflect.ValueOf(models2)
			_, er0b := InterfaceSlice(models2)
			failIndices = ToArrayIndex(s0, failIndices)
			return successIndices, failIndices, er0b
		}
	} else {
		models2 = models
	}
	s := reflect.ValueOf(models2)
	_models, er1 := InterfaceSlice(models2)
	if er1 != nil {
		// Return full fail
		failIndices = ToArrayIndex(s, failIndices)
		return successIndices, failIndices, er1
	}
	_, er2 := InsertManyWithSize(ctx, w.db, w.tableName, _models, 0, w.BuildParam)

	if er2 == nil {
		// Return full success
		successIndices = ToArrayIndex(s, successIndices)
		return successIndices, failIndices, er2
	} else {
		// Return full fail
		failIndices = ToArrayIndex(s, failIndices)
	}
	return successIndices, failIndices, er2
}

func ToArrayIndex(value reflect.Value, indices []int) []int {
	for i := 0; i < value.Len(); i++ {
		indices = append(indices, i)
	}
	return indices
}
