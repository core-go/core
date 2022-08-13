package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type SearchBuilder struct {
	Collection *mongo.Collection
	BuildQuery func(m interface{}) (bson.D, bson.M)
	GetSort    func(m interface{}) string
	BuildSort  func(s string, modelType reflect.Type) bson.D
	Map        func(ctx context.Context, model interface{}) (interface{}, error)
}

func NewSearchBuilderWithSort(db *mongo.Database, collectionName string, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, buildSort func(string, reflect.Type) bson.D, options ...func(context.Context, interface{}) (interface{}, error)) *SearchBuilder {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) > 0 && options[0] != nil {
		mp = options[0]
	}
	collection := db.Collection(collectionName)
	builder := &SearchBuilder{Collection: collection, BuildQuery: buildQuery, GetSort: getSort, BuildSort: buildSort, Map: mp}
	return builder
}
func NewSearchBuilder(db *mongo.Database, collectionName string, buildQuery func(interface{}) (bson.D, bson.M), getSort func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) *SearchBuilder {
	return NewSearchBuilderWithSort(db, collectionName, buildQuery, getSort, BuildSort, options...)
}
func (b *SearchBuilder) Search(ctx context.Context, m interface{}, results interface{}, limit int64, options ...int64) (int64, string, error) {
	query, fields := b.BuildQuery(m)

	var sort = bson.D{}
	s := b.GetSort(m)
	modelType := reflect.TypeOf(results).Elem().Elem()
	sort = b.BuildSort(s, modelType)
	var skip int64 = 0
	if len(options) > 0 && options[0] > 0 {
		skip = options[0]
	}
	count, err := BuildSearchResult(ctx, b.Collection, results, query, fields, sort, limit, skip, b.Map)
	return count, "", err
}
