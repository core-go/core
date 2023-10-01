package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strings"
)

func FindIdField(modelType reflect.Type) (int, string, string) {
	return FindField(modelType, "_id")
}
func FindField(modelType reflect.Type, bsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		bsonTag := field.Tag.Get("bson")
		tags := strings.Split(bsonTag, ",")
		json := field.Name
		if tag1, ok1 := field.Tag.Lookup("json"); ok1 {
			json = strings.Split(tag1, ",")[0]
		}
		for _, tag := range tags {
			if strings.TrimSpace(tag) == bsonName {
				return i, field.Name, json
			}
		}
	}
	return -1, "", ""
}
func Exist(ctx context.Context, collection *mongo.Collection, id interface{}, objectId bool) (bool, error) {
	query := bson.M{"_id": id}
	if objectId {
		objId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return false, err
		}
		query = bson.M{"_id": objId}
	}
	x := collection.FindOne(ctx, query)
	if x.Err() != nil {
		if fmt.Sprint(x.Err()) == "mongo: no documents in result" {
			return false, nil
		} else {
			return false, x.Err()
		}
	}
	return true, nil
}
