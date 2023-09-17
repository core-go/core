package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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
func FindOne(ctx context.Context, collection *firestore.CollectionRef, docID string, modelType reflect.Type) (interface{}, error) {
	idx, _, _ := FindIdField(modelType)
	if idx < 0 {
		log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex GetById, ExistsById, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
	}
	return FindOneWithIdIndexAndTracking(ctx, collection, docID, modelType, idx, -1, -1)
}
func FindOneWithIdIndex(ctx context.Context, collection *firestore.CollectionRef, docID string, modelType reflect.Type, idIndex int) (interface{}, error) {
	return FindOneWithIdIndexAndTracking(ctx, collection, docID, modelType, idIndex, -1, -1)
}
func FindOneWithTracking(ctx context.Context, collection *firestore.CollectionRef, docID string, modelType reflect.Type, createdTimeIndex int, updatedTimeIndex int) (interface{}, error) {
	idx, _, _ := FindIdField(modelType)
	if idx < 0 {
		log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex GetById, ExistsById, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
	}
	return FindOneWithIdIndexAndTracking(ctx, collection, docID, modelType, idx, createdTimeIndex, updatedTimeIndex)
}
func FindOneAndDecode(ctx context.Context, collection *firestore.CollectionRef, docID string, result interface{}) (bool, error) {
	modelType := reflect.Indirect(reflect.ValueOf(result)).Type()
	//modelType := reflect.TypeOf(result)
	idx, _, _ := FindIdField(modelType)
	if idx < 0 {
		log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex GetById, ExistsById, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
	}
	return FindOneAndDecodeWithIdIndexAndTracking(ctx, collection, docID, result, idx, -1, -1)
}
func FindOneAndDecodeWithIdIndex(ctx context.Context, collection *firestore.CollectionRef, docID string, result interface{}, idIndex int) (interface{}, error) {
	return FindOneAndDecodeWithIdIndexAndTracking(ctx, collection, docID, result, idIndex, -1, -1)
}
func FindOneAndDecodeWithTracking(ctx context.Context, collection *firestore.CollectionRef, docID string, result interface{}, createdTimeIndex int, updatedTimeIndex int) (interface{}, error) {
	modelType := reflect.TypeOf(result)
	idx, _, _ := FindIdField(modelType)
	if idx < 0 {
		log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex GetById, ExistsById, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
	}
	return FindOneAndDecodeWithIdIndexAndTracking(ctx, collection, docID, modelType, idx, createdTimeIndex, updatedTimeIndex)
}
func FindOneMap(ctx context.Context, collection *firestore.CollectionRef, docID string) (bool, map[string]interface{}, error) {
	doc, err := collection.Doc(docID).Get(ctx)
	if err != nil {
		return false, nil, err
	}
	return true, doc.Data(), nil
}
func FindOneDoc(ctx context.Context, collection *firestore.CollectionRef, docID string) (bool, *firestore.DocumentSnapshot, error) {
	doc, err := collection.Doc(docID).Get(ctx)
	if err != nil {
		return false, nil, err
	}
	return true, doc, nil
}
func FindOneWithQueries(ctx context.Context, collection *firestore.CollectionRef, where []Query, modelType reflect.Type, createdTimeIndex int, updatedTimeIndex int) (interface{}, error) {
	return FindOneWithQueriesAndTracking(ctx, collection, where, modelType, -1, -1)
}
func FindOneWithQueriesAndTracking(ctx context.Context, collection *firestore.CollectionRef, where []Query, modelType reflect.Type, createdTimeIndex int, updatedTimeIndex int) (interface{}, error) {
	iter := GetDocuments(ctx, collection, where, 1)
	idx, _, _ := FindIdField(modelType)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		result := reflect.New(modelType).Interface()
		err = doc.DataTo(&result)
		if err != nil {
			return nil, err
		}
		BindCommonFields(result, doc, idx, createdTimeIndex, updatedTimeIndex)
		return result, nil
	}
	return nil, status.Errorf(codes.NotFound, "not found")
}
func FindByField(ctx context.Context, collection *firestore.CollectionRef, values []string, modelType reflect.Type, jsonName string) (interface{}, []error) {
	return FindByFieldWithTracking(ctx, collection, values, modelType, jsonName, -1, -1)
}
func FindByFieldWithTracking(ctx context.Context, collection *firestore.CollectionRef, values []string, modelType reflect.Type, jsonName string, createdTimeIndex int, updatedTimeIndex int) (interface{}, []error) {
	idx, _, firestoreField := GetFieldByJson(modelType, jsonName)
	iter := collection.Where(firestoreField, "in", values).Documents(ctx)
	var result []interface{}
	var failure []error
	var keySuccess []string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			failure = append(failure, err)
			break
		}
		data := reflect.New(modelType).Interface()
		err = doc.DataTo(&data)
		if err != nil {
			failure = append(failure, err)
			break
		}
		BindCommonFields(data, doc, idx, createdTimeIndex, updatedTimeIndex)
		keySuccess = append(keySuccess, doc.Ref.ID)
		result = append(result, data)
	}
	// keyFailure := difference(keySuccess, values)
	return result, failure
}
func FindAndDecode(ctx context.Context, collection *firestore.CollectionRef, ids []string, result interface{}, jsonField string) ([]string, []string, []error) {
	return FindAndDecodeWithTracking(ctx, collection, ids, result, jsonField, -1, -1)
}
func FindAndDecodeWithTracking(ctx context.Context, collection *firestore.CollectionRef, ids []string, result interface{}, jsonField string, createdTimeIndex int, updatedTimeIndex int) ([]string, []string, []error) {
	var failure []error
	var keySuccess []string
	var keyFailure []string
	if reflect.TypeOf(result).Kind() != reflect.Slice {
		failure = append(failure, errors.New("result must be a slice"))
		return keySuccess, keyFailure, failure
	}
	modelType := reflect.TypeOf(result).Elem()
	idx, _, firestoreField := GetFieldByJson(modelType, jsonField)
	iter := collection.Where(firestoreField, "in", ids).Documents(ctx)
	data := reflect.New(modelType).Interface()

	var sliceData []interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			failure = append(failure, err)
			break
		}
		doc.DataTo(&data)
		if err != nil {
			failure = append(failure, err)
			keyFailure = append(keyFailure, doc.Ref.ID)
			break
		}
		BindCommonFields(data, doc, idx, createdTimeIndex, updatedTimeIndex)
		keySuccess = append(keySuccess, doc.Ref.ID)
		sliceData = append(sliceData, data)
	}
	valueResult := reflect.ValueOf(result)
	valueData := reflect.ValueOf(sliceData)
	reflect.Copy(valueResult, valueData)
	result = valueResult.Interface()
	// keyFailure := difference(keySuccess, ids)
	return keySuccess, keyFailure, failure
}

