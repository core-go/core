package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"reflect"
)

type Writer struct {
	*Loader
	maps             map[string]string
	versionField     string
	versionJson      string
	versionFirestore string
	versionIndex     int
}

func NewWriter(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *Writer {
	return NewWriterWithVersion(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, "", options...)
}
func NewWriterWithVersion(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, versionField string, options ...string) *Writer {
	loader := NewLoader(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, options...)
	maps := MakeFirestoreMap(modelType)
	if len(versionField) > 0 {
		index, versionJson, versionFirestore := FindFieldByName(modelType, versionField)
		if index >= 0 {
			return &Writer{Loader: loader, maps: maps, versionField: versionField, versionIndex: index, versionJson: versionJson, versionFirestore: versionFirestore}
		}
	}
	return &Writer{Loader: loader, maps: maps, versionIndex: -1}
}

func (s *Writer) Insert(ctx context.Context, model interface{}) (int64, error) {
	mv := reflect.ValueOf(model)
	id := reflect.Indirect(mv).Field(s.idIndex).Interface().(string)
	if s.versionIndex >= 0 {
		return InsertOneWithVersion(ctx, s.Collection, id, model, s.versionIndex)
	}
	return InsertOne(ctx, s.Collection, id, model)
}

func (s *Writer) Update(ctx context.Context, model interface{}) (int64, error) {
	mv := reflect.ValueOf(model)
	id := reflect.Indirect(mv).Field(s.idIndex).Interface().(string)
	if s.versionIndex >= 0 {
		return UpdateOneWithVersion(ctx, s.Collection, model, s.versionIndex, s.versionField, s.idIndex)
	}
	return UpdateOne(ctx, s.Collection, id, model)
}

func (s *Writer) Patch(ctx context.Context, data map[string]interface{}) (int64, error) {
	id := data[s.jsonIdName]
	if s.versionIndex >= 0 {
		return PatchOneWithVersion(ctx, s.Collection, id.(string), data, s.maps, s.versionJson)
	}
	delete(data, s.jsonIdName)
	return PatchOne(ctx, s.Collection, id.(string), data, s.maps)
}

func (s *Writer) Save(ctx context.Context, model interface{}) (int64, error) {
	mv := reflect.ValueOf(model)
	id := reflect.Indirect(mv).Field(s.idIndex).Interface().(string)
	if s.versionIndex >= 0 {
		return SaveOneWithVersion(ctx, s.Collection, id, model, s.versionIndex, s.versionField)
	}
	exist, er1 := Exist(ctx, s.Collection, id)
	if er1 != nil {
		return 0, er1
	}
	if !exist {
		return InsertOne(ctx, s.Collection, id, model)
	}
	return UpdateOne(ctx, s.Collection, id, model)
}

func (s *Writer) Delete(ctx context.Context, id interface{}) (int64, error) {
	sid := id.(string)
	return DeleteOne(ctx, s.Collection, sid)
}
