package service

import (
	"context"
	"reflect"
	"strings"
	"time"
)

type TrackingConfig struct {
	Authorization string `mapstructure:"authorization" json:"authorization,omitempty" gorm:"column:authorization" bson:"authorization,omitempty" dynamodbav:"authorization,omitempty" firestore:"authorization,omitempty"`
	User          string `mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	CreatedBy     string `mapstructure:"created_by" json:"createdBy,omitempty" gorm:"column:createdby" bson:"createdBy,omitempty" dynamodbav:"createdBy,omitempty" firestore:"createdBy,omitempty"`
	CreationTime  string `mapstructure:"creation_time" json:"creationTime,omitempty" gorm:"column:creationtime" bson:"creationTime,omitempty" dynamodbav:"creationTime,omitempty" firestore:"creationTime,omitempty"`
	UpdatedBy     string `mapstructure:"updated_by" json:"updatedBy,omitempty" gorm:"column:updatedby" bson:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty" firestore:"updatedBy,omitempty"`
	UpdateTime    string `mapstructure:"update_time" json:"updateTime,omitempty" gorm:"column:updatetime" bson:"updateTime,omitempty" dynamodbav:"updateTime,omitempty" firestore:"updateTime,omitempty"`
}
type DefaultModelBuilder struct {
	GenerateId     func(ctx context.Context, model interface{}) (int, error)
	Authorization  string
	Key            string
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

func NewModelBuilderByConfig(generateId func(ctx context.Context, model interface{}) (int, error), modelType reflect.Type, c TrackingConfig) *DefaultModelBuilder {
	return NewModelBuilder(generateId, modelType, c.Authorization, c.User, c.CreatedBy, c.CreationTime, c.UpdatedBy, c.UpdateTime)
}
func NewModelBuilder(generateId func(ctx context.Context, model interface{}) (int, error), modelType reflect.Type, authorization string, key string, createdByName, createdAtName, updatedByName, updatedAtName string) *DefaultModelBuilder {
	createdByIndex := FindFieldIndex(modelType, createdByName)
	createdAtIndex := FindFieldIndex(modelType, createdAtName)
	updatedByIndex := FindFieldIndex(modelType, updatedByName)
	updatedAtIndex := FindFieldIndex(modelType, updatedAtName)

	return &DefaultModelBuilder{
		GenerateId:     generateId,
		Authorization:  authorization,
		Key:            key,
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
	if c.GenerateId != nil {
		c.GenerateId(ctx, obj)
	}
	valueModelObject := reflect.Indirect(reflect.ValueOf(obj))
	if valueModelObject.Kind() == reflect.Ptr {
		valueModelObject = reflect.Indirect(valueModelObject)
	}
	userId := FromContext(ctx, c.Authorization, c.Key)
	if valueModelObject.Kind() == reflect.Struct {
		if c.createdByIndex >= 0 {
			createdByField := reflect.Indirect(valueModelObject).Field(c.createdByIndex)
			if createdByField.Kind() == reflect.Ptr {
				createdByField.Set(reflect.ValueOf(&userId))
			} else {
				createdByField.Set(reflect.ValueOf(userId))
			}
		}
		if c.createdAtIndex >= 0 {
			createdAtField := reflect.Indirect(valueModelObject).Field(c.createdAtIndex)
			t := time.Now()
			if createdAtField.Kind() == reflect.Ptr {
				createdAtField.Set(reflect.ValueOf(&t))
			} else {
				createdAtField.Set(reflect.ValueOf(t))
			}
		}

		if c.updatedByIndex >= 0 {
			updatedByField := reflect.Indirect(valueModelObject).Field(c.updatedByIndex)
			if updatedByField.Kind() == reflect.Ptr {
				updatedByField.Set(reflect.ValueOf(&userId))
			} else {
				updatedByField.Set(reflect.ValueOf(userId))
			}
		}
		if c.updatedAtIndex >= 0 {
			updatedAtField := reflect.Indirect(valueModelObject).Field(c.updatedAtIndex)
			t := time.Now()
			if updatedAtField.Kind() == reflect.Ptr {
				updatedAtField.Set(reflect.ValueOf(&t))
			} else {
				updatedAtField.Set(reflect.ValueOf(t))
			}
		}
	} else if valueModelObject.Kind() == reflect.Map {
		var createdByTag, createdAtTag, updatedByTag, updatedAtTag string
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

func (c *DefaultModelBuilder) BuildToUpdate(ctx context.Context, obj interface{}) interface{} {
	valueModelObject := reflect.Indirect(reflect.ValueOf(obj))
	if valueModelObject.Kind() == reflect.Ptr {
		valueModelObject = reflect.Indirect(valueModelObject)
	}
	userId := FromContext(ctx, c.Authorization, c.Key)
	if valueModelObject.Kind() == reflect.Struct {
		if c.updatedByIndex >= 0 {
			updatedByField := reflect.Indirect(valueModelObject).Field(c.updatedByIndex)
			if updatedByField.Kind() == reflect.Ptr {
				updatedByField.Set(reflect.ValueOf(&userId))
			} else {
				updatedByField.Set(reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			updatedAtField := valueModelObject.Field(c.updatedAtIndex)
			t := time.Now()
			if updatedAtField.Kind() == reflect.Ptr {
				updatedAtField.Set(reflect.ValueOf(&t))
				//updatedAtField = reflect.Indirect(updatedAtField)
			} else {
				updatedAtField.Set(reflect.ValueOf(t))
			}
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
func FromContext(ctx context.Context, authorization string, key string) string {
	if len(authorization) > 0 {
		token := ctx.Value(authorization)
		if token != nil {
			if authorizationToken, exist := token.(map[string]interface{}); exist {
				return FromMap(key, authorizationToken)
			}
		}
		return ""
	} else {
		u := ctx.Value(key)
		if u != nil {
			v, ok := u.(string)
			if ok {
				return v
			}
		}
		return ""
	}
}
func FromMap(key string, data map[string]interface{}) string {
	u := data[key]
	if u != nil {
		v, ok := u.(string)
		if ok {
			return v
		}
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
