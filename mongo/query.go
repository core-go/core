package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

func BuildSearchResult(ctx context.Context, collection *mongo.Collection, results interface{}, query bson.D, fields bson.M, sort bson.D, limit int64, skip int64, opts ...func(context.Context, interface{}) (interface{}, error)) (int64, error) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(opts) > 0 {
		mp = opts[0]
	}
	optionsFind := options.Find()
	optionsFind.Projection = fields
	optionsFind.SetSkip(skip)
	optionsFind.SetLimit(limit)
	if sort != nil {
		optionsFind.SetSort(sort)
	}

	cursor, er0 := collection.Find(ctx, query, optionsFind)
	if er0 != nil {
		return 0, er0
	}

	er1 := cursor.All(ctx, results)
	if er1 != nil {
		return 0, er1
	}
	options := options.Count()
	count, er2 := collection.CountDocuments(ctx, query, options)
	if er2 != nil {
		return 0, er2
	}
	if mp == nil {
		return count, nil
	}
	_, er3 := MapModels(ctx, results, mp)
	return count, er3
}

func BuildSort(s string, modelType reflect.Type) bson.D {
	var sort = bson.D{}
	if len(s) == 0 {
		return sort
	}
	sorts := strings.Split(s, ",")
	for i := 0; i < len(sorts); i++ {
		sortField := strings.TrimSpace(sorts[i])
		if len(sortField) > 0 {
			fieldName := sortField
			c := sortField[0:1]
			if c == "-" || c == "+" {
				fieldName = sortField[1:]
			}
			columnName := GetBsonNameForSort(modelType, fieldName)
			if len(columnName) > 0 {
				sortType := GetSortType(c)
				sort = append(sort, bson.E{Key: columnName, Value: sortType})
			}
		}
	}
	return sort
}

func GetBsonNameForSort(modelType reflect.Type, sortField string) string {
	sortField = strings.TrimSpace(sortField)
	idx, fieldName, name := GetFieldByJson(modelType, sortField)
	if len(name) > 0 {
		return name
	}
	if idx >= 0 {
		return fieldName
	}
	return sortField
}

func GetSortType(sortType string) int {
	if sortType == "-" {
		return -1
	} else {
		return 1
	}
}
