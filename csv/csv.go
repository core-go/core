package csv

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)

const DateLayout string = "2006-01-02 15:04:05 +0700 +07"

func BuildCsv(rows []string, fields []string, valueOfmodels reflect.Value, embedFieldName string) []string {
	if lengthResult := valueOfmodels.Len(); lengthResult > 0 {
		model := valueOfmodels.Index(0).Interface()

		firstLayerIndexes, secondLayerIndexes := findIndexByTagJsonOrEmbededFieldName(model, fields, embedFieldName)

		for i := 0; i < lengthResult; i++ {
			var cols []string
			valueOfmodel := valueOfmodels.Index(i)
			for _, fieldName := range fields {
				if index, exist := firstLayerIndexes[fieldName]; exist {
					valueOfFieldName := valueOfmodel.Field(index)
					cols = AppendColumns(valueOfFieldName, cols)
				} else if index, exist := secondLayerIndexes[fieldName]; exist {
					embedFieldValue := reflect.Indirect(valueOfmodel.Field(firstLayerIndexes[embedFieldName]))
					valueOfFieldName := embedFieldValue.Field(index)
					cols = AppendColumns(valueOfFieldName, cols)
				}
			}
			rows = append(rows, strings.Join(cols, ","))
		}
	}
	return rows
}

func AppendColumns(value reflect.Value, cols []string) []string {
	kind := value.Kind()
	if kind == reflect.Ptr && value.IsNil() {
		cols = append(cols, "")
	} else {
		i := value.Interface()
		if kind == reflect.Ptr {
			i = reflect.Indirect(reflect.ValueOf(i)).Interface()
		}
		v, okS := i.(string)
		if okS {
			c := strings.Contains(v, `"`)
			if c || strings.Contains(v, ",") {
				if c {
					v = strings.ReplaceAll(v, `"`, `""`)
				}
				v = "\"" + v + "\""
				cols = append(cols, v)
			} else {
				cols = append(cols, v)
			}
		} else {
			d, okD := i.(time.Time)
			if okD {
				v := d.Format(DateLayout)
				cols = append(cols, v)
			} else {
				f, okBF := i.(big.Float)
				if okBF {
					cols = append(cols, f.String())
				} else {
					bi, okBI := i.(big.Int)
					if okBI {
						cols = append(cols, fmt.Sprintf("%v", bi.String()))
					} else {
						cols = append(cols, fmt.Sprintf("%v", i))
					}
				}
			}
		}
	}
	return cols
}

func findIndexByTagJsonOrEmbededFieldName(model interface{}, jsonNames []string, embedFieldName string) (firstLayerIndex map[string]int, secondLayerIndexes map[string]int) {
	tmp := make([]string, len(jsonNames))
	copy(tmp, jsonNames)

	firstLayerIndex = map[string]int{}
	secondLayerIndexes = map[string]int{}
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	numField := modelValue.NumField()

	for i := 0; i < numField; i++ {
		if jsonTag, exist := modelValue.Type().Field(i).Tag.Lookup("json"); exist {
			for j, name := range tmp {
				tags := strings.Split(jsonTag, ",")
				for _, tag := range tags {
					if strings.Compare(strings.TrimSpace(tag), name) == 0 {
						firstLayerIndex[name] = i
						tmp = append(tmp[:j], tmp[j+1:]...)
						break
					}
				}
			}
		}
		if modelValue.Type().Field(i).Name == embedFieldName {
			firstLayerIndex[embedFieldName] = i
			for j, name := range tmp {
				embedValue := reflect.Indirect(modelValue.Field(i))
				if index, _ := findIndexByTagJson(embedValue.Type(), name); index != -1 {
					secondLayerIndexes[name] = index
					tmp = append(tmp[:j], tmp[j+1:]...)
					break
				}
			}
		}
	}
	return
}

func findIndexByTagJson(modelType reflect.Type, jsonName string) (int, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		jsonTag := field.Tag.Get("json")
		tags := strings.Split(jsonTag, ",")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), jsonName) == 0 {
				return i, field.Name
			}
		}
	}
	return -1, ""
}
