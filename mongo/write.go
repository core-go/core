package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strings"
)

//For Get By Id
func FindFieldIndex(modelType reflect.Type, fieldName string) int {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			return i
		}
	}
	return -1
}
func MakeBsonMap(modelType reflect.Type) map[string]string {
	maps := make(map[string]string)
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		key1 := field.Name
		if tag0, ok0 := field.Tag.Lookup("json"); ok0 {
			if strings.Contains(tag0, ",") {
				a := strings.Split(tag0, ",")
				key1 = a[0]
			} else {
				key1 = tag0
			}
		}
		if tag, ok := field.Tag.Lookup("bson"); ok {
			if tag != "-" {
				if strings.Contains(tag, ",") {
					a := strings.Split(tag, ",")
					if key1 == "-" {
						key1 = a[0]
					}
					maps[key1] = a[0]
				} else {
					if key1 == "-" {
						key1 = tag
					}
					maps[key1] = tag
				}
			}
		} else {
			if key1 == "-" {
				key1 = field.Name
			}
			maps[key1] = key1
		}
	}
	return maps
}
func InsertOne(ctx context.Context, collection *mongo.Collection, model interface{}) (int64, error) {
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		errMsg := err.Error()
		if strings.Index(errMsg, "duplicate key error collection:") >= 0 {
			return 0, nil
		} else {
			return 0, err
		}
	} else {
		if idValue, ok := result.InsertedID.(primitive.ObjectID); ok {
			valueOfModel := reflect.Indirect(reflect.ValueOf(model))
			typeOfModel := valueOfModel.Type()
			idIndex, _, _ := FindIdField(typeOfModel)
			if idIndex != -1 {
				mapObjectIdToModel(idValue, valueOfModel, idIndex)
			}
		}
		return 1, err
	}
}
func InsertOneWithVersion(ctx context.Context, collection *mongo.Collection, model interface{}, versionIndex int) (int64, error) {
	var defaultVersion interface{}
	modelType := reflect.TypeOf(model).Elem()
	versionType := modelType.Field(versionIndex).Type
	switch versionType.String() {
	case "int":
		defaultVersion = int(1)
	case "int32":
		defaultVersion = int32(1)
	case "int64":
		defaultVersion = int64(1)
	default:
		panic("not support type's version")
	}
	model, err := setValue(model, versionIndex, defaultVersion)
	if err != nil {
		return 0, err
	}
	return InsertOne(ctx, collection, model)
}

