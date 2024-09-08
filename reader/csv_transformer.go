package reader

import (
	"context"
	"reflect"
)

func NewCSVTransformer(modelType reflect.Type) (*CSVTransformer, error) {
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &CSVTransformer{modelType: modelType, formatCols: formatCols}, nil
}

type CSVTransformer struct {
	modelType  reflect.Type
	formatCols map[int]Delimiter
}

func (f CSVTransformer) Transform(ctx context.Context, record []string, res interface{}) error {
	return ScanLine(record, res, f.formatCols)
}
