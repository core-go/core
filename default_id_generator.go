package service

import (
	"context"
	"reflect"
	"strings"
)

type DefaultIdGenerator struct {
	GenerateId func() (string, error)
	idTogether bool
	emptyOnly  bool
}

func NewDefaultIdGenerator(generate func() (string, error)) *DefaultIdGenerator {
	return NewIdGenerator(generate, true, true)
}

func NewIdGenerator(generate func() (string, error), idTogether bool, emptyOnly bool) *DefaultIdGenerator {
	x := DefaultIdGenerator{GenerateId: generate, idTogether: idTogether, emptyOnly: emptyOnly}
	return &x
}

func (s *DefaultIdGenerator) Generate(ctx context.Context, model interface{}) (int, error) {
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
			id, err := s.GenerateId()
			if err != nil {
				return 0, err
			}

			fieldName := reflect.Indirect(valueObject).FieldByName(field.Name)
			if fieldName.Kind() == reflect.Ptr {
				SetValue(model, i, &id)
				idFields = append(idFields, id)
			} else {
				SetValue(model, i, id)
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
