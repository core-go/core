package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"reflect"
)

type Searcher struct {
	search func(context.Context, interface{}, interface{}, int64, string) (string, error)
}

func NewSearcher(search func(context.Context, interface{}, interface{}, int64, string) (string, error)) *Searcher {
	return &Searcher{search: search}
}

func (s *Searcher) Search(ctx context.Context, m interface{}, results interface{}, pageSize int64, nextPageToken string) (string, error) {
	return s.search(ctx, m, results, pageSize, nextPageToken)
}

func NewSearcherWithQueryAndSort(client *firestore.Client, collectionName string, modelType reflect.Type, buildQuery func(interface{}) ([]Query, []string), getSort func(interface{}) string, buildSort func(s string, modelType reflect.Type) map[string]firestore.Direction, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *Searcher {
	builder := NewSearchBuilderWithQuery(client, collectionName, modelType, buildQuery, getSort, buildSort, createdTimeFieldName, updatedTimeFieldName, options...)
	return NewSearcher(builder.Search)
}

func NewSearcherWithQuery(client *firestore.Client, collectionName string, modelType reflect.Type, buildQuery func(interface{}) ([]Query, []string), getSort func(interface{}) string, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *Searcher {
	return NewSearcherWithQueryAndSort(client, collectionName, modelType, buildQuery, getSort, BuildSort, createdTimeFieldName, updatedTimeFieldName, options...)
}
