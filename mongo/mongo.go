package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

func CreateUniqueIndex(collection *mongo.Collection, fieldName string) (string, error) {
	keys := bson.D{{Key: fieldName, Value: 1}}
	indexName, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		},
	)
	return indexName, err
}
func difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}
func FindByIds(ctx context.Context, collection *mongo.Collection, ids []string, idObjectId bool, modelType reflect.Type) (interface{}, []string, error) {
	modelsType := reflect.Zero(reflect.SliceOf(modelType)).Type()
	result := reflect.New(modelsType).Interface()
	res := reflect.Indirect(reflect.ValueOf(result))
	if !idObjectId {
		find, errFind := collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if errFind != nil {
			return nil, ids, errFind
		}
		err := find.All(ctx, result)
		keySuccess := []string{""}
		_, fieldName, _ := FindIdField(modelType)
		for i := 0; i < res.Len(); i++ {
			key := res.Index(i).FieldByName(fieldName).String()
			keySuccess = append(keySuccess, key)
		}
		keys := difference(keySuccess, ids)
		_ = find.Close(ctx)
		return result, keys, err
	} else {
		id := make([]primitive.ObjectID, 0)
		for _, val := range ids {
			item, err := primitive.ObjectIDFromHex(val)
			if err != nil {
				return false, nil, err
			}
			id = append(id, item)
		}
		find, errFind := collection.Find(ctx, bson.M{"_id": bson.M{"$in": id}})
		if errFind != nil {
			return false, ids, errFind
		}
		err := find.All(ctx, result)
		keySuccess := make([]primitive.ObjectID, 0)
		_, fieldName, _ := FindIdField(modelType)
		for i := 0; i < res.Len(); i++ {
			key, _ := primitive.ObjectIDFromHex(res.Index(i).FieldByName(fieldName).Interface().(primitive.ObjectID).Hex())
			keySuccess = append(keySuccess, key)
		}
		keyToStr := []string{""}
		for index, _ := range keySuccess {
			key := keySuccess[index].Hex()
			keyToStr = append(keyToStr, key)
		}
		keys := difference(keyToStr, ids)
		return result, keys, err
	}
}
func FindByIdsAndDecode(ctx context.Context, collection *mongo.Collection, ids []string, idObjectId bool, result interface{}) ([]string, error) {
	var keys []string
	if !idObjectId {
		res := reflect.Indirect(reflect.ValueOf(result))
		find, errFind := collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if errFind != nil {
			return ids, errFind
		}
		find.All(ctx, result)
		keySuccess := []string{""}
		//		_, fieldName := FindIdField(modelType)
		for i := 0; i < res.Len(); i++ {
			keySuccess = append(keySuccess, res.Index(i).Field(0).String())
		}
		keys = difference(keySuccess, ids)
		_ = find.Close(ctx)
	} else {
		res := reflect.Indirect(reflect.ValueOf(result))
		id := make([]primitive.ObjectID, 0)
		for _, val := range ids {
			item, err := primitive.ObjectIDFromHex(val)
			if err != nil {
				return nil, err
			}
			id = append(id, item)
		}
		find, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": id}})
		if err != nil {
			return ids, err
		}
		find.All(ctx, result)
		keySuccess := make([]primitive.ObjectID, 0)
		//_, fieldName := FindIdField(modelType)
		for i := 0; i < res.Len(); i++ {
			key, _ := primitive.ObjectIDFromHex(res.Index(i).Field(0).Interface().(primitive.ObjectID).Hex())
			keySuccess = append(keySuccess, key)
		}
		keyToStr := []string{""}
		for index, _ := range keySuccess {
			key := keySuccess[index].Hex()
			keyToStr = append(keyToStr, key)
		}
		keys = difference(keyToStr, ids)
	}
	if result != nil {
		return keys, nil
	}
	return keys, errors.New("no result return")
}
func Find(ctx context.Context, collection *mongo.Collection, query bson.M, modelType reflect.Type) (interface{}, error) {
	cur, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	modelsType := reflect.Zero(reflect.SliceOf(modelType)).Type()
	arr := reflect.New(modelsType).Interface()
	er2 := cur.All(ctx, arr)
	_ = cur.Close(ctx)
	return arr, er2
}
func PatchMany(ctx context.Context, collection *mongo.Collection, models interface{}, idName string) (*mongo.BulkWriteResult, error) {
	models_ := make([]mongo.WriteModel, 0)
	ids := make([]interface{}, 0)
	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(models)
		length := values.Len()
		if length > 0 {
			if index := findIndex(values.Index(0).Interface(), idName); index != -1 {
				for i := 0; i < length; i++ {
					row := values.Index(i).Interface()
					updateModel := mongo.NewUpdateOneModel().SetUpdate(values.Index(i))
					v, err1 := getValue(row, index)
					if err1 == nil && v != nil {
						if reflect.TypeOf(v).String() != "string" {
							updateModel = mongo.NewUpdateOneModel().SetUpdate(bson.M{
								"$set": row,
							}).SetFilter(bson.M{"_id": v})
						} else {
							if idStr, ok := v.(string); ok {
								updateModel = mongo.NewUpdateOneModel().SetUpdate(bson.M{
									"$set": row,
								}).SetFilter(bson.M{"_id": idStr})
							}
						}
						ids = append(ids, v)
					}
					models_ = append(models_, updateModel)
				}
			}
		}
	}
	var defaultOrdered = false
	return collection.BulkWrite(ctx, models_, &options.BulkWriteOptions{Ordered: &defaultOrdered})
}
func PatchManyAndGetSuccessList(ctx context.Context, collection *mongo.Collection, models interface{}, idName string) (interface{}, interface{}, []interface{}, error) {
	models_ := make([]mongo.WriteModel, 0)
	ids := make([]interface{}, 0)
	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(models)
		length := values.Len()
		if length > 0 {
			if index := findIndex(values.Index(0).Interface(), idName); index != -1 {
				for i := 0; i < length; i++ {
					row := values.Index(i).Interface()
					updateModel := mongo.NewUpdateOneModel().SetUpdate(values.Index(i))
					v, err1 := getValue(row, index)
					if err1 == nil && v != nil {
						if reflect.TypeOf(v).String() != "string" {
							updateModel = mongo.NewUpdateOneModel().SetUpdate(bson.M{
								"$set": row,
							}).SetFilter(bson.M{"_id": v})
						} else {
							if idStr, ok := v.(string); ok {
								updateModel = mongo.NewUpdateOneModel().SetUpdate(bson.M{
									"$set": row,
								}).SetFilter(bson.M{"_id": idStr})
							}
						}
						ids = append(ids, v)
					}
					models_ = append(models_, updateModel)
				}
			}
		}
	}
	var defaultOrdered = false
	_, er0 := collection.BulkWrite(ctx, models_, &options.BulkWriteOptions{Ordered: &defaultOrdered})
	if er0 != nil {
		return nil, nil, nil, er0
	}
	successIdList := make([]interface{}, 0)
	_options := options.FindOptions{Projection: bson.M{"_id": 1}}
	cur, er1 := collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}}, &_options)
	if er1 != nil {
		return nil, nil, nil, er1
	}
	er2 := cur.All(ctx, &successIdList)
	if er2 != nil {
		return nil, nil, nil, er2
	}
	successIdList = mapArrayInterface(successIdList)
	failList, failIdList := diffModelArray(models, successIdList, idName)
	return successIdList, failList, failIdList, er0
}
func mapArrayInterface(successIdList []interface{}) []interface{} {
	arr := make([]interface{}, 0)
	for _, value := range successIdList {
		if primitiveE, ok := value.(primitive.D); ok {
			for _, itemPrimitiveE := range primitiveE {
				arr = append(arr, itemPrimitiveE.Value)
			}
		}
	}
	return arr
}
func diffModelArray(modelsAll interface{}, successIdList interface{}, idName string) (interface{}, []interface{}) {
	modelsB := make([]interface{}, 0)
	modelBId := make([]interface{}, 0)
	switch reflect.TypeOf(modelsAll).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(modelsAll)
		length := values.Len()
		if length > 0 {
			if index := findIndex(values.Index(0).Interface(), idName); index != -1 {
				for i := 0; i < length; i++ {
					itemValue := values.Index(i)
					id, _ := getValue(itemValue.Interface(), index)
					if !existInArrayInterface(successIdList, id) {
						modelsB = append(modelsB, itemValue.Interface())
						modelBId = append(modelBId, id)
					}
				}
			}
		}
	}
	return modelsB, modelBId
}
func existInArrayInterface(arr interface{}, valueID interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		values := reflect.ValueOf(arr)
		for i := 0; i < values.Len(); i++ {
			itemValueID := values.Index(i).Interface()
			if itemValueID == valueID {
				return true
			}
		}
	}
	return false
}
func UpdateMaps(ctx context.Context, collection *mongo.Collection, maps []map[string]interface{}, idName string) (*mongo.BulkWriteResult, error) {
	if idName == "" {
		idName = "_id"
	}
	models_ := make([]mongo.WriteModel, 0)
	for _, row := range maps {
		v, _ := row[idName]
		if v != nil {
			updateModel := mongo.NewReplaceOneModel().SetReplacement(bson.M{
				"$set": row,
			}).SetFilter(bson.M{"_id": v})
			models_ = append(models_, updateModel)
		}
	}
	res, err := collection.BulkWrite(ctx, models_)
	return res, err
}
func UpsertMaps(ctx context.Context, collection *mongo.Collection, maps []map[string]interface{}, idName string) (*mongo.BulkWriteResult, error) {
	models_ := make([]mongo.WriteModel, 0)
	for _, row := range maps {
		id, _ := row[idName]
		if id != nil || (reflect.TypeOf(id).String() == "string") || (reflect.TypeOf(id).String() == "string" && len(id.(string)) > 0) {
			updateModel := mongo.NewUpdateOneModel().SetUpdate(bson.M{
				"$set": row,
			}).SetUpsert(true).SetFilter(bson.M{"_id": id})
			models_ = append(models_, updateModel)
		} else {
			insertModel := mongo.NewInsertOneModel().SetDocument(row)
			models_ = append(models_, insertModel)
		}
	}
	res, err := collection.BulkWrite(ctx, models_)
	return res, err
}
func FindFieldByName(modelType reflect.Type, fieldName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			name1 := fieldName
			name2 := fieldName
			tag1, ok1 := field.Tag.Lookup("json")
			tag2, ok2 := field.Tag.Lookup("bson")
			if ok1 {
				name1 = strings.Split(tag1, ",")[0]
			}
			if ok2 {
				name2 = strings.Split(tag2, ",")[0]
			}
			return i, name1, name2
		}
	}
	return -1, fieldName, fieldName
}
//For Search and Patch
func GetBsonName(modelType reflect.Type, fieldName string) string {
	field, found := modelType.FieldByName(fieldName)
	if !found {
		return fieldName
	}
	if tag, ok := field.Tag.Lookup("bson"); ok {
		return strings.Split(tag, ",")[0]
	}
	return fieldName
}
func GetBsonNameByIndex(modelType reflect.Type, fieldIndex int) string {
	if tag, ok := modelType.Field(fieldIndex).Tag.Lookup("bson"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}
