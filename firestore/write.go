package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"reflect"
	"strings"
	"time"
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
func MapToFirestore(json map[string]interface{}, docMap map[string]interface{}, maps map[string]string) map[string]interface{} {
	//docMap := doc.Data()
	for k, v := range json {
		fk, ok := maps[k]
		if ok {
			docMap[fk] = v
		}
	}
	return docMap
}
func Create(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}) (int64, string, *time.Time, error) {
	var docRef *firestore.DocumentRef
	rid := id
	// TODO apply idField.IsZero() for golang 13 or above
	if len(id) > 0 {
		docRef = collection.Doc(id)
	} else {
		docRef = collection.NewDoc()
		rid = docRef.ID
	}
	res, err := docRef.Create(ctx, model)

	if err != nil {
		if strings.Contains(err.Error(), "Document already exists") {
			return 0, rid, nil, err
		} else {
			return -1, rid, nil, err
		}
	}
	return 1, rid, &res.UpdateTime, nil
}
func Save(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}) (int64, *time.Time, error) {
	res, err := collection.Doc(id).Set(ctx, model)
	if err != nil {
		return -1, nil, err
	}
	return 1, &res.UpdateTime, nil
}
func Update(ctx context.Context, collection *firestore.CollectionRef, id string, model interface{}) (int64, *time.Time, error) {
	docRef := collection.Doc(id)
	_, er0 := docRef.Get(ctx)
	if er0 != nil {
		if strings.HasSuffix(er0.Error(), " not found") {
			return 0, nil, nil
		}
		return -1, nil, er0
	}
	res, err := docRef.Set(ctx, model)
	if err != nil {
		return -1, nil, err
	}
	return 1, &res.UpdateTime, nil
}
func Delete(ctx context.Context, collection *firestore.CollectionRef, id string) (int64, error) {
	_, err := collection.Doc(id).Delete(ctx, firestore.Exists)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return 0, nil
		}
		return 0, err
	}
	return 1, err
}
