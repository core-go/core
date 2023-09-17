package firestore

import (
	"cloud.google.com/go/firestore"
	"reflect"
)

func NewRepository(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *GenericWriter {
	return NewGenericWriterWithVersion(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, "", options...)
}
func NewRepositoryWithVersion(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, versionField string, options ...string) *GenericWriter {
	return NewGenericWriterWithVersion(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, versionField, options...)
}
func NewAdapter(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, options ...string) *Writer {
	return NewWriterWithVersion(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, "", options...)
}
func NewAdapterWithVersion(client *firestore.Client, collectionName string, modelType reflect.Type, createdTimeFieldName string, updatedTimeFieldName string, versionField string, options ...string) *Writer {
	return NewWriterWithVersion(client, collectionName, modelType, createdTimeFieldName, updatedTimeFieldName, versionField, options...)
}
