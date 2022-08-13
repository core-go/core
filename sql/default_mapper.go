package sql

import (
	"context"
	"encoding/json"
	"reflect"
)

type DefaultMapper struct{}

func NewMapper() *DefaultMapper {
	return &DefaultMapper{}
}

func (d *DefaultMapper) DbToModel(ctx context.Context, model interface{}) (interface{}, error) {
	var value reflect.Value
	if val, ok := model.(reflect.Value); ok {
		value = reflect.Indirect(val)
	} else {
		m := reflect.Indirect(reflect.ValueOf(model)).Interface()
		value = reflect.Indirect(reflect.ValueOf(m))
	}
	if value.Kind() == reflect.Struct {
		DbToModel(value)
	}
	return model, nil
}

func (d *DefaultMapper) DbToModels(ctx context.Context, models interface{}) (interface{}, error) {
	vo := reflect.Indirect(reflect.ValueOf(models))
	if vo.Kind() == reflect.Slice {
		if vo.Len() == 0 {
			return 0, nil
		}
		for i := 0; i < vo.Len(); i++ {
			model := vo.Index(i).Addr()
			d.DbToModel(ctx, model)
		}
	}
	return models, nil
}

func (d *DefaultMapper) ModelToDb(ctx context.Context, model interface{}) (interface{}, error) {
	var value reflect.Value
	if val, ok := model.(reflect.Value); ok {
		value = reflect.Indirect(val)
	} else {
		value = reflect.Indirect(reflect.ValueOf(model))
		kind := value.Kind()
		if kind == reflect.Ptr || kind == reflect.Interface {
			value = reflect.Indirect(reflect.ValueOf(value.Interface()))
		}
	}
	if value.Kind() == reflect.Struct {
		ModelToDb(value)
	}
	return model, nil
}

func (d *DefaultMapper) ModelsToDb(ctx context.Context, models interface{}) (interface{}, error) {
	vo := reflect.Indirect(reflect.ValueOf(models))
	if vo.Kind() == reflect.Slice {
		if vo.Len() == 0 {
			return 0, nil
		}
		for i := 0; i < vo.Len(); i++ {
			model := vo.Index(i).Addr()
			d.ModelToDb(ctx, model)
		}
	}
	return models, nil
}

func DbToModel(value reflect.Value) {
	mapTag := GetIndexTagFields("field", value.Type())
	for i, fieldName := range mapTag {
		j := GetIndexField(fieldName, "col", value.Type())
		if value.CanAddr() {
			data := value.Field(j).Addr().Interface()
			json.Unmarshal([]byte(value.Field(i).String()), data)
		} else {
			data := value.Field(j)
			json.Unmarshal([]byte(value.Field(i).String()), &data)
		}
	}
}

func ModelToDb(value reflect.Value) {
	mapTag := GetIndexTagFields("col", value.Type())
	for i, fieldName := range mapTag {
		j := GetIndexField(fieldName, "field", value.Type())
		data := value.Field(j)
		if parsedToClob, err := json.Marshal(value.Field(i).Interface()); err == nil && data.CanSet() {
			data.SetString(string(parsedToClob))
		}
	}
}

func GetIndexTagFields(key string, modelType reflect.Type) map[int]string {
	mapTag := make(map[int]string)
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		tagField := f.Tag.Get(key)
		if tagField != "" {
			mapTag[i] = tagField
		}
	}
	return mapTag
}

func GetIndexField(field, key string, modelType reflect.Type) (index int) {
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		v := f.Tag.Get(key)
		if v == field {
			return i
		}
	}
	return -1
}
