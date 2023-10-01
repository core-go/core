package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

// For Batch Update
func initArrayResults(modelsType reflect.Type) interface{} {
	return reflect.New(modelsType).Interface()
}

func MapToMongoObjects(model interface{}, idName string, idObjectId bool, modelType reflect.Type, newId bool) (interface{}, interface{}) {
	var results = initArrayResults(modelType)
	var ids = make([]interface{}, 0)
	switch reflect.TypeOf(model).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(model)
		for i := 0; i < values.Len(); i++ {
			model, id := MapToMongoObject(values.Index(i).Interface(), idName, idObjectId, newId)
			ids = append(ids, id)
			results = appendToArray(results, model)
		}
	}
	return results, ids
}
func MapToMongoObject(model interface{}, idName string, objectId bool, newId bool) (interface{}, interface{}) {
	if index := findIndex(model, idName); index != -1 {
		id, _ := getValue(model, index)
		if objectId {
			if newId && (id == nil) {
				setValue(model, index, primitive.NewObjectID())
			} else {
				objectId, err := primitive.ObjectIDFromHex(id.(string))
				if err == nil {
					setValue(model, index, objectId)
				}
			}
		} else {
			setValue(model, index, id)
		}
		return model, id
	}
	return model, nil
}
