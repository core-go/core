package query

import (
	"reflect"
	"strings"

	"github.com/core-go/search"
)

func UseQuery[T any, F any]() func(F) map[string]interface{} {
	b := NewBuilder[T, F]()
	return b.BuildQuery
}

type Builder[T any, F any] struct {
	ModelType reflect.Type
}

func NewBuilder[T any, F any]() *Builder[T, F] {
	var t T
	modelType := reflect.TypeOf(t)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	return &Builder[T, F]{ModelType: modelType}
}
func (b *Builder[T, F]) BuildQuery(filter F) map[string]interface{} {
	return Build(filter, b.ModelType)
}

func Build(filter interface{}, resultModelType reflect.Type) map[string]interface{} {
	query := map[string]interface{}{}
	if _, ok := filter.(*search.Filter); ok {
		return query
	}
	value := reflect.Indirect(reflect.ValueOf(filter))
	numField := value.NumField()
	for i := 0; i < numField; i++ {
		fieldValue := value.Field(i).Interface()
		if v, ok := fieldValue.(*search.Filter); ok {
			if v.Excluding != nil && len(v.Excluding) > 0 {
				_, _, columnName := getFieldByBson(value.Type(), "_id")
				if len(columnName) > 0 {
					actionDateQuery := map[string]interface{}{}
					actionDateQuery["$nin"] = v.Excluding
					query[columnName] = actionDateQuery
				}
			}
			continue
		} else if rangeTime, ok := fieldValue.(*search.TimeRange); ok && rangeTime != nil {
			_, columnName := findFieldByName(resultModelType, value.Type().Field(i).Name)
			actionDateQuery := map[string]interface{}{}
			if rangeTime.Min != nil {
				actionDateQuery["$gte"] = *rangeTime.Min
			}
			if rangeTime.Max != nil {
				actionDateQuery["$lte"] = *rangeTime.Max
			} else if rangeTime.Top != nil {
				actionDateQuery["$lt"] = rangeTime.Top
			}
			if len(actionDateQuery) > 0 {
				query[columnName] = actionDateQuery
			}
		} else if numberRange, ok := fieldValue.(*search.NumberRange); ok && numberRange != nil {
			_, columnName := findFieldByName(resultModelType, value.Type().Field(i).Name)
			amountQuery := map[string]interface{}{}

			if numberRange.Min != nil {
				amountQuery["$gte"] = *numberRange.Min
			} else if numberRange.Lower != nil {
				amountQuery["$gt"] = *numberRange.Lower
			}
			if numberRange.Max != nil {
				amountQuery["$lte"] = *numberRange.Max
			} else if numberRange.Top != nil {
				amountQuery["$lt"] = *numberRange.Top
			}

			if len(amountQuery) > 0 {
				query[columnName] = amountQuery
			}
		} else if numberRange, ok := fieldValue.(*search.Int64Range); ok && numberRange != nil {
			_, columnName := findFieldByName(resultModelType, value.Type().Field(i).Name)
			amountQuery := map[string]interface{}{}

			if numberRange.Min != nil {
				amountQuery["$gte"] = *numberRange.Min
			} else if numberRange.Lower != nil {
				amountQuery["$gt"] = *numberRange.Lower
			}
			if numberRange.Max != nil {
				amountQuery["$lte"] = *numberRange.Max
			} else if numberRange.Top != nil {
				amountQuery["$lt"] = *numberRange.Top
			}
			if len(amountQuery) > 0 {
				query[columnName] = amountQuery
			}
		} else if numberRange, ok := fieldValue.(*search.Int32Range); ok && numberRange != nil {
			_, columnName := findFieldByName(resultModelType, value.Type().Field(i).Name)
			amountQuery := map[string]interface{}{}

			if numberRange.Min != nil {
				amountQuery["$gte"] = *numberRange.Min
			} else if numberRange.Lower != nil {
				amountQuery["$gt"] = *numberRange.Lower
			}
			if numberRange.Max != nil {
				amountQuery["$lte"] = *numberRange.Max
			} else if numberRange.Top != nil {
				amountQuery["$lt"] = *numberRange.Top
			}
			if len(amountQuery) > 0 {
				query[columnName] = amountQuery
			}
		} else if numberRange, ok := fieldValue.(*search.IntRange); ok && numberRange != nil {
			_, columnName := findFieldByName(resultModelType, value.Type().Field(i).Name)
			amountQuery := map[string]interface{}{}

			if numberRange.Min != nil {
				amountQuery["$gte"] = *numberRange.Min
			} else if numberRange.Lower != nil {
				amountQuery["$gt"] = *numberRange.Lower
			}
			if numberRange.Max != nil {
				amountQuery["$lte"] = *numberRange.Max
			} else if numberRange.Top != nil {
				amountQuery["$lt"] = *numberRange.Top
			}
			if len(amountQuery) > 0 {
				query[columnName] = amountQuery
			}
		} else if value.Field(i).Kind().String() == "slice" {
			actionDateQuery := map[string]interface{}{}
			_, columnName := findFieldByName(resultModelType, value.Type().Field(i).Name)
			actionDateQuery["$in"] = fieldValue
			query[columnName] = actionDateQuery
		} else {
			t := value.Field(i).Kind().String()
			if _, ok := fieldValue.(*search.Filter); t == "bool" || (strings.Contains(t, "int") && fieldValue != 0) || (strings.Contains(t, "float") && fieldValue != 0) || (!ok && t == "string" && value.Field(i).Len() > 0) || (!ok && t == "ptr" &&
				value.Field(i).Pointer() != 0) {
				_, columnName := findFieldByName(resultModelType, value.Type().Field(i).Name)
				if len(columnName) > 0 {
					query[columnName] = fieldValue
				}
			}
		}
	}
	return query
}

func findFieldByName(modelType reflect.Type, fieldName string) (index int, jsonTagName string) {
	numField := modelType.NumField()
	for index := 0; index < numField; index++ {
		field := modelType.Field(index)
		if field.Name == fieldName {
			jsonTagName := fieldName
			if jsonTag, ok := field.Tag.Lookup("json"); ok {
				jsonTagName = strings.Split(jsonTag, ",")[0]
			}
			return index, jsonTagName
		}
	}
	return -1, fieldName
}
func getFieldByBson(modelType reflect.Type, bsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag1, ok1 := field.Tag.Lookup("bson")
		if ok1 && strings.Split(tag1, ",")[0] == bsonName {
			if tag2, ok2 := field.Tag.Lookup("json"); ok2 {
				json := strings.Split(tag2, ",")[0]
				return i, field.Name, json
			}
			return i, field.Name, ""
		}
	}
	return -1, bsonName, bsonName
}
