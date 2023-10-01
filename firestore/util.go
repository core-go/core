package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"reflect"
	"strings"
)

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}
func FindField(modelType reflect.Type, firestoreName string) (int, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == firestoreName {
			firestoreTag := field.Tag.Get("firestore")
			tags := strings.Split(firestoreTag, ",")
			for _, tag := range tags {
				if strings.Compare(strings.TrimSpace(tag), firestoreName) == 0 {
					return i, field.Name
				}
			}
		}
	}
	return -1, ""
}
func GetFirestoreName(modelType reflect.Type, fieldName string) string {
	field, _ := modelType.FieldByName(fieldName)
	bsonTag := field.Tag.Get("firestore")
	tags := strings.Split(bsonTag, ",")
	if len(tags) > 0 {
		return tags[0]
	}
	return fieldName
}
func findId(queries []Query) string {
	for _, p := range queries {
		if p.Path == "_id" || p.Path == "" {
			return p.Value.(string)
		}
	}
	return ""
}
func MapFieldId(value interface{}, fieldNameId string, doc *firestore.DocumentSnapshot) {
	// fmt.Println(reflect.TypeOf(value))
	rv := reflect.Indirect(reflect.ValueOf(value))
	fv := rv.FieldByName(fieldNameId)
	if fv.IsValid() && fv.CanAddr() { //TODO handle set , now error no set id
		fv.Set(reflect.ValueOf(doc.Ref.ID))
	}
}
func FindFieldName(modelType reflect.Type, firestoreName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		bsonTag := field.Tag.Get("firestore")
		tags := strings.Split(bsonTag, ",")
		json := field.Name
		if tag1, ok1 := field.Tag.Lookup("json"); ok1 {
			json = strings.Split(tag1, ",")[0]
		}
		for _, tag := range tags {
			if strings.TrimSpace(tag) == firestoreName {
				return i, field.Name, json
			}
		}
	}
	return -1, "", ""
}

// Update
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
func getIdValueFromMap(m map[string]interface{}) string {
	if id, exist := m["id"].(string); exist {
		return id
	}
	return ""
}
