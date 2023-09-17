package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"reflect"
	"strings"
)

func Connect(ctx context.Context, credentials []byte) (*firestore.Client, error) {
	app, er1 := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(credentials))
	if er1 != nil {
		log.Fatalf("Could not create admin client: %v", er1)
		return nil, er1
	}

	client, er2 := app.Firestore(ctx)
	if er2 != nil {
		log.Fatalf("Could not create data operations client: %v", er2)
		return nil, er2
	}
	return client, nil
}
func BindCommonFields(result interface{}, doc *firestore.DocumentSnapshot, idIndex int, createdTimeIndex int, updatedTimeIndex int) {
	rv := reflect.Indirect(reflect.ValueOf(result))
	fv := rv.Field(idIndex)
	fv.Set(reflect.ValueOf(doc.Ref.ID))

	if createdTimeIndex >= 0 {
		cv := rv.Field(createdTimeIndex)
		cv.Set(reflect.ValueOf(&doc.CreateTime))
	}
	if updatedTimeIndex >= 0 {
		uv := rv.Field(updatedTimeIndex)
		uv.Set(reflect.ValueOf(&doc.UpdateTime))
	}
}
func FindFieldByName(modelType reflect.Type, fieldName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			name1 := fieldName
			name2 := fieldName
			tag1, ok1 := field.Tag.Lookup("json")
			tag2, ok2 := field.Tag.Lookup("firestore")
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
// for get all and search
func appendToArray(arr interface{}, item interface{}) interface{} {
	arrValue := reflect.ValueOf(arr)
	elemValue := arrValue.Elem()

	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = reflect.Indirect(itemValue)
	}
	elemValue.Set(reflect.Append(elemValue, itemValue))
	return arr
}
func FindIdField(modelType reflect.Type) (int, string, string) {
	return findBsonField(modelType, "_id")
}
func findBsonField(modelType reflect.Type, bsonName string) (int, string, string) {
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
func SetValue(model interface{}, index int, value interface{}) (interface{}, error) {
	v := reflect.Indirect(reflect.ValueOf(model))
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	v.Field(index).Set(reflect.ValueOf(value))
	return model, nil
}
