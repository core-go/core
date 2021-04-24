package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FieldLoader struct {
	Collection *mongo.Collection
	Name       string
}

func NewFieldLoader(db *mongo.Database, collectionName string, name string) *FieldLoader {
	collection := db.Collection(collectionName)
	return &FieldLoader{
		Collection: collection,
		Name:       name,
	}
}

func (l *FieldLoader) Values(ctx context.Context, ids []string) ([]string, error) {
	var array []string
	var finalResult []bson.M
	query := bson.M{l.Name: bson.M{"$in": ids}}

	findOptions := options.Find() // build a `findOptions`
	findOptions.SetSort(map[string]int{l.Name: 1})
	findOptions.SetProjection(map[string]int{l.Name: 1, "_id": 0})
	result, err := l.Collection.Find(ctx, query, findOptions)
	result.All(ctx, &finalResult)

	if err != nil {
		return array, err
	}

	for _, model := range finalResult {
		array = append(array, model[l.Name].(string))
	}
	return array, nil
}
