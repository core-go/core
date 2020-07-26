package service

import (
	"context"
	"reflect"
	"strings"
	"time"
)

type DefaultModelBuilder struct {
	IdGenerator    IdGenerator
	modelType      reflect.Type
	createdByName  string
	createdAtName  string
	updatedByName  string
	updatedAtName  string
	createdByIndex int
	createdAtIndex int
	updatedByIndex int
	updatedAtIndex int
}

func NewModelBuilder(generator IdGenerator, modelType reflect.Type, createdByName, createdAtName, updatedByName, updatedAtName string) *DefaultModelBuilder {
	createdByIndex := FindFieldIndex(modelType, createdByName)
	createdAtIndex := FindFieldIndex(modelType, createdAtName)
	updatedByIndex := FindFieldIndex(modelType, updatedByName)
	updatedAtIndex := FindFieldIndex(modelType, updatedAtName)

	return &DefaultModelBuilder{
		IdGenerator:    generator,
		modelType:      modelType,
		createdByName:  createdByName,
		createdAtName:  createdAtName,
		updatedByName:  updatedByName,
		updatedAtName:  updatedAtName,
		createdByIndex: createdByIndex,
		createdAtIndex: createdAtIndex,
		updatedByIndex: updatedByIndex,
		updatedAtIndex: updatedAtIndex,
	}
}

func (c *DefaultModelBuilder) BuildToInsert(ctx context.Context, obj interface{}) interface{} {
	if c.IdGenerator != nil {
		c.IdGenerator.Generate(ctx, obj)
	}
	valueModelObject := reflect.Indirect(reflect.ValueOf(obj))
	if valueModelObject.Kind() == reflect.Ptr {
		valueModelObject = reflect.Indirect(valueModelObject)
	}
	userId := GetUserIdFromContext(ctx)
	if valueModelObject.Kind() == reflect.Struct {
		if c.createdByIndex >= 0 {
			createdByField := reflect.Indirect(valueModelObject).Field(c.createdByIndex)
			if createdByField.Kind() == reflect.Ptr {
				createdByField = reflect.Indirect(createdByField)
			}
			createdByField.Set(reflect.ValueOf(userId))
		}

		if c.createdAtIndex >= 0 {
			createdAtField := reflect.Indirect(valueModelObject).Field(c.createdAtIndex)
			if createdAtField.Kind() == reflect.Ptr {
				createdAtField = reflect.Indirect(createdAtField)
			}
			createdAtField.Set(reflect.ValueOf(time.Now()))
		}
	} else if valueModelObject.Kind() == reflect.Map {
		var createdByTag, createdAtTag string
		if c.createdByIndex >= 0 {
			if createdByTag = GetBsonName(c.modelType, c.createdByIndex); createdByTag == "" || createdByTag == "-" {
				createdByTag = GetJsonName(c.modelType, c.createdByIndex)
			}
			if createdByTag != "" && createdByTag != "-" {
				valueModelObject.SetMapIndex(reflect.ValueOf(createdByTag), reflect.ValueOf(userId))
			}
		}

		if c.createdAtIndex >= 0 {
			if createdAtTag = GetBsonName(c.modelType, c.createdAtIndex); createdAtTag == "" || createdAtTag == "-" {
				createdAtTag = GetJsonName(c.modelType, c.createdAtIndex)
			}
			if createdAtTag != "" && createdAtTag != "-" {
				valueModelObject.SetMapIndex(reflect.ValueOf(createdAtTag), reflect.ValueOf(time.Now()))
			}
		}
	}

	return obj
}

func (c *DefaultModelBuilder) BuildToUpdate(ctx context.Context, obj interface{}) interface{} {
	valueModelObject := reflect.Indirect(reflect.ValueOf(obj))
	if valueModelObject.Kind() == reflect.Ptr {
		valueModelObject = reflect.Indirect(valueModelObject)
	}
	userId := GetUserIdFromContext(ctx)
	if valueModelObject.Kind() == reflect.Struct {
		if c.updatedByIndex >= 0 {
			updatedByField := reflect.Indirect(valueModelObject).Field(c.updatedByIndex)
			if updatedByField.Kind() == reflect.Ptr {
				updatedByField = reflect.Indirect(updatedByField)
			}
			updatedByField.Set(reflect.ValueOf(userId))
		}

		if c.updatedAtIndex >= 0 {
			updatedAtField := reflect.Indirect(valueModelObject).Field(c.updatedAtIndex)
			if updatedAtField.Kind() == reflect.Ptr {
				updatedAtField = reflect.Indirect(updatedAtField)
			}
			updatedAtField.Set(reflect.ValueOf(time.Now()))
		}
	} else if valueModelObject.Kind() == reflect.Map {
		var updatedByTag, updatedAtTag string
		if c.updatedByIndex >= 0 {
			if updatedByTag = GetBsonName(c.modelType, c.updatedByIndex); updatedByTag == "" || updatedByTag == "-" {
				updatedByTag = GetJsonName(c.modelType, c.updatedByIndex)
			}
			if updatedByTag != "" && updatedByTag != "-" {
				valueModelObject.SetMapIndex(reflect.ValueOf(updatedByTag), reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			if updatedAtTag = GetBsonName(c.modelType, c.updatedAtIndex); updatedAtTag == "" || updatedAtTag == "-" {
				updatedAtTag = GetJsonName(c.modelType, c.updatedAtIndex)
			}
			if updatedAtTag != "" && updatedAtTag != "-" {
				valueModelObject.SetMapIndex(reflect.ValueOf(updatedAtTag), reflect.ValueOf(time.Now()))
			}
		}
	}

	return obj
}

func (c *DefaultModelBuilder) BuildToPatch(ctx context.Context, obj interface{}) interface{} {
	return c.BuildToUpdate(ctx, obj)
}

func (c *DefaultModelBuilder) BuildToSave(ctx context.Context, obj interface{}) interface{} {
	return c.BuildToUpdate(ctx, obj)
}

func GetUserIdFromContext(ctx context.Context) string {
	token := ctx.Value("authorization")
	if authorizationToken, ok := token.(map[string]interface{}); ok {
		u := authorizationToken["userId"]
		if u != nil {
			userId, _ := u.(string)
			return userId
		} else {
			u = authorizationToken["userid"]
			if u != nil {
				userId, _ := u.(string)
				return userId
			} else {
				u = authorizationToken["uid"]
				userId, _ := u.(string)
				return userId
			}
		}
		return GetUserNameFromToken(authorizationToken)
	}
	return ""
}

func GetUserNameFromToken(token map[string]interface{}) string {
	u := token["username"]
	if u != nil {
		userName, _ := u.(string)
		return userName
	} else {
		u = token["userName"]
		userName, _ := u.(string)
		return userName
	}
	return ""
}

func GetBsonName(modelType reflect.Type, index int) string {
	field := modelType.Field(index)
	if tag, ok := field.Tag.Lookup("bson"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func GetJsonName(modelType reflect.Type, index int) string {
	field := modelType.Field(index)
	if tag, ok := field.Tag.Lookup("json"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func FindFieldIndex(modelType reflect.Type, fieldName string) int {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			return i
		}
	}
	return -1
}
