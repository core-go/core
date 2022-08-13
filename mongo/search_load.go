package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

func NewSearchLoaderWithQueryAndSort(db *mongo.Database, collection string, modelType reflect.Type, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, buildSort func(string, reflect.Type) bson.D, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader) {
	return NewMongoSearchLoaderWithQueryAndSort(db, collection, modelType, false, buildQuery, getSort, buildSort, options...)
}
func NewMongoSearchLoaderWithQuery(db *mongo.Database, collection string, modelType reflect.Type, idObjectId bool, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader) {
	return NewMongoSearchLoaderWithQueryAndSort(db, collection, modelType, idObjectId, buildQuery, getSort, BuildSort, options...)
}
func NewSearchLoaderWithQuery(db *mongo.Database, collection string, modelType reflect.Type, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader) {
	return NewMongoSearchLoaderWithQueryAndSort(db, collection, modelType, false, buildQuery, getSort, BuildSort, options...)
}
func NewMongoSearchLoaderWithQueryAndSort(db *mongo.Database, collection string, modelType reflect.Type, idObjectId bool, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, buildSort func(string, reflect.Type) bson.D, options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) > 0 && options[0] != nil {
		mp = options[0]
	}
	loader := NewMongoLoader(db, collection, modelType, idObjectId, mp)
	builder := NewSearchBuilderWithSort(db, collection, buildQuery, getSort, buildSort, mp)
	searcher := NewSearcher(builder.Search)
	return searcher, loader
}
func NewMongoSearchLoader(db *mongo.Database, collection string, modelType reflect.Type, idObjectId bool, search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) > 0 && options[0] != nil {
		mp = options[0]
	}
	loader := NewMongoLoader(db, collection, modelType, idObjectId, mp)
	searcher := NewSearcher(search)
	return searcher, loader
}
func NewSearchLoader(db *mongo.Database, collection string, modelType reflect.Type, search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), options ...func(context.Context, interface{}) (interface{}, error)) (*Searcher, *Loader) {
	return NewMongoSearchLoader(db, collection, modelType, false, search, options...)
}
