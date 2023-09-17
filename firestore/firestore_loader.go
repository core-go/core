package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"reflect"
)

type FirestoreLoader struct {
	Client           *firestore.Client
	Collection       *firestore.CollectionRef
	modelType        reflect.Type
	idIndex          int
	jsonIdName       string
	collectionName   string
	createdTimeIndex int
	updatedTimeIndex int
}

func NewFirestoreLoader(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *FirestoreLoader {
	idx := -1
	var idFieldName string
	if len(options) > 0 && len(options[0]) > 0 {
		idFieldName = options[0]
	}
	var jsonIdName string
	if len(idFieldName) == 0 {
		idx, _, jsonIdName = FindIdField(modelType)
		if idx < 0 {
			log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex Load, Exist, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
		}
	} else {
		idx, jsonIdName, _ = FindFieldByName(modelType, idFieldName)
		if idx < 0 {
			log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex Load, Exist, Save, Update) because don't have any fields of " + modelType.Name())
		}
	}
	ctIdx := -1
	if len(createdTimeFieldName) >= 0 {
		ctIdx, _, _ = FindFieldByName(modelType, createdTimeFieldName)
	}
	utIdx := -1
	if len(updatedTimeFieldName) >= 0 {
		utIdx, _, _ = FindFieldByName(modelType, updatedTimeFieldName)
	}
	return &FirestoreLoader{Client: client, Collection: client.Collection(collectionName), modelType: modelType, idIndex: idx, collectionName: collectionName, jsonIdName: jsonIdName, createdTimeIndex: ctIdx, updatedTimeIndex: utIdx}
}

func (s *FirestoreLoader) Id() string {
	return s.jsonIdName
}

func (s *FirestoreLoader) All(ctx context.Context) (interface{}, error) {
	query := make([]Query, 0)
	return FindWithIdIndexAndTracking(ctx, s.Collection, query, s.modelType, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *FirestoreLoader) Load(ctx context.Context, id interface{}) (interface{}, error) {
	sid := id.(string)
	return FindOneWithIdIndexAndTracking(ctx, s.Collection, sid, s.modelType, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *FirestoreLoader) Get(ctx context.Context, id string, result interface{}) (bool, error) {
	return FindOneAndDecodeWithIdIndexAndTracking(ctx, s.Collection, id, result, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *FirestoreLoader) LoadAndDecode(ctx context.Context, id string, result interface{}) (bool, error) {
	return FindOneAndDecodeWithIdIndexAndTracking(ctx, s.Collection, id, result, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *FirestoreLoader) Exist(ctx context.Context, id string) (bool, error) {
	return Exist(ctx, s.Collection, id)
}
