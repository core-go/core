package service

import (
	"reflect"
	"strings"
)

type DefaultIdGenerator struct {
	idTogether bool
	shortId    bool
}

func NewDefaultIdGenerator() *DefaultIdGenerator {
	generator := DefaultIdGenerator{true, true}
	return &generator
}

func NewIdGenerator(idTogether bool, emptyOnly bool, shortId bool) *DefaultIdGenerator {
	generator := DefaultIdGenerator{idTogether, shortId}
	return &generator
}

func (s *DefaultIdGenerator) Generate(model interface{}) (int, error) {
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
		if bool == false {
			bsonTag := field.Tag.Get("bson")
			tags := strings.Split(bsonTag, ",")
			bool := havePrimaryKeyMongo(tags)
			if bool {
				startPrimaryKey = i
				idTag := field.Tag.Get("id")
				idTags := strings.Split(idTag, ";")
				if idTags[0] == "manual" {
					continue
				}
				var id string
				if s.shortId == true {
					shortId, err := ShortId()
					if err != nil {
						return 0, err
					}
					id = shortId
				} else {
					id = RandomId()
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
