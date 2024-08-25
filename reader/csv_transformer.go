package reader

import (
	"context"
	"reflect"
)

func NewCSVTransformer[T any]() (*CSVTransformer[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &CSVTransformer[T]{formatCols: formatCols}, nil
}

type CSVTransformer[T any] struct {
	formatCols map[int]Delimiter
}

func (f CSVTransformer[T]) Transform(ctx context.Context, record []string) (T, error) {
	var res T
	err := ScanLine(record, &res, f.formatCols)
	return res, err
}
