package firestore

import (
	"cloud.google.com/go/firestore"
	"reflect"
)

func NewSearchLoader(client *firestore.Client, collectionName string, modelType reflect.Type, buildQuery func(interface{}) ([]Query, []string), getSort func(interface{}) string, createdTimeFieldName string, updatedTimeFieldName string, options ...string) (*Searcher, *Loader) {
	return NewSearchLoaderWithSort(client, collectionName, modelType, buildQuery, getSort, BuildSort, createdTimeFieldName, updatedTimeFieldName, options...)
}
func NewSearchLoaderWithSort(client *firestore.Client, collectionName string, modelType reflect.Type, buildQuery func(interface{}) ([]Query, []string), getSort func(interface{}) string, buildSort func(s string, modelType reflect.Type) map[string]firestore.Direction, createdTimeFieldName string, updatedTimeFieldName string, options ...string) (*Searcher, *Loader) {
	loader := NewLoader(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, options...)
	searcher := NewSearcherWithQueryAndSort(client, collectionName, modelType, buildQuery, getSort, buildSort, createdTimeFieldName, updatedTimeFieldName, options...)
	return searcher, loader
}
