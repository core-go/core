package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
	"reflect"
)

func BuildQuery(collection *firestore.CollectionRef, queries []Query, limit int, selectFields ...string) firestore.Query {
	var q firestore.Query
	if limit != 0 {
		q = collection.Limit(limit)
	}
	if len(queries) == 0 {
		return collection.Select(selectFields...)
	}
	for i, p := range queries {
		if i == 0 {
			q = collection.Where(p.Path, p.Operator, p.Value)
		}
		q = q.Where(p.Path, p.Operator, p.Value)
	}
	return q
}
func GetDocuments(ctx context.Context, collection *firestore.CollectionRef, where []Query, limit int) *firestore.DocumentIterator {
	if len(where) > 0 {
		return BuildQuery(collection, where, limit).Documents(ctx)
	}
	if limit != 0 {
		return collection.Limit(limit).Documents(ctx)
	}
	return collection.Documents(ctx)
}
func Exist(ctx context.Context, collection *firestore.CollectionRef, docID string) (bool, error) {
	_, err := collection.Doc(docID).Get(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}
func Find(ctx context.Context, collection *firestore.CollectionRef, where []Query, modelType reflect.Type) (interface{}, error) {
	idx, _, _ := FindIdField(modelType)
	if idx < 0 {
		log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex GetById, ExistsById, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
	}
	return FindWithIdIndexAndTracking(ctx, collection, where, modelType, idx, -1, -1)
}
func FindWithIdIndex(ctx context.Context, collection *firestore.CollectionRef, where []Query, modelType reflect.Type, idIndex int) (interface{}, error) {
	return FindWithIdIndexAndTracking(ctx, collection, where, modelType, idIndex, -1, -1)
}
func FindWithTracking(ctx context.Context, collection *firestore.CollectionRef, where []Query, modelType reflect.Type, createdTimeIndex int, updatedTimeIndex int) (interface{}, error) {
	idx, _, _ := FindIdField(modelType)
	if idx < 0 {
		log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex GetById, ExistsById, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
	}
	return FindWithIdIndexAndTracking(ctx, collection, where, modelType, idx, createdTimeIndex, updatedTimeIndex)
}
func FindWithIdIndexAndTracking(ctx context.Context, collection *firestore.CollectionRef, where []Query, modelType reflect.Type, idIndex int, createdTimeIndex int, updatedTimeIndex int) (interface{}, error) {
	iter := GetDocuments(ctx, collection, where, 0)
	modelsType := reflect.Zero(reflect.SliceOf(modelType)).Type()
	arr := reflect.New(modelsType).Interface()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		result := reflect.New(modelType).Interface()
		err = doc.DataTo(&result)
		if err != nil {
			return nil, err
		}
		BindCommonFields(result, doc, idIndex, createdTimeIndex, updatedTimeIndex)
		//SetValue(result, idIndex, doc.Ref.ID)
		arr = appendToArray(arr, result)
	}
	return arr, nil
}
func FindOneWithIdIndexAndTracking(ctx context.Context, collection *firestore.CollectionRef, docID string, modelType reflect.Type, idIndex int, createdTimeIndex int, updatedTimeIndex int) (interface{}, error) {
	doc, er1 := collection.Doc(docID).Get(ctx)
	if er1 != nil {
		return nil, er1
	}
	result := reflect.New(modelType).Interface()
	er2 := doc.DataTo(&result)
	if er2 != nil {
		return nil, er2
	}
	BindCommonFields(result, doc, idIndex, createdTimeIndex, updatedTimeIndex)
	//SetValue(result, idIndex, doc.Ref.ID)
	return result, nil
}
func FindOneAndDecodeWithIdIndexAndTracking(ctx context.Context, collection *firestore.CollectionRef, docID string, result interface{}, idIndex int, createdTimeIndex int, updatedTimeIndex int) (bool, error) {
	doc, err := collection.Doc(docID).Get(ctx)
	if err != nil {
		return false, err
	}
	err = doc.DataTo(result)
	if err != nil {
		return false, err
	}
	BindCommonFields(result, doc, idIndex, createdTimeIndex, updatedTimeIndex)
	return true, nil
}
