package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"reflect"
	"strings"
)

func MakeFirestoreMap(modelType reflect.Type) map[string]string {
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
		if tag, ok := field.Tag.Lookup("firestore"); ok {
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
func DeleteOne(ctx context.Context, collection *firestore.CollectionRef, docID string) (int64, error) {
	_, err := collection.Doc(docID).Delete(ctx, firestore.Exists)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return 0, nil
		}
		return 0, err
	}
	return 1, err
}
func InsertOne(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}) (int64, error) {
	var doc *firestore.DocumentRef
	// TODO apply idField.IsZero() for golang 13 or above
	if len(id) > 0 {
		doc = collection.NewDoc()
	} else {
		doc = collection.Doc(id)
	}
	_, err := doc.Create(ctx, model)
	if err != nil {
		if strings.Index(err.Error(), "Document already exists") >= 0 {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return 1, nil
}
func InsertOneWithVersion(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}, versionIndex int) (int64, error) {
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
	model, err := setValueWithIndex(model, versionIndex, defaultVersion)
	if err != nil {
		return 0, err
	}
	return InsertOne(ctx, collection, id, model)
}
func FindOneMap(ctx context.Context, collection *firestore.CollectionRef, docID string) (bool, map[string]interface{}, error) {
	doc, err := collection.Doc(docID).Get(ctx)
	if err != nil {
		return false, nil, err
	}
	return true, doc.Data(), nil
}
func UpdateOneWithVersion(ctx context.Context, collection *firestore.CollectionRef, model interface{}, versionIndex int, versionFieldName string, idIndex int) (int64, error) {
	id := getIdValueFromModel(model, idIndex)
	if len(id) == 0 {
		return 0, fmt.Errorf("cannot update one an Object that do not have id field")
	}
	itemExist, oldModel, err := FindOneMap(ctx, collection, id)
	if err != nil {
		return 0, err
	}
	if itemExist {
		currentVersion := getFieldValueAtIndex(model, versionIndex)
		oldVersion := oldModel[versionFieldName]
		switch reflect.TypeOf(currentVersion).String() {
		case "int":
			oldVersion = int(oldVersion.(int64))
		case "int32":
			oldVersion = int32(oldVersion.(int64))
		}
		if currentVersion == oldVersion {
			updateModelVersion(model, versionIndex)
			_, err := collection.Doc(id).Set(ctx, model)
			if err != nil {
				return 0, err
			} else {
				return 1, nil
			}
		} else {
			return -1, fmt.Errorf("wrong version")
		}
	} else {
		return 0, fmt.Errorf("not found")
	}
}
func updateModelVersion(model interface{}, versionIndex int) {
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	currentVersion := getFieldValueAtIndex(model, versionIndex)

	switch reflect.ValueOf(currentVersion).Kind().String() {
	case "int":
		nextVersion := reflect.ValueOf(currentVersion.(int) + 1)
		modelValue.Field(versionIndex).Set(nextVersion)
	case "int32":
		nextVersion := reflect.ValueOf(currentVersion.(int32) + 1)
		modelValue.Field(versionIndex).Set(nextVersion)
	case "int64":
		nextVersion := reflect.ValueOf(currentVersion.(int64) + 1)
		modelValue.Field(versionIndex).Set(nextVersion)
	default:
		panic("version's type not supported")
	}
}
func PatchOne(ctx context.Context, collection *firestore.CollectionRef, id string, json map[string]interface{}, maps map[string]string) (int64, error) {
	if len(id) == 0 {
		return 0, fmt.Errorf("cannot patch one an Object that do not have id field")
	}
	docRef := collection.Doc(id)
	doc, er1 := docRef.Get(ctx)
	if er1 != nil{
		return -1, er1
	}
	fs := MapToFirestore(json, doc, maps)
	_, er2 := docRef.Set(ctx, fs)
	if er2 != nil {
		return 0, er2
	}
	return 1, nil
}
func PatchOneWithVersion(ctx context.Context, collection *firestore.CollectionRef, id string, json map[string]interface{}, maps map[string]string, jsonVersion string) (int64, error) {
	itemExist, doc, err := FindOneDoc(ctx, collection, id)
	if err != nil {
		return 0, err
	}
	if itemExist {
		fs := MapToFirestore(json, doc, maps)
		currentVersion := json[jsonVersion]
		firestoreVersion, ok := maps[jsonVersion]
		if !ok {
			return -1, fmt.Errorf("cannot map version between json and firestore")
		}
		oldVersion := fs[firestoreVersion]
		switch currentVersion.(type) {
		case int:
			oldVersion = int(oldVersion.(int64))
		case int32:
			oldVersion = int32(oldVersion.(int64))
		}
		if currentVersion == oldVersion {
			updateMapVersion(fs, firestoreVersion)
			_, err := collection.Doc(id).Set(ctx, fs)
			if err != nil {
				return 0, err
			} else {
				return 1, nil
			}
		} else {
			return -1, fmt.Errorf("wrong version")
		}
	} else {
		return 0, fmt.Errorf("not found")
	}
}
func FindOneDoc(ctx context.Context, collection *firestore.CollectionRef, docID string) (bool, *firestore.DocumentSnapshot, error) {
	doc, err := collection.Doc(docID).Get(ctx)
	if err != nil {
		return false, nil, err
	}
	return true, doc, nil
}
func MapToFirestore(json map[string]interface{}, doc *firestore.DocumentSnapshot, maps map[string]string) map[string]interface{} {
	fs := doc.Data()
	for k, v := range json {
		fk, ok := maps[k]
		if ok {
			fs[fk] = v
		}
	}
	return fs
}
func updateMapVersion(m map[string]interface{}, version string) {
	if currentVersion, exist := m[version]; exist {
		switch currentVersion.(type) {
		case int:
			m[version] = currentVersion.(int) + 1
		case int32:
			m[version] = currentVersion.(int32) + 1
		case int64:
			m[version] = currentVersion.(int64) + 1
		default:
			panic("version's type not supported")
		}
	}
}
func SaveOneWithVersion(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}, versionIndex int, versionFieldName string) (int64, error) {
	exist, oldModel, err := FindOneMap(ctx, collection, id)
	if err != nil {
		if errNotFound := strings.Contains(err.Error(), "not found"); !errNotFound {
			return 0, err
		}
	}
	if exist {
		currentVersion := getFieldValueAtIndex(model, versionIndex)
		oldVersion := oldModel[versionFieldName]
		switch reflect.TypeOf(currentVersion).String() {
		case "int":
			oldVersion = int(oldVersion.(int64))
		case "int32":
			oldVersion = int32(oldVersion.(int64))
		}
		if currentVersion == oldVersion {
			updateModelVersion(model, versionIndex)
			_, err := collection.Doc(id).Set(ctx, model)
			if err != nil {
				return 0, err
			} else {
				return 1, nil
			}
		} else {
			return -1, fmt.Errorf("wrong version")
		}
	} else {
		return InsertOneWithVersion(ctx, collection, id, model, versionIndex)
	}
}
