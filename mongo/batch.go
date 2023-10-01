package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"strings"
)

func appendToArray(arr interface{}, item interface{}) interface{} {
	arrValue := reflect.ValueOf(arr)
	elemValue := reflect.Indirect(arrValue)

	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = reflect.Indirect(itemValue)
	}
	elemValue.Set(reflect.Append(elemValue, itemValue))
	return arr
}
func InsertMany(ctx context.Context, collection *mongo.Collection, models interface{}) (bool, error) {
	arr := make([]interface{}, 0)
	values := reflect.Indirect(reflect.ValueOf(models))
	length := values.Len()
	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice:
		for i := 0; i < length; i++ {
			arr = append(arr, values.Index(i).Interface())
		}
	}

	if len(arr) > 0 {
		res, err := collection.InsertMany(ctx, arr)
		if err != nil {
			if strings.Index(err.Error(), "duplicate key error collection:") >= 0 {
				return true, nil
			} else {
				return false, err
			}
		}

		valueOfModel := reflect.Indirect(reflect.ValueOf(arr[0]))
		idIndex, _, _ := FindIdField(valueOfModel.Type())
		if idIndex >= 0 {
			for i, _ := range arr {
				if idValue, ok := res.InsertedIDs[i].(primitive.ObjectID); ok {
					mapObjectIdToModel(idValue, values.Index(i), idIndex)
				}
			}
		}
	}
	return false, nil
}
func InsertManySkipErrors(ctx context.Context, collection *mongo.Collection, models interface{}) (interface{}, interface{}, error) {
	arr := make([]interface{}, 0)
	indexFailArr := make([]int, 0)
	modelsType := reflect.TypeOf(models)
	insertedFails := reflect.New(modelsType).Interface()
	idName := ""
	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(models)
		if values.Len() == 0 {
			return insertedFails, insertedFails, nil
		}
		_, name, _ := FindIdField(reflect.TypeOf(values.Index(0).Interface()))
		idName = name
		for i := 0; i < values.Len(); i++ {
			arr = append(arr, values.Index(i).Interface())
		}
	}
	var defaultOrdered = false
	rs, err := collection.InsertMany(ctx, arr, &options.InsertManyOptions{Ordered: &defaultOrdered})
	if err != nil {
		values := reflect.ValueOf(models)
		insertedSuccess := reflect.New(modelsType).Interface()
		if bulkWriteException, ok := err.(mongo.BulkWriteException); ok {
			for _, writeError := range bulkWriteException.WriteErrors {
				appendToArray(insertedFails, values.Index(writeError.Index).Interface())
				indexFailArr = append(indexFailArr, writeError.Index)
			}
			if rs != nil && len(idName) > 0 {
				insertedSuccess = mapIdInObjects(models, indexFailArr, rs.InsertedIDs, modelsType, idName)
			}
			return insertedSuccess, insertedFails, err
		} else {
			for i := 0; i < values.Len(); i++ {
				appendToArray(insertedFails, values.Index(i).Interface())
			}
			return insertedSuccess, insertedFails, err
		}
	}
	if len(idName) > 0 {
		insertedSuccess := mapIdInObjects(models, indexFailArr, rs.InsertedIDs, modelsType, idName)
		return insertedSuccess, nil, err
	}
	return nil, nil, err
}
func mapIdInObjects(models interface{}, arrayFailIndexIgnore []int, insertedIDs []interface{}, modelsType reflect.Type, fieldName string) interface{} {
	insertedSuccess := reflect.New(modelsType).Interface()
	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(models)
		length := values.Len()
		if length > 0 && length == len(insertedIDs) {
			if index := findIndex(values.Index(0).Interface(), fieldName); index != -1 {
				for i := 0; i < length; i++ {
					if !existInArray(arrayFailIndexIgnore, i) {
						if id, ok := insertedIDs[i].(primitive.ObjectID); ok {
							itemValue := values.Index(i)
							var errSet error
							var vSet interface{}
							switch reflect.Indirect(itemValue).FieldByName(fieldName).Kind() {
							case reflect.String:
								idString := id.Hex()
								vSet, errSet = setValue(itemValue, index, idString)
								break
							default:
								vSet, errSet = setValue(itemValue, index, id)
								break
							}
							if errSet == nil {
								appendToArray(insertedSuccess, vSet)
							} else {
								appendToArray(insertedSuccess, itemValue.Interface())
								log.Println("Error map Id: ", errSet)
							}
						}
					}
				}
			}
		}
	}
	return insertedSuccess
}
func existInArray(arr []int, value interface{}) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
func InArray(value int, arr []int) bool {
	for i := 0; i < len(arr); i++ {
		if value == arr[i] {
			return true
		}
	}
	return false
}
func UpdateMany(ctx context.Context, collection *mongo.Collection, models interface{}, idName string) (*mongo.BulkWriteResult, error) {
	models_ := make([]mongo.WriteModel, 0)
	if reflect.TypeOf(models).Kind() == reflect.Slice {
		values := reflect.ValueOf(models)
		length := values.Len()
		if length > 0 {
			if index := findIndex(values.Index(0).Interface(), idName); index != -1 {
				for i := 0; i < length; i++ {
					row := values.Index(i).Interface()
					v, er0 := getValue(row, index)
					if er0 != nil {
						return nil, er0
					}
					updateQuery := bson.M{
						"$set": row,
					}
					updateModel := mongo.NewUpdateOneModel().SetUpdate(updateQuery).SetFilter(bson.M{"_id": v})
					models_ = append(models_, updateModel)
				}
			}
		}
	}
	res, err := collection.BulkWrite(ctx, models_)
	return res, err
}

// Patch
func PatchMaps(ctx context.Context, collection *mongo.Collection, maps []map[string]interface{}, idName string) (*mongo.BulkWriteResult, error) {
	if idName == "" {
		idName = "_id"
	}
	writeModels := make([]mongo.WriteModel, 0)
	for _, row := range maps {
		v, _ := row[idName]
		if v != nil {
			updateModel := mongo.NewUpdateOneModel().SetUpdate(bson.M{
				"$set": row,
			}).SetFilter(bson.M{"_id": v})
			writeModels = append(writeModels, updateModel)
		}
	}
	res, err := collection.BulkWrite(ctx, writeModels)
	return res, err
}
func UpsertMany(ctx context.Context, collection *mongo.Collection, model interface{}, idName string) (*mongo.BulkWriteResult, error) { //Patch
	models := make([]mongo.WriteModel, 0)
	switch reflect.TypeOf(model).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(model)

		n := values.Len()
		if n > 0 {
			if index := findIndex(values.Index(0).Interface(), idName); index != -1 {
				for i := 0; i < n; i++ {
					row := values.Index(i).Interface()
					id, er0 := getValue(row, index)
					if er0 != nil {
						return nil, er0
					}
					if id != nil || (reflect.TypeOf(id).String() == "string") || (reflect.TypeOf(id).String() == "string" && len(id.(string)) > 0) { // if exist
						updateModel := mongo.NewReplaceOneModel().SetUpsert(true).SetReplacement(row).SetFilter(bson.M{"_id": id})
						models = append(models, updateModel)
					} else {
						insertModel := mongo.NewInsertOneModel().SetDocument(row)
						models = append(models, insertModel)
					}
				}
			}
		}
	}
	rs, err := collection.BulkWrite(ctx, models)
	return rs, err
}
