package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

func NewMongoViewRepository(db *mongo.Database, collectionName string, modelType reflect.Type, idObjectId bool, options ...func(context.Context, interface{}) (interface{}, error)) *Loader {
	return NewMongoLoader(db, collectionName, modelType, idObjectId, options...)
}

func NewViewRepository(db *mongo.Database, collectionName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) *Loader {
	return NewMongoLoader(db, collectionName, modelType, false, options...)
}
func NewMongoRepositoryWithVersion(db *mongo.Database, collectionName string, modelType reflect.Type, idObjectId bool, versionField string, options ...Mapper) *Writer {
	return NewMongoWriterWithVersion(db, collectionName, modelType, idObjectId, versionField, options...)
}
func NewRepositoryWithVersion(db *mongo.Database, collectionName string, modelType reflect.Type, versionField string, options ...Mapper) *Writer {
	return NewMongoRepositoryWithVersion(db, collectionName, modelType, false, versionField, options...)
}
func NewRepository(db *mongo.Database, collectionName string, modelType reflect.Type, options ...Mapper) *Writer {
	return NewMongoRepositoryWithVersion(db, collectionName, modelType, false, "", options...)
}
func NewMongoViewAdapter(db *mongo.Database, collectionName string, modelType reflect.Type, idObjectId bool, options ...func(context.Context, interface{}) (interface{}, error)) *Loader {
	return NewMongoLoader(db, collectionName, modelType, idObjectId, options...)
}

func NewViewAdapter(db *mongo.Database, collectionName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) *Loader {
	return NewMongoLoader(db, collectionName, modelType, false, options...)
}
func NewMongoAdapterWithVersion(db *mongo.Database, collectionName string, modelType reflect.Type, idObjectId bool, versionField string, options ...Mapper) *Writer {
	return NewMongoWriterWithVersion(db, collectionName, modelType, idObjectId, versionField, options...)
}
func NewAdapterWithVersion(db *mongo.Database, collectionName string, modelType reflect.Type, versionField string, options ...Mapper) *Writer {
	return NewMongoAdapterWithVersion(db, collectionName, modelType, false, versionField, options...)
}
func NewAdapter(db *mongo.Database, collectionName string, modelType reflect.Type, options ...Mapper) *Writer {
	return NewMongoAdapterWithVersion(db, collectionName, modelType, false, "", options...)
}
