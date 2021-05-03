package builder

import "reflect"

func setValue(model interface{}, index int, value interface{}) (interface{}, error) {
	v := reflect.Indirect(reflect.ValueOf(model))
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	v.Field(index).Set(reflect.ValueOf(value))
	return model, nil
}
