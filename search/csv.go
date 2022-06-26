package search

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)

const DateLayout string = "2006-01-02 15:04:05 +0700 +07"

func BuildCsv(rows []string, fields []string, valueOfmodels reflect.Value, embedFieldName string, opts...map[string]int) []string {
	if lengthResult := valueOfmodels.Len(); lengthResult > 0 {
		model := valueOfmodels.Index(0).Interface()

		var firstLayerIndexes map[string]int
		var secondLayerIndexes map[string]int
		if len(opts) > 0 && opts[0] != nil {
			firstLayerIndexes = opts[0]
			if len(opts) > 1 && opts[1] != nil {
				secondLayerIndexes = opts[1]
			}
		} else {
			firstLayerIndexes, secondLayerIndexes = BuildJsonMap(model, fields, embedFieldName)
		}

		for i := 0; i < lengthResult; i++ {
			var cols []string
			valueOfmodel := valueOfmodels.Index(i)
			for _, fieldName := range fields {
				index0, exist0 := secondLayerIndexes[fieldName]
				if exist0 {
					fmt.Sprintf("%d ", index0)
				} else {
					fmt.Sprintf("not exist")
				}
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
						if kind == reflect.Struct {
							if v2 := reflect.Indirect(reflect.ValueOf(i)); v2.NumField() == 1 {
								field := v2.Field(0)
								fv := field.Interface()
								k := field.Kind()
								if k == reflect.Ptr {
									fv = reflect.Indirect(reflect.ValueOf(fv)).Interface()
								}
								if sv, ok := fv.(big.Float); ok {
									cols = append(cols, sv.String())
								} else if svi, ok := fv.(big.Int); ok {
									cols = append(cols, svi.Text(10))
								} else {
									cols = append(cols, fmt.Sprintf("%v", fv))
								}
							} else {
								cols = append(cols, fmt.Sprintf("%v", i))
							}
						} else {
							cols = append(cols, fmt.Sprintf("%v", i))
						}
					}
				}
			}
		}
	}
	return cols
}

func BuildJsonMap(model interface{}, jsonNames []string, embedFieldName string) (firstLayerIndex map[string]int, secondLayerIndexes map[string]int) {
	tmp := make([]string, len(jsonNames))
	copy(tmp, jsonNames)

	firstLayerIndex = map[string]int{}
	secondLayerIndexes = map[string]int{}
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	numField := modelValue.NumField()
	modelType := modelValue.Type()
	for i := 0; i < numField; i++ {
		if jsonTag, exist := modelType.Field(i).Tag.Lookup("json"); exist {
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
		if modelType.Field(i).Name == embedFieldName {
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
