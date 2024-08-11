package firestore

import (
	"cloud.google.com/go/firestore"
	"reflect"
	"strings"
)

func BindCommonFields(res interface{}, doc *firestore.DocumentSnapshot, idIndex int, createdTimeIndex int, updatedTimeIndex int) {
	rv := reflect.Indirect(reflect.ValueOf(res))
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
