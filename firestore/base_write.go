package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"reflect"
)

func getIdValueFromModel(model interface{}, idIndex int) string {
	if id, exist := getFieldValueAtIndex(model, idIndex).(string); exist {
		return id
	}
	return ""
}
func getFieldValueAtIndex(model interface{}, index int) interface{} {
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	return modelValue.Field(index).Interface()
}
func setValueWithIndex(model interface{}, index int, value interface{}) (interface{}, error) {
	vo := reflect.Indirect(reflect.ValueOf(model))
	numField := vo.NumField()
	if index >= 0 && index < numField {
		vo.Field(index).Set(reflect.ValueOf(value))
		return model, nil
	}
	return nil, fmt.Errorf("error no found field index: %v", index)
}
func UpdateOne(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}) (int64, error) {
	if len(id) == 0 {
		return 0, fmt.Errorf("cannot update one an object that do not have id field")
	}
	_, err := collection.Doc(id).Set(ctx, model)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
func Insert(ctx context.Context, collection *firestore.CollectionRef, idIndex int, model interface{}) error {
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	idField := modelValue.Field(idIndex)
	if reflect.Indirect(idField).Kind() != reflect.String {
		return fmt.Errorf("the ID field must be string")
	}
	var doc *firestore.DocumentRef
	// TODO apply idField.IsZero() for golang 13 or above
	if idField.Len() == 0 {
		doc = collection.NewDoc()
		idField.Set(reflect.ValueOf(doc.ID))
	} else {
		doc = collection.Doc(idField.String())
	}
	_, err := doc.Create(ctx, model)
	return err
}
