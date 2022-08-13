package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type Writer struct {
	*Loader
	maps         map[string]string
	versionField string
	versionIndex int
	Mapper       Mapper
}

func NewMongoWriterWithVersion(db *mongo.Database, collectionName string, modelType reflect.Type, idObjectId bool, versionField string, options ...Mapper) *Writer {
	var mapper Mapper
	var loader *Loader
	if len(options) > 0 {
		mapper = options[0]
		loader = NewMongoLoader(db, collectionName, modelType, idObjectId, mapper.DbToModel)
	} else {
		loader = NewMongoLoader(db, collectionName, modelType, idObjectId)
	}
	if len(versionField) > 0 {
		index := FindFieldIndex(modelType, versionField)
		if index >= 0 {
			return &Writer{Loader: loader, maps: MakeBsonMap(modelType), versionField: versionField, versionIndex: index, Mapper: mapper}
		}
	}
	return &Writer{Loader: loader, maps: MakeBsonMap(modelType), versionField: "", versionIndex: -1, Mapper: mapper}
}
func NewWriterWithVersion(db *mongo.Database, collectionName string, modelType reflect.Type, versionField string, options ...Mapper) *Writer {
	return NewMongoWriterWithVersion(db, collectionName, modelType, false, versionField, options...)
}
func NewWriter(db *mongo.Database, collectionName string, modelType reflect.Type, options ...Mapper) *Writer {
	return NewMongoWriterWithVersion(db, collectionName, modelType, false, "", options...)
}

func (m *Writer) Insert(ctx context.Context, model interface{}) (int64, error) {
	if m.Mapper != nil {
		m2, err := m.Mapper.ModelToDb(ctx, model)
		if err != nil {
			return 0, err
		}
		if m.versionIndex >= 0 {
			return InsertOneWithVersion(ctx, m.Collection, m2, m.versionIndex)
		}
		return InsertOne(ctx, m.Collection, m2)
	}
	if m.versionIndex >= 0 {
		return InsertOneWithVersion(ctx, m.Collection, model, m.versionIndex)
	}
	return InsertOne(ctx, m.Collection, model)
}

func (m *Writer) Update(ctx context.Context, model interface{}) (int64, error) {
	if m.Mapper != nil {
		m2, err := m.Mapper.ModelToDb(ctx, model)
		if err != nil {
			return 0, err
		}
		if m.versionIndex >= 0 {
			return UpdateByIdAndVersion(ctx, m.Collection, m2, m.versionIndex)
		}
		idQuery := BuildQueryByIdFromObject(m2)
		return UpdateOne(ctx, m.Collection, m2, idQuery)
	}
	if m.versionIndex >= 0 {
		return UpdateByIdAndVersion(ctx, m.Collection, model, m.versionIndex)
	}
	idQuery := BuildQueryByIdFromObject(model)
	return UpdateOne(ctx, m.Collection, model, idQuery)
}

func (m *Writer) Patch(ctx context.Context, model map[string]interface{}) (int64, error) {
	if m.Mapper != nil {
		m2, err := m.Mapper.ModelToDb(ctx, model)
		if err != nil {
			return 0, err
		}
		m3, ok1 := m2.(map[string]interface{})
		if !ok1 {
			return 0, fmt.Errorf("result of LocationToBson must be a map[string]interface{}")
		}
		if m.versionIndex >= 0 {
			return PatchByIdAndVersion(ctx, m.Collection, m3, m.maps, m.jsonIdName, m.versionField)
		}
		jsonName0 := GetJsonByIndex(m.modelType, m.idIndex)
		idQuery := BuildQueryByIdFromMap(m3, jsonName0)
		b0 := MapToBson(m3, m.maps)
		return PatchOne(ctx, m.Collection, b0, idQuery)
	}
	if m.versionIndex >= 0 {
		return PatchByIdAndVersion(ctx, m.Collection, model, m.maps, m.jsonIdName, m.versionField)
	}
	jsonName := GetJsonByIndex(m.modelType, m.idIndex)
	idQuery := BuildQueryByIdFromMap(model, jsonName)
	b := MapToBson(model, m.maps)
	return PatchOne(ctx, m.Collection, b, idQuery)
}

func (m *Writer) Save(ctx context.Context, model interface{}) (int64, error) {
	if m.Mapper != nil {
		m2, err := m.Mapper.ModelToDb(ctx, model)
		if err != nil {
			return 0, err
		}
		if m.versionIndex >= 0 {
			return UpsertOneWithVersion(ctx, m.Collection, m2, m.versionIndex)
		}
		idQuery := BuildQueryByIdFromObject(m2)
		return UpsertOne(ctx, m.Collection, idQuery, m2)
	}
	if m.versionIndex >= 0 {
		return UpsertOneWithVersion(ctx, m.Collection, model, m.versionIndex)
	}
	idQuery := BuildQueryByIdFromObject(model)
	return UpsertOne(ctx, m.Collection, idQuery, model)
}

func (m *Writer) Delete(ctx context.Context, id interface{}) (int64, error) {
	query := bson.M{"_id": id}
	return DeleteOne(ctx, m.Collection, query)
}
