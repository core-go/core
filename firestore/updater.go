package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"reflect"
)

type Updater struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
	IdName     string
	idx        int
	modelType  reflect.Type
	modelsType reflect.Type
	Map        func(ctx context.Context, model interface{}) (interface{}, error)
}

func NewUpdaterWithIdName(client *firestore.Client, collectionName string, modelType reflect.Type, fieldName string, options ...func(context.Context, interface{}) (interface{}, error)) *Updater {
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
	return &Updater{client: client, collection: collection, IdName: fieldName, idx: idx, modelType: modelType, modelsType: modelsType, Map: mp}
}

func NewUpdater(client *firestore.Client, collectionName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) *Updater {
	return NewUpdaterWithIdName(client, collectionName, modelType, "", options...)
}

func (w *Updater) Write(ctx context.Context, model interface{}) error {
	id := getIdValueFromModel(model, w.idx)
	if w.Map != nil {
		m2, er0 := w.Map(ctx, model)
		if er0 != nil {
			return er0
		}
		_, er1 := UpdateOne(ctx, w.collection, id, m2)
		return er1
	}
	_, er2 := UpdateOne(ctx, w.collection, id, model)
	return er2
}
