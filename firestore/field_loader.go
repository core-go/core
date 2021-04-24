package firestore

import (
	"context"
	"sort"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FieldLoader struct {
	Collection *firestore.CollectionRef
	Name       string
}

func NewFieldLoader(client *firestore.Client, collectionName string, name string) *FieldLoader {
	collection := client.Collection(collectionName)
	return &FieldLoader{
		Collection: collection,
		Name:       name,
	}
}

func (l *FieldLoader) Values(ctx context.Context, ids []string) ([]string, error) {
	var array []string
	iter := l.Collection.Select(l.Name).Where(l.Name, "in", ids).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var model map[string]interface{}
		err = doc.DataTo(&model)
		if err != nil {
			return nil, err
		}
		array = append(array, model[l.Name].(string))
	}
	sort.Strings(array)
	return array, nil
}