// Update
func BuildQueryByIdFromObject(object interface{}, modelType reflect.Type, idIndex int) (query []Query) {
	value := reflect.Indirect(reflect.ValueOf(object)).Field(idIndex).Interface()
	return BuildQueryById(value)
}

func BuildQueryById(id interface{}) (query []Query) {
	query = []Query{{Path: "_id", Operator: "==", Value: id.(string)}}
	return query
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
func SaveOne(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}) (int64, error) {
	oldModel := reflect.New(reflect.TypeOf(model))
	if len(id) == 0 {
		return InsertOne(ctx, collection, id, model)
	}
	exist, err := FindOneAndDecode(ctx, collection, id, &oldModel)
	if err != nil {
		if errNotFound := strings.Contains(err.Error(), "not found"); !errNotFound {
			return 0, err
		}
	}
	if exist {
		return UpdateOne(ctx, collection, id, model)
	} else {
		return InsertOne(ctx, collection, id, model)
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

func getIdValueFromModel(model interface{}, idIndex int) string {
	if id, exist := getFieldValueAtIndex(model, idIndex).(string); exist {
		return id
	}
	return ""
}

func getIdValueFromMap(m map[string]interface{}) string {
	if id, exist := m["id"].(string); exist {
		return id
	}
	return ""
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
