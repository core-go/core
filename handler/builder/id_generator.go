package builder

import (
	"context"
	"reflect"
	"strings"
)

type IdGenerator[T any] struct {
	GenerateId func(ctx context.Context) (string, error)
	idTogether bool
	emptyOnly  bool
}

func NewIdGenerator[T any](generate func(context.Context) (string, error), options ...bool) *IdGenerator[T] {
	var idTogether, emptyOnly bool
	if len(options) > 0 {
		idTogether = options[0]
	} else {
		idTogether = true
	}
	if len(options) > 1 {
		emptyOnly = options[1]
	} else {
		emptyOnly = true
	}
	x := IdGenerator[T]{GenerateId: generate, idTogether: idTogether, emptyOnly: emptyOnly}
	return &x
}

func (s *IdGenerator[T]) Generate(ctx context.Context, model *T) (int, error) {
	valueObject := reflect.Indirect(reflect.ValueOf(model))
	if valueObject.Kind() == reflect.Ptr {
		valueObject = reflect.Indirect(valueObject)
	}

	startPrimaryKey := -1
	modelType := valueObject.Type()
	numField := modelType.NumField()
	var idFields []string
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		bool := havePrimaryKeySql(tags)

		bsonTag := field.Tag.Get("bson")
		bsonTags := strings.Split(bsonTag, ",")
		boolMongo := havePrimaryKeyMongo(bsonTags)
		if boolMongo || bool {
			startPrimaryKey = i
			idTag := field.Tag.Get("id")
			idTags := strings.Split(idTag, ";")
			if idTags[0] == "manual" {
				continue
			}
			id, err := s.GenerateId(ctx)
			if err != nil {
				return 0, err
			}

			fieldName := reflect.Indirect(valueObject).FieldByName(field.Name)
			if fieldName.Kind() == reflect.Ptr {
				setValue(model, i, &id)
				idFields = append(idFields, id)
			} else {
				setValue(model, i, id)
				idFields = append(idFields, id)
			}
		} else {
			if s.idTogether && startPrimaryKey > -1 {
				break
			}
		}

	}
	return len(idFields), nil
}
func havePrimaryKeySql(tags []string) bool {
	for _, tag := range tags {
		if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
			return true
		}
	}
	return false
}

func havePrimaryKeyMongo(tags []string) bool {
	for _, tag := range tags {
		if strings.Compare(strings.TrimSpace(tag), "_id") == 0 {
			return true
		}
	}
	return false
}