//For Update
func BuildQueryId(model interface{}, fieldname string) bson.M {
	query := bson.M{}
	if i := findIndex(model, fieldname); i != -1 {
		id, _ := getValue(model, i)
		query = bson.M{
			"_id": id,
		}
	}
	return query
}
func GetBsonNameByModelIndex(model interface{}, fieldIndex int) string {
	t := reflect.TypeOf(model).Elem()
	if tag, ok := t.Field(fieldIndex).Tag.Lookup("bson"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}
func BuildQueryByIdFromObject(object interface{}) bson.M {
	vo := reflect.Indirect(reflect.ValueOf(object))
	if idIndex, _, _ := FindIdField(vo.Type()); idIndex >= 0 {
		value := vo.Field(idIndex).Interface()
		return bson.M{"_id": value}
	} else {
		panic("id field not found")
	}
}

//Version
func copyMap(originalMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range originalMap {
		newMap[k] = v
	}
	return newMap
}
func BuildIdAndVersionQuery(query map[string]interface{}, model interface{}, versionField string) map[string]interface{} {
	index := findIndex(model, versionField)
	return BuildIdAndVersionQueryByVersionIndex(query, model, index)
}
func BuildIdAndVersionQueryByVersionIndex(query map[string]interface{}, model interface{}, versionIndex int) map[string]interface{} {
	newMap := copyMap(query)
	vo := reflect.Indirect(reflect.ValueOf(model))
	if versionIndex >= 0 && versionIndex < vo.NumField() {
		var valueOfCurrentVersion reflect.Value
		valueOfCurrentVersion = vo.Field(versionIndex)
		versionColumnName := GetBsonNameByModelIndex(model, versionIndex)
		newMap[versionColumnName] = valueOfCurrentVersion.Interface()
		switch valueOfCurrentVersion.Kind().String() {
		case "int":
			{
				nextVersion := reflect.ValueOf(valueOfCurrentVersion.Interface().(int) + 1)
				vo.Field(versionIndex).Set(nextVersion)
			}
		case "int32":
			{
				nextVersion := reflect.ValueOf(valueOfCurrentVersion.Interface().(int32) + 1)
				vo.Field(versionIndex).Set(nextVersion)
			}
		case "int64":
			{
				nextVersion := reflect.ValueOf(valueOfCurrentVersion.Interface().(int64) + 1)
				vo.Field(versionIndex).Set(nextVersion)
			}
		default:
			panic("not support type's version")
		}
		return newMap
	} else {
		panic("invalid versionIndex")
	}
}
func Update(ctx context.Context, collection *mongo.Collection, model interface{}, fieldname string) error {
	query := BuildQueryId(model, fieldname)
	defaultObjID, _ := primitive.ObjectIDFromHex("000000000000")
	if idValue := query["_id"]; !(idValue == "" || idValue == 0 || idValue == defaultObjID) {
		_, err := UpdateOne(ctx, collection, model, query)
		return err
	}
	return errors.New("require field _id")
}
func UpdateOne(ctx context.Context, collection *mongo.Collection, model interface{}, query bson.M) (int64, error) { //Patch
	updateQuery := bson.M{
		"$set": model,
	}
	result, err := collection.UpdateOne(ctx, query, updateQuery)
	if result.ModifiedCount > 0 {
		return result.ModifiedCount, err
	} else if result.UpsertedCount > 0 {
		return result.UpsertedCount, err
	} else {
		return result.MatchedCount, err
	}
}
func UpdateByIdAndVersion(ctx context.Context, collection *mongo.Collection, model interface{}, versionIndex int) (int64, error) {
	idQuery := BuildQueryByIdFromObject(model)
	versionQuery := BuildIdAndVersionQueryByVersionIndex(idQuery, model, versionIndex)
	rowAffect, er1 := UpdateOne(ctx, collection, model, versionQuery)
	if er1 != nil {
		return 0, er1
	}
	if rowAffect == 0 {
		isExist, er2 := Exist(ctx, collection, idQuery["_id"], false)
		if er2 != nil {
			return 0, er2
		}
		if isExist {
			return -1, nil
		} else {
			return 0, nil
		}
	}
	return rowAffect, er1
}

//For Patch
func GetJsonByIndex(modelType reflect.Type, fieldIndex int) string {
	if tag, ok := modelType.Field(fieldIndex).Tag.Lookup("json"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}
func BuildQueryByIdFromMap(m map[string]interface{}, idName string) bson.M {
	if idValue, exist := m[idName]; exist {
		return bson.M{"_id": idValue}
	} else {
		panic("id field not found")
	}
}
func MapToBson(object map[string]interface{}, objectMap map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range object {
		field, ok := objectMap[key]
		if ok {
			result[field] = value
		}
	}
	return result
}
func PatchOne(ctx context.Context, collection *mongo.Collection, model interface{}, query bson.M) (int64, error) {
	updateQuery := bson.M{
		"$set": model,
	}
	result, err := collection.UpdateOne(ctx, query, updateQuery)
	if err != nil {
		return 0, err
	}
	if result.ModifiedCount > 0 {
		return result.ModifiedCount, err
	} else if result.UpsertedCount > 0 {
		return result.UpsertedCount, err
	} else {
		return result.MatchedCount, err
	}
}
func PatchByIdAndVersion(ctx context.Context, collection *mongo.Collection, model map[string]interface{}, maps map[string]string, idName string, versionField string) (int64, error) {
	idQuery := BuildQueryByIdFromMap(model, idName)
	versionQuery := BuildIdAndVersionQueryByMap(idQuery, model, maps, versionField)
	b := MapToBson(model, maps)
	rowAffect, er1 := PatchOne(ctx, collection, b, versionQuery)
	if er1 != nil {
		return 0, er1
	}
	if rowAffect == 0 {
		isExist, er2 := Exist(ctx, collection, idQuery["_id"], false)
		if er2 != nil {
			return 0, er2
		}
		if isExist {
			return -1, nil
		}
		return 0, nil
	}
	return rowAffect, er1
}
func BuildIdAndVersionQueryByMap(query map[string]interface{}, v map[string]interface{}, maps map[string]string, versionField string) map[string]interface{} {
	newMap := copyMap(query)
	if currentVersion, exist := v[versionField]; exist {
		newMap[maps[versionField]] = currentVersion
		switch versionValue := currentVersion.(type) {
		case int:
			{
				v[versionField] = versionValue + 1
			}
		case int32:
			{
				v[versionField] = versionValue + 1
			}
		case int64:
			{
				v[versionField] = versionValue + 1
			}
		default:
			panic("not support type's version")
		}
	}
	return newMap
}
func Upsert(ctx context.Context, collection *mongo.Collection, model interface{}, fieldname string) error {
	query := BuildQueryId(model, fieldname)
	_, err := UpsertOne(ctx, collection, query, model)
	return err
}
func UpsertOne(ctx context.Context, collection *mongo.Collection, filter bson.M, model interface{}) (int64, error) {
	defaultObjID, _ := primitive.ObjectIDFromHex("000000000000")

	if idValue := filter["_id"]; idValue == "" || idValue == 0 || idValue == defaultObjID {
		return InsertOne(ctx, collection, model)
	} else {
		isExisted, err := Exist(ctx, collection, idValue, false)
		if err != nil {
			return 0, err
		}
		if isExisted {
			update := bson.M{
				"$set": model,
			}
			result := collection.FindOneAndUpdate(ctx, filter, update)
			if result.Err() != nil {
				if fmt.Sprint(result.Err()) == "mongo: no documents in result" {
					return 0, nil
				} else {
					return 0, result.Err()
				}
			}
			return 1, result.Err()
		} else {
			return InsertOne(ctx, collection, model)
		}
	}
}
func UpsertOneWithVersion(ctx context.Context, collection *mongo.Collection, model interface{}, versionIndex int) (int64, error) {
	idQuery := BuildQueryByIdFromObject(model)
	defaultObjID, _ := primitive.ObjectIDFromHex("000000000000")

	if idValue := idQuery["_id"]; idValue == "" || idValue == 0 || idValue == defaultObjID {
		return InsertOneWithVersion(ctx, collection, model, versionIndex)
	} else {
		isExisted, err := Exist(ctx, collection, idValue, false)
		if err != nil {
			return 0, err
		}
		if isExisted {
			versionQuery := BuildIdAndVersionQueryByVersionIndex(idQuery, model, versionIndex)
			update := bson.M{
				"$set": model,
			}
			result := collection.FindOneAndUpdate(ctx, versionQuery, update)
			if result.Err() != nil {
				if fmt.Sprint(result.Err()) == "mongo: no documents in result" {
					return -1, nil
				} else {
					return 0, result.Err()
				}
			}
			return 1, result.Err()
		} else {
			return InsertOneWithVersion(ctx, collection, model, versionIndex)
		}
	}
}
func DeleteOne(ctx context.Context, coll *mongo.Collection, query bson.M) (int64, error) {
	result, err := coll.DeleteOne(ctx, query)
	if result == nil {
		return 0, err
	}
	return result.DeletedCount, err
}
