package impt

import (
	"context"
	"database/sql"
	"io"
	"reflect"
)

func NewImportRepository(db *sql.DB, modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error) error) error,
) *Importer {
	return NewImporter(db, modelType, transform, write, read)
}
func NewImportAdapter(db *sql.DB, modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error) error) error,
) *Importer {
	return NewImporter(db, modelType, transform, write, read)
}
func NewImportService(db *sql.DB, modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error) error) error,
) *Importer {
	return NewImporter(db, modelType, transform, write, read)
}
func NewImporter(db *sql.DB, modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error) error) error,
) *Importer {
	return &Importer{DB: db, modelType: modelType, Transform: transform, Write: write, Read: read}
}

type Importer struct {
	DB        *sql.DB
	modelType reflect.Type
	Transform func(ctx context.Context, lines []string) (interface{}, error)
	Read      func(next func(lines []string, err error) error) error
	Write     func(ctx context.Context, data interface{}, endLineFlag bool) error
}

func (s *Importer) Import(ctx context.Context) (err error) {
	err = s.Read(func(lines []string, err error) error {
		if err == io.EOF {
			err = s.Write(ctx, nil, true)
			return nil
		}
		itemStruct, err := s.Transform(ctx, lines)
		if err != nil {
			return err
		}
		err = s.Write(ctx, itemStruct, false)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
