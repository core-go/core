package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"reflect"
)

type Inserter struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
	IdName     string
	idx        int
	modelType  reflect.Type
	modelsType reflect.Type
	Map        func(ctx context.Context, model interface{}) (interface{}, error)
}

func NewInserterWithIdName(client *firestore.Client, collectionName string, modelType reflect.Type, fieldName string, options ...func(context.Context, interface{}) (interface{}, error)) *Inserter {
	var idx int
	if len(fieldName) == 0 {
		idx, fieldName, _ = FindIdField(modelType)
		if idx < 0 {
			log.Println("Require Id value (Ex Load, Exist, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
		}
	} else {
		idx, _, _ = FindFieldByName(modelType, fieldName)
	}

	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	modelsType := reflect.Zero(reflect.SliceOf(modelType)).Type()
	collection := client.Collection(collectionName)
	return &Inserter{client: client, collection: collection, IdName: fieldName, idx: idx, modelType: modelType, modelsType: modelsType, Map: mp}
}

func NewInserter(client *firestore.Client, collectionName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) *Inserter {
	return NewInserterWithIdName(client, collectionName, modelType, "", options...)
}

func (w *Inserter) Write(ctx context.Context, model interface{}) error {
	if w.Map != nil {
		m2, er0 := w.Map(ctx, model)
		if er0 != nil {
			return er0
		}
		return Insert(ctx, w.collection, w.idx, m2)
	}
	return Insert(ctx, w.collection, w.idx, model)
}
