package export

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
)

func NewExportAdapter[T any](db *sql.DB,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, *T) string,
	write func(p []byte) (n int, err error),
	close func() error,
	opts ...func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
) (*Exporter[T], error) {
	return NewExporter[T](db, buildQuery, transform, write, close, opts...)
}

func NewExportService[T any](db *sql.DB,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, *T) string,
	write func(p []byte) (n int, err error),
	close func() error,
	opts ...func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
) (*Exporter[T], error) {
	return NewExporter[T](db, buildQuery, transform, write, close, opts...)
}

func NewExporter[T any](db *sql.DB,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, *T) string,
	write func(p []byte) (n int, err error),
	close func() error,
	opts ...func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
) (*Exporter[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	fieldsIndex, err := GetColumnIndexes(modelType)
	if err != nil {
		return nil, err
	}
	columns := GetColumnsSelect(modelType)
	var toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	if len(opts) > 0 {
		toArray = opts[0]
	}
	return &Exporter[T]{DB: db, columns: columns, Map: fieldsIndex, BuildQuery: buildQuery, Transform: transform, Write: write, Close: close, Array: toArray}, nil
}

type Exporter[T any] struct {
	DB         *sql.DB
	Map        map[string]int
	columns    []string
	Transform  func(context.Context, *T) string
	BuildQuery func(context.Context) (string, []interface{})
	Write      func(p []byte) (n int, err error)
	Close      func() error
	Array      func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *Exporter[T]) Export(ctx context.Context) (int64, error) {
	query, p := s.BuildQuery(ctx)
	rows, err := s.DB.QueryContext(ctx, query, p...)
	if err != nil {
		return 0, err
	}
	return s.ScanAndWrite(ctx, rows)
}

func (s *Exporter[T]) ScanAndWrite(ctx context.Context, rows *sql.Rows) (int64, error) {
	defer s.Close()

	var i int64
	i = 0
	for rows.Next() {
		var obj T
		r, swapValues := StructScan(&obj, s.columns, s.Map, s.Array)
		if err := rows.Scan(r...); err != nil {
			return i, err
		}
		SwapValuesToBool(&obj, &swapValues)
		err1 := s.TransformAndWrite(ctx, s.Write, &obj)
		if err1 != nil {
			return i, err1
		}
		i = i + 1
	}
	return i, nil
}

func (s *Exporter[T]) TransformAndWrite(ctx context.Context, write func(p []byte) (n int, err error), model *T) error {
	line := s.Transform(ctx, model)
	_, er := write([]byte(line))
	return er
}
