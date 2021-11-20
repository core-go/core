package convert

import "reflect"

func Clone(origin interface{}) interface{} {
	originValue := reflect.Indirect(reflect.ValueOf(origin))
	originType := reflect.TypeOf(origin)

	resultType := reflect.TypeOf(origin)
	result := reflect.New(resultType)
	numFields := originType.NumField()
	for i := 0; i < numFields; i++ {
		field := originType.Field(i)
		value := originValue.FieldByName(field.Name)
		f := result.Elem().Field(i)
		if value.Kind() == reflect.String {
			f.SetString(value.String())
		} else if value.Kind() == reflect.Int {
			f.SetInt(value.Int())
		} else if value.Kind() == reflect.Float64 {
			f.SetFloat(value.Float())
		} else if value.Kind() == reflect.Bool {
			f.SetBool(value.Bool())
		} else if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			} else {
				val := value.Interface()
				switch val.(type) {
				case *string:
					strVal, ok := val.(*string)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&strVal)))
					}
				case *int:
					intVal, ok := val.(*int)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&intVal)))
					}
				case *float64:
					floatVal, ok := val.(*float64)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&floatVal)))
					}
				case *bool:
					boolVal, ok := val.(*bool)
					if ok {
						f.Set(reflect.Indirect(reflect.ValueOf(&boolVal)))
					}
				}
			}
		} else if value.Kind() == reflect.Struct {
			data := Clone(value.Interface())
			f.Set(reflect.Indirect(reflect.ValueOf(data)))
		}
	}
	return result.Interface()
}
