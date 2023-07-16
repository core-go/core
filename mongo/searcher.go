package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type Searcher struct {
	search func(ctx context.Context, searchModel interface{}, results interface{}, limit int64, skip int64) (int64, error)
}

func NewSearcher(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error)) *Searcher {
	return &Searcher{search: search}
}

func (s *Searcher) Search(ctx context.Context, m interface{}, results interface{}, limit int64, skip int64) (int64, error) {
	return s.search(ctx, m, results, limit, skip)
}

func NewSearcherWithQuery(db *mongo.Database, collectionName string, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) *Searcher {
	return NewSearcherWithQueryAndSort(db, collectionName, buildQuery, getSort, BuildSort, options...)
}
func NewSearcherWithQueryAndSort(db *mongo.Database, collectionName string, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, buildSort func(string, reflect.Type) bson.D, options ...func(context.Context, interface{}) (interface{}, error)) *Searcher {
	builder := NewSearchBuilderWithSort(db, collectionName, buildQuery, getSort, buildSort, options...)
	return NewSearcher(builder.Search)
}
