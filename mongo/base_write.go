package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"reflect"
)

//For Insert
func mapObjectIdToModel(id primitive.ObjectID, valueOfModel reflect.Value, idIndex int) {
	switch reflect.Indirect(valueOfModel).Field(idIndex).Kind() {
	case reflect.String:
		if _, err := setValue(valueOfModel, idIndex, id.Hex()); err != nil {
			log.Println("Err: " + err.Error())
		}
		break
	default:
		if _, err := setValue(valueOfModel, idIndex, id); err != nil {
			log.Println("Err: " + err.Error())
		}
		break
	}
}
func setValue(model interface{}, index int, value interface{}) (interface{}, error) {
	vo := reflect.Indirect(reflect.ValueOf(model))
	switch reflect.ValueOf(model).Kind() {
	case reflect.Ptr:
		{
			vo.Field(index).Set(reflect.ValueOf(value))
			return model, nil
		}
	default:
		if modelWithTypeValue, ok := model.(reflect.Value); ok {
			_, err := setValueWithTypeValue(modelWithTypeValue, index, value)
			return modelWithTypeValue.Interface(), err
		}
	}
	return model, nil
}
func setValueWithTypeValue(model reflect.Value, index int, value interface{}) (reflect.Value, error) {
	trueValue := reflect.Indirect(model)
	switch trueValue.Kind() {
	case reflect.Struct:
		{
			val := reflect.Indirect(reflect.ValueOf(value))
			if trueValue.Field(index).Kind() == val.Kind() {
				trueValue.Field(index).Set(reflect.ValueOf(value))
				return trueValue, nil
			} else {
				return trueValue, fmt.Errorf("value's kind must same as field's kind")
			}
		}
	default:
		return trueValue, nil
	}
}
// Update
func getValue(model interface{}, index int) (interface{}, error) {
	vo := reflect.Indirect(reflect.ValueOf(model))
	return vo.Field(index).Interface(), nil
}
func findIndex(model interface{}, fieldName string) int {
	modelType := reflect.Indirect(reflect.ValueOf(model))
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		if modelType.Type().Field(i).Name == fieldName {
			return i
		}
	}
	return -1
}
