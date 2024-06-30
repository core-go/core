package builder

import (
	"context"
	"reflect"
	"strings"
	"time"
)

type TrackingConfig struct {
	Authorization string `yaml:"authorization" mapstructure:"authorization" json:"authorization,omitempty" gorm:"column:authorization" bson:"authorization,omitempty" dynamodbav:"authorization,omitempty" firestore:"authorization,omitempty"`
	User          string `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	CreatedBy     string `yaml:"created_by" mapstructure:"created_by" json:"createdBy,omitempty" gorm:"column:createdby" bson:"createdBy,omitempty" dynamodbav:"createdBy,omitempty" firestore:"createdBy,omitempty"`
	CreatedAt     string `yaml:"created_at" mapstructure:"created_at" json:"createdAt,omitempty" gorm:"column:createdat" bson:"createdAt,omitempty" dynamodbav:"createdAt,omitempty" firestore:"createdAt,omitempty"`
	UpdatedBy     string `yaml:"updated_by" mapstructure:"updated_by" json:"updatedBy,omitempty" gorm:"column:updatedby" bson:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty" firestore:"updatedBy,omitempty"`
	UpdatedAt     string `yaml:"updated_at" mapstructure:"updated_at" json:"updatedAt,omitempty" gorm:"column:updatedat" bson:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty" firestore:"updatedAt,omitempty"`
}
type Builder[T any] struct {
	GenerateId     func(ctx context.Context, model *T) (int, error)
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

func NewBuilderWithIdAndConfig[T any](generateId func(context.Context) (string, error), c TrackingConfig) *Builder[T] {
	if generateId != nil {
		idGenerator := NewIdGenerator[T](generateId)
		return NewBuilderByConfig[T](idGenerator.Generate, c)
	} else {
		return NewBuilderByConfig[T](nil, c)
	}
}
func NewBuilderByConfig[T any](generateId func(context.Context, *T) (int, error), c TrackingConfig) *Builder[T] {
	return NewBuilder[T](generateId, c.CreatedBy, c.CreatedAt, c.UpdatedBy, c.UpdatedAt, c.User, c.Authorization)
}
func NewBuilderWithId[T any](generateId func(context.Context) (string, error), options ...string) *Builder[T] {
	if generateId != nil {
		idGenerator := NewIdGenerator[T](generateId)
		return NewBuilder[T](idGenerator.Generate, options...)
	} else {
		return NewBuilder[T](nil, options...)
	}
}
func NewBuilder[T any](generateId func(context.Context, *T) (int, error), opts ...string) *Builder[T] {
	var t T
	modelType := reflect.TypeOf(t)
	if modelType.Kind() != reflect.Struct {
		panic("T must be a struct")
	}
	var createdByName, createdAtName, updatedByName, updatedAtName, key, authorization string
	if len(opts) > 0 {
		createdByName = opts[0]
	}
	if len(opts) > 1 {
		createdAtName = opts[1]
	}
	if len(opts) > 2 {
		updatedByName = opts[2]
	}
	if len(opts) > 3 {
		updatedAtName = opts[3]
	}
	if len(opts) > 4 && len(opts[4]) > 0 {
		key = opts[4]
	} else {
		key = "userId"
	}
	if len(opts) > 5 {
		authorization = opts[5]
	}
	createdByIndex := findFieldIndex(modelType, createdByName)
	createdAtIndex := findFieldIndex(modelType, createdAtName)
	updatedByIndex := findFieldIndex(modelType, updatedByName)
	updatedAtIndex := findFieldIndex(modelType, updatedAtName)

	return &Builder[T]{
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

func (c *Builder[T]) Create(ctx context.Context, obj *T) error {
	if c.GenerateId != nil {
		_, er0 := c.GenerateId(ctx, obj)
		if er0 != nil {
			return er0
		}
	}
	v := reflect.Indirect(reflect.ValueOf(obj))
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	userId := fromContext(ctx, c.Key, c.Authorization)
	if v.Kind() == reflect.Struct {
		if c.createdByIndex >= 0 {
			createdByField := reflect.Indirect(v).Field(c.createdByIndex)
			if createdByField.Kind() == reflect.Ptr {
				createdByField.Set(reflect.ValueOf(&userId))
			} else {
				createdByField.Set(reflect.ValueOf(userId))
			}
		}

		if c.createdAtIndex >= 0 {
			createdAtField := reflect.Indirect(v).Field(c.createdAtIndex)
			t := time.Now()
			if createdAtField.Kind() == reflect.Ptr {
				createdAtField.Set(reflect.ValueOf(&t))
			} else {
				createdAtField.Set(reflect.ValueOf(t))
			}
		}

		if c.updatedByIndex >= 0 {
			updatedByField := reflect.Indirect(v).Field(c.updatedByIndex)
			if updatedByField.Kind() == reflect.Ptr {
				updatedByField.Set(reflect.ValueOf(&userId))
			} else {
				updatedByField.Set(reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			updatedAtField := v.Field(c.updatedAtIndex)
			t := time.Now()
			if updatedAtField.Kind() == reflect.Ptr {
				updatedAtField.Set(reflect.ValueOf(&t))
				//updatedAtField = reflect.Indirect(updatedAtField)
			} else {
				updatedAtField.Set(reflect.ValueOf(t))
			}
		}
	} else if v.Kind() == reflect.Map {
		var createdByTag, createdAtTag string
		if c.createdByIndex >= 0 {
			if createdByTag = getJsonName(c.modelType, c.createdByIndex); createdByTag == "" || createdByTag == "-" {
				createdByTag = getBsonName(c.modelType, c.createdByIndex)
			}
			if createdByTag != "" && createdByTag != "-" {
				v.SetMapIndex(reflect.ValueOf(createdByTag), reflect.ValueOf(userId))
			}
		}
		if c.createdAtIndex >= 0 {
			if createdAtTag = getJsonName(c.modelType, c.createdAtIndex); createdAtTag == "" || createdAtTag == "-" {
				createdAtTag = getBsonName(c.modelType, c.createdAtIndex)
			}
			if createdAtTag != "" && createdAtTag != "-" {
				v.SetMapIndex(reflect.ValueOf(createdAtTag), reflect.ValueOf(time.Now()))
			}
		}
		var updatedByTag, updatedAtTag string
		if c.updatedByIndex >= 0 {
			if updatedByTag = getJsonName(c.modelType, c.updatedByIndex); updatedByTag == "" || updatedByTag == "-" {
				updatedByTag = getBsonName(c.modelType, c.updatedByIndex)
			}
			if updatedByTag != "" && updatedByTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedByTag), reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			if updatedAtTag = getJsonName(c.modelType, c.updatedAtIndex); updatedAtTag == "" || updatedAtTag == "-" {
				updatedAtTag = getBsonName(c.modelType, c.updatedAtIndex)
			}
			if updatedAtTag != "" && updatedAtTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedAtTag), reflect.ValueOf(time.Now()))
			}
		}
	}
	return nil
}

func (c *Builder[T]) Update(ctx context.Context, obj *T) error {
	v := reflect.Indirect(reflect.ValueOf(obj))
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	userId := fromContext(ctx, c.Key, c.Authorization)
	if v.Kind() == reflect.Struct {
		if c.updatedByIndex >= 0 {
			updatedByField := reflect.Indirect(v).Field(c.updatedByIndex)
			if updatedByField.Kind() == reflect.Ptr {
				updatedByField.Set(reflect.ValueOf(&userId))
			} else {
				updatedByField.Set(reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			updatedAtField := v.Field(c.updatedAtIndex)
			t := time.Now()
			if updatedAtField.Kind() == reflect.Ptr {
				updatedAtField.Set(reflect.ValueOf(&t))
				//updatedAtField = reflect.Indirect(updatedAtField)
			} else {
				updatedAtField.Set(reflect.ValueOf(t))
			}
		}
	} else if v.Kind() == reflect.Map {
		var updatedByTag, updatedAtTag string
		if c.updatedByIndex >= 0 {
			if updatedByTag = getJsonName(c.modelType, c.updatedByIndex); updatedByTag == "" || updatedByTag == "-" {
				updatedByTag = getBsonName(c.modelType, c.updatedByIndex)
			}
			if updatedByTag != "" && updatedByTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedByTag), reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			if updatedAtTag = getJsonName(c.modelType, c.updatedAtIndex); updatedAtTag == "" || updatedAtTag == "-" {
				updatedAtTag = getBsonName(c.modelType, c.updatedAtIndex)
			}
			if updatedAtTag != "" && updatedAtTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedAtTag), reflect.ValueOf(time.Now()))
			}
		}
	}

	return nil
}

func fromContext(ctx context.Context, key string, authorization string) string {
	if len(authorization) > 0 {
		token := ctx.Value(authorization)
		if token != nil {
			if authorizationToken, exist := token.(map[string]interface{}); exist {
				return fromMap(key, authorizationToken)
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
func fromMap(key string, data map[string]interface{}) string {
	u := data[key]
	if u != nil {
		v, ok := u.(string)
		if ok {
			return v
		}
	}
	return ""
}
func getBsonName(modelType reflect.Type, index int) string {
	field := modelType.Field(index)
	if tag, ok := field.Tag.Lookup("bson"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func getJsonName(modelType reflect.Type, index int) string {
	field := modelType.Field(index)
	if tag, ok := field.Tag.Lookup("json"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func findFieldIndex(modelType reflect.Type, fieldName string) int {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			return i
		}
	}
	return -1
}
