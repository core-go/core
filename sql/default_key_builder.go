package sql

import (
	"fmt"
	"reflect"
	"strings"
)

type DefaultKeyBuilder struct {
	PositionPrimaryKeysMap map[reflect.Type][]int
}

func NewDefaultKeyBuilder() *DefaultKeyBuilder {
	return &DefaultKeyBuilder{PositionPrimaryKeysMap: make(map[reflect.Type][]int)}
}

func (b *DefaultKeyBuilder) getPositionPrimaryKeys(modelType reflect.Type) []int {
	if b.PositionPrimaryKeysMap[modelType] == nil {
		var positions []int

		numField := modelType.NumField()
		for i := 0; i < numField; i++ {
			gorm := strings.Split(modelType.Field(i).Tag.Get("gorm"), ";")
			for _, value := range gorm {
				if value == "primary_key" {
					positions = append(positions, i)
					break
				}
			}
		}

		b.PositionPrimaryKeysMap[modelType] = positions
	}

	return b.PositionPrimaryKeysMap[modelType]
}

func (b *DefaultKeyBuilder) BuildKey(object interface{}) string {
	ids := make(map[string]interface{})
	objectValue := reflect.Indirect(reflect.ValueOf(object))
	positions := b.getPositionPrimaryKeys(objectValue.Type())
	var values []string
	for _, position := range positions {
		if _, colName, ok := GetFieldByIndex(objectValue.Type(), position); ok {
			ids[colName] = fmt.Sprint(objectValue.Field(position).Interface())
			values = append(values, fmt.Sprint(objectValue.Field(position).Interface()))
		}
	}
	return strings.Join(values, "-")
}

func (b *DefaultKeyBuilder) BuildKeyFromMap(keyMap map[string]interface{}, idNames []string) string {
	var values []string
	for _, key := range idNames {
		if keyVal, exist := keyMap[key]; !exist {
			values = append(values, "")
		} else {
			str, ok := keyVal.(string)
			if !ok {
				return ""
			}
			values = append(values, str)
		}
	}
	return strings.Join(values, "-")
}
