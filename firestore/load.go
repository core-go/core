package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"strings"
)

func Exist(ctx context.Context, collection *firestore.CollectionRef, id string) (bool, error) {
	docRef := collection.Doc(id)
	_, err := docRef.Get(ctx)
	if err != nil {
		if strings.HasSuffix(err.Error(), " not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func Load(ctx context.Context, collection *firestore.CollectionRef, id string, res interface{}) (bool, *firestore.DocumentSnapshot, error) {
	docRef := collection.Doc(id)
	doc, er1 := docRef.Get(ctx)
	if er1 != nil {
		if strings.HasSuffix(er1.Error(), " not found") {
			return false, doc, nil
		}
		return false, doc, er1
	}
	er2 := doc.DataTo(res)
	return true, doc, er2
}
