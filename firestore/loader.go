package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"reflect"
)

type Loader struct {
	Client           *firestore.Client
	Collection       *firestore.CollectionRef
	modelType        reflect.Type
	idIndex          int
	jsonIdName       string
	collectionName   string
	createdTimeIndex int
	updatedTimeIndex int
}

func NewLoader(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *Loader {
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
	return &Loader{Client: client, Collection: client.Collection(collectionName), modelType: modelType, idIndex: idx, collectionName: collectionName, jsonIdName: jsonIdName, createdTimeIndex: ctIdx, updatedTimeIndex: utIdx}
}

func (s *Loader) Id() string {
	return s.jsonIdName
}

func (s *Loader) All(ctx context.Context) (interface{}, error) {
	query := make([]Query, 0)
	return FindWithIdIndexAndTracking(ctx, s.Collection, query, s.modelType, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *Loader) Load(ctx context.Context, id interface{}) (interface{}, error) {
	sid := id.(string)
	return FindOneWithIdIndexAndTracking(ctx, s.Collection, sid, s.modelType, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *Loader) Get(ctx context.Context, id interface{}, result interface{}) (bool, error) {
	sid := id.(string)
	return FindOneAndDecodeWithIdIndexAndTracking(ctx, s.Collection, sid, result, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *Loader) LoadAndDecode(ctx context.Context, id interface{}, result interface{}) (bool, error) {
	sid := id.(string)
	return FindOneAndDecodeWithIdIndexAndTracking(ctx, s.Collection, sid, result, s.idIndex, s.createdTimeIndex, s.updatedTimeIndex)
}

func (s *Loader) Exist(ctx context.Context, id interface{}) (bool, error) {
	sid := id.(string)
	return Exist(ctx, s.Collection, sid)
}

/*
func (s *Loader) LoadByIds(ctx context.Context, ids []string) (interface{}, []string, []error) {
	return FindByIds(ctx, s.Collection, ids, s.modelType, s.jsonIdName)
}

func (s *Loader) LoadByIdsAndDecode(ctx context.Context, ids []string, result interface{}) ([]string, []error) {
	return FindByIdsAndDecode(ctx, s.Collection, ids, result, s.jsonIdName)
}
*/
