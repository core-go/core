package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"
)

type BatchUpdateService struct {
	collection *mongo.Collection
	modelType  reflect.Type
	idName     string
	idObjectId bool
}

func NewBatchUpdateService(db *mongo.Database, collection string, modelType reflect.Type, idObjectId bool) *BatchUpdateService {
	_, idName, _ := FindIdField(modelType)
	if len(idName) == 0 {
		log.Println(modelType.Name() + " repository can't use functions that need Id value (Ex GetById, ExistsById, Save, Update) because don't have any fields of " + modelType.Name() + " struct define _id bson tag.")
	}
	return &BatchUpdateService{db.Collection(collection), modelType, idName, idObjectId}
}

func NewDefaultBatchUpdateService(db *mongo.Database, collection string, modelType reflect.Type) *BatchUpdateService {
	return NewBatchUpdateService(db, collection, modelType, false)
}

func (m *BatchUpdateService) InsertMany(ctx context.Context, models interface{}) (int64, bool, error) {
	objects, _ := MapToMongoObjects(models, m.idName, m.idObjectId, m.modelType, true)
	values := reflect.ValueOf(models)
	length := int64(values.Len())
	duplicate, err := InsertMany(ctx, m.collection, objects)
	return length, duplicate, err
}

func (m *BatchUpdateService) UpdateMany(ctx context.Context, models interface{}) (int64, error) {
	objects, _ := MapToMongoObjects(models, m.idName, m.idObjectId, m.modelType, false)
	rs, err := UpdateMany(ctx, m.collection, objects, m.idName)
	if err != nil {
		return 0, err
	}
	return rs.ModifiedCount + rs.UpsertedCount + rs.MatchedCount, err
}

func (m *BatchUpdateService) SaveMany(ctx context.Context, models interface{}) (int64, error) {
	objects, _ := MapToMongoObjects(models, m.idName, m.idObjectId, m.modelType, false)
	rs, err := UpsertMany(ctx, m.collection, objects, m.idName)
	if err != nil {
		return 0, err
	}
	return rs.InsertedCount + rs.ModifiedCount + rs.UpsertedCount + rs.MatchedCount, err
}
