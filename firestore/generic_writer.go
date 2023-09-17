package firestore

import (
	"context"
	"reflect"

	"cloud.google.com/go/firestore"
)

type Repository interface {
	Get(ctx context.Context, id string, result interface{}) (bool, error)
	Exist(ctx context.Context, id string) (bool, error)
	Insert(ctx context.Context, model interface{}) (int64, error)
	Update(ctx context.Context, model interface{}) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

type GenericWriter struct {
	*FirestoreLoader
	maps             map[string]string
	versionField     string
	versionJson      string
	versionFirestore string
	versionIndex     int
}

func NewGenericWriter(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *Writer {
	return NewWriterWithVersion(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, "", options...)
}
func NewGenericWriterWithVersion(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, versionField string, options ...string) *GenericWriter {
	loader := NewFirestoreLoader(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, options...)
	maps := MakeFirestoreMap(modelType)
	if len(versionField) > 0 {
		index, versionJson, versionFirestore := FindFieldByName(modelType, versionField)
		if index >= 0 {
			return &GenericWriter{FirestoreLoader: loader, maps: maps, versionField: versionField, versionIndex: index, versionJson: versionJson, versionFirestore: versionFirestore}
		}
	}
	return &GenericWriter{FirestoreLoader: loader, maps: maps, versionIndex: -1}
}

func (s *GenericWriter) Insert(ctx context.Context, model interface{}) (int64, error) {
	mv := reflect.ValueOf(model)
	id := reflect.Indirect(mv).Field(s.idIndex).Interface().(string)
	if s.versionIndex >= 0 {
		return InsertOneWithVersion(ctx, s.Collection, id, model, s.versionIndex)
	}
	return InsertOne(ctx, s.Collection, id, model)
}

func (s *GenericWriter) Update(ctx context.Context, model interface{}) (int64, error) {
	mv := reflect.ValueOf(model)
	id := reflect.Indirect(mv).Field(s.idIndex).Interface().(string)
	if s.versionIndex >= 0 {
		return UpdateOneWithVersion(ctx, s.Collection, model, s.versionIndex, s.versionField, s.idIndex)
	}
	return UpdateOne(ctx, s.Collection, id, model)
}

func (s *GenericWriter) Patch(ctx context.Context, data map[string]interface{}) (int64, error) {
	id := data[s.jsonIdName]
	if s.versionIndex >= 0 {
		return PatchOneWithVersion(ctx, s.Collection, id.(string), data, s.maps, s.versionJson)
	}
	delete(data, s.jsonIdName)
	return PatchOne(ctx, s.Collection, id.(string), data, s.maps)
}

func (s *GenericWriter) Save(ctx context.Context, model interface{}) (int64, error) {
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

func (s *GenericWriter) Delete(ctx context.Context, id string) (int64, error) {
	return DeleteOne(ctx, s.Collection, id)
}
