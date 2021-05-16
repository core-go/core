package csv

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	layoutDate string = "2006-01-02 15:04:05 +0700 +07"
	layout     string = "2006-01-02T15:04:05"
)

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
	const e = ""
	const s = "string"
	const in = "int"
	const f = "float64"
	const x = "\""
	const y = "\"\""
	var v = fmt.Sprintf("%v", value)
	if v == "" || v == "0" || v == "<nil>" {
		cols = append(cols, "")
	} else {
		if fmt.Sprintf("%v", value.Kind()) == s {
			if strings.Contains(v, ",") {
				//a := "\"" + string(strings.ReplaceAll(v, x, y)) + "\""
				cols = append(cols, "")
			} else {
				cols = append(cols, fmt.Sprintf("%v", v))
			}
		} else if fmt.Sprintf("%v", value.Kind()) == "ptr" || fmt.Sprintf("%v", value.Kind()) == "struct" {
			fieldDate, err := time.Parse(layoutDate, v)
			if err != nil {
				fmt.Println("err", fmt.Sprintf("%v", err))
				cols = append(cols, fmt.Sprintf("%v", fmt.Sprintf("%v", v)))
			} else {
				cols = append(cols, fmt.Sprintf("%v", fieldDate.UTC().Format(layout)))
			}
		} else if fmt.Sprintf("%v", value.Kind()) == in || fmt.Sprintf("%v", value.Kind()) == f {
			cols = append(cols, fmt.Sprintf("%v", v))
		} else {
			cols = append(cols, fmt.Sprintf("%v", ""))
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
