package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

func NewExportRepository(db *mongo.Collection, modelType reflect.Type,
	buildQuery func(context.Context) bson.D,
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
	opts... func(context.Context) *options.FindOptions,
) *Exporter {
	return NewExporter(db, modelType, buildQuery, transform, write, close, opts...)
}
func NewExportAdapter(db *mongo.Collection, modelType reflect.Type,
	buildQuery func(context.Context) bson.D,
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
	opts... func(context.Context) *options.FindOptions,
) *Exporter {
	return NewExporter(db, modelType, buildQuery, transform, write, close, opts...)
}
func NewExportService(db *mongo.Collection, modelType reflect.Type,
	buildQuery func(context.Context) bson.D,
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
	opts... func(context.Context) *options.FindOptions,
) *Exporter {
	return NewExporter(db, modelType, buildQuery, transform, write, close, opts...)
}

func NewExporter(db *mongo.Collection, modelType reflect.Type,
	buildQuery func(context.Context) bson.D,
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
	opts... func(context.Context) *options.FindOptions,
) *Exporter {
	var opt func(context.Context) *options.FindOptions
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	return &Exporter{Collection: db, modelType: modelType, Write: write, Close: close, Transform: transform, BuildQuery: buildQuery, BuildFindOptions: opt}
}

type Exporter struct {
	Collection       *mongo.Collection
	modelType        reflect.Type
	Transform        func(context.Context, interface{}) string
	BuildQuery       func(context.Context) bson.D
	BuildFindOptions func(context.Context) *options.FindOptions
	Write            func(p []byte) (n int, err error)
	Close            func() error
}

func (s *Exporter) Export(ctx context.Context) error {
	query := s.BuildQuery(ctx)
	var cursor *mongo.Cursor
	var err error
	if s.BuildFindOptions != nil {
		optionsFind := s.BuildFindOptions(ctx)
		cursor, err = s.Collection.Find(ctx, query, optionsFind)
	} else {
		cursor, err = s.Collection.Find(ctx, query)
	}
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		initModel := reflect.New(s.modelType).Interface()
		err = cursor.Decode(initModel)
		if err != nil {
			return err
		}
		s.TransformAndWrite(ctx, s.Write, initModel)
	}
	return cursor.Err()
}

func (s *Exporter) TransformAndWrite(ctx context.Context, write func(p []byte) (n int, err error), model interface{}) error {
	line := s.Transform(ctx, model)
	_, er := write([]byte(line))
	return er
}
