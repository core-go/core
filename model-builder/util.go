package builder

import "reflect"

func setValue(model interface{}, index int, value interface{}) (interface{}, error) {
	valueModelObject := reflect.Indirect(reflect.ValueOf(model))
	if valueModelObject.Kind() == reflect.Ptr {
		valueModelObject = reflect.Indirect(valueModelObject)
	}

	valueModelObject.Field(index).Set(reflect.ValueOf(value))
	return model, nil
}
