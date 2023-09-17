package query

import (
	f "github.com/core-go/core/firestore"
	"github.com/core-go/core/search"
	"log"
	"reflect"
	"strings"
)

type Builder struct {
	ModelType reflect.Type
}
func UseQuery(resultModelType reflect.Type) func(interface{}) ([]f.Query, []string) {
	b := NewBuilder(resultModelType)
	return b.BuildQuery
}
func NewBuilder(resultModelType reflect.Type) *Builder {
	return &Builder{ModelType: resultModelType}
}
func (b *Builder) BuildQuery(sm interface{}) ([]f.Query, []string) {
	return BuildQueryByType(sm, b.ModelType)
}

func BuildQueryByType(sm interface{}, resultModelType reflect.Type) ([]f.Query, []string) {
	var query = make([]f.Query, 0)
	fields := make([]string, 0)

	if _, ok := sm.(*search.Filter); ok {
		return query, fields
	}

	value := reflect.Indirect(reflect.ValueOf(sm))
	numField := value.NumField()
	var keyword string
	keywordFormat := map[string]string{
		"prefix":  "==",
		"contain": "==",
		"equal":   "==",
	}
	for i := 0; i < numField; i++ {
		field := value.Field(i)
		kind := field.Kind()
		x := field.Interface()
		ps := false
		var psv string
		if kind == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			s0, ok0 := x.(*string)
			if ok0 {
				if s0 == nil || len(*s0) == 0 {
					continue
				}
				ps = true
				psv = *s0
			}
			field = field.Elem()
			kind = field.Kind()
		}
		s0, ok0 := x.(string)
		if ok0 {
			if len(s0) == 0 {
				continue
			}
			psv = s0
		}
		ks := kind.String()
		if v, ok := x.(*search.Filter); ok {
			if len(v.Fields) > 0 {
				for _, key := range v.Fields {
					i, _, columnName := getFieldByJson(resultModelType, key)
					if len(columnName) <= 0 {
						fields = fields[len(fields):]
						break
					} else if i == -1 {
						columnName = key
					}
					fields = append(fields, columnName)
				}
			}
			if len(v.Q) > 0 {
				keyword = strings.TrimSpace(v.Q)
			}
			continue
		} else if ps || ks == "string" {
			var keywordQuery f.Query
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			// var operator string
			operator := "=="
			var searchValue interface{}
			if len(psv) > 0 {
				const defaultKey = "contain"
				if key, ok := value.Type().Field(i).Tag.Lookup("operator"); ok && len(key) > 0 {
					operator = key
				} else {
					if key, ok := value.Type().Field(i).Tag.Lookup("match"); ok {
						if format, exist := keywordFormat[key]; exist {
							operator = format
						} else {
							log.Panicf("match not support \"%v\" format\n", key)
						}
					} else if format, exist := keywordFormat[defaultKey]; exist {
						operator = format
					}
				}
				searchValue = psv
			} else if len(keyword) > 0 {
				if key, ok := value.Type().Field(i).Tag.Lookup("keyword"); ok {
					if format, exist := keywordFormat[key]; exist {
						operator = format
					} else {
						log.Panicf("keyword not support \"%v\" format\n", key)
					}
				}
				searchValue = keyword
			}
			if len(columnName) > 0 && len(operator) > 0 {
				keywordQuery = f.Query{Path: columnName, Operator: operator, Value: searchValue}
				query = append(query, keywordQuery)
			}
		} else if rangeTime, ok := x.(*search.TimeRange); ok && rangeTime != nil {
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			actionTimeQuery := make([]f.Query, 0)
			if rangeTime.Min == nil {
				actionTimeQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeTime.Max}}
			} else if rangeTime.Max == nil {
				actionTimeQuery = []f.Query{{Path: columnName, Operator: ">=", Value: rangeTime.Min}}
			} else {
				actionTimeQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeTime.Max}, {Path: columnName, Operator: ">=", Value: rangeTime.Min}}
			}
			query = append(query, actionTimeQuery...)
		} else if rangeTime, ok := x.(search.TimeRange); ok {
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			actionTimeQuery := make([]f.Query, 0)
			if rangeTime.Min == nil {
				actionTimeQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeTime.Max}}
			} else if rangeTime.Max == nil {
				actionTimeQuery = []f.Query{{Path: columnName, Operator: ">=", Value: rangeTime.Min}}
			} else {
				actionTimeQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeTime.Max}, {Path: columnName, Operator: ">=", Value: rangeTime.Min}}
			}
			query = append(query, actionTimeQuery...)
		} else if rangeDate, ok := x.(*search.DateRange); ok && rangeDate != nil {
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			actionDateQuery := make([]f.Query, 0)
			if rangeDate.Min == nil && rangeDate.Max == nil {
				continue
			} else if rangeDate.Min == nil {
				actionDateQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeDate.Max}}
			} else if rangeDate.Max == nil {
				actionDateQuery = []f.Query{{Path: columnName, Operator: ">=", Value: rangeDate.Min}}
			} else {
				actionDateQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeDate.Max}, {Path: columnName, Operator: ">=", Value: rangeDate.Min}}
			}
			query = append(query, actionDateQuery...)
		} else if rangeDate, ok := x.(search.DateRange); ok {
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			actionDateQuery := make([]f.Query, 0)
			if rangeDate.Min == nil && rangeDate.Max == nil {
				continue
			} else if rangeDate.Min == nil {
				actionDateQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeDate.Max}}
			} else if rangeDate.Max == nil {
				actionDateQuery = []f.Query{{Path: columnName, Operator: ">=", Value: rangeDate.Min}}
			} else {
				actionDateQuery = []f.Query{{Path: columnName, Operator: "<=", Value: rangeDate.Max}, {Path: columnName, Operator: ">=", Value: rangeDate.Min}}
			}
			query = append(query, actionDateQuery...)
		} else if numberRange, ok := x.(*search.NumberRange); ok && numberRange != nil {
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			amountQuery := make([]f.Query, 0)

			if numberRange.Min != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: ">=", Value: *numberRange.Min})
			} else if numberRange.Lower != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: ">", Value: *numberRange.Lower})
			}
			if numberRange.Max != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: "<=", Value: *numberRange.Max})
			} else if numberRange.Upper != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: "<", Value: *numberRange.Upper})
			}

			if len(amountQuery) > 0 {
				query = append(query, amountQuery...)
			}
		} else if numberRange, ok := x.(search.NumberRange); ok {
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			amountQuery := make([]f.Query, 0)

			if numberRange.Min != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: ">=", Value: *numberRange.Min})
			} else if numberRange.Lower != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: ">", Value: *numberRange.Lower})
			}
			if numberRange.Max != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: "<=", Value: *numberRange.Max})
			} else if numberRange.Upper != nil {
				amountQuery = append(amountQuery, f.Query{Path: columnName, Operator: "<", Value: *numberRange.Upper})
			}

			if len(amountQuery) > 0 {
				query = append(query, amountQuery...)
			}
		} else if ks == "slice" && reflect.Indirect(reflect.ValueOf(x)).Len() > 0 {
			columnName := getFirestoreName(resultModelType, value.Type().Field(i).Name)
			q := f.Query{Path: columnName, Operator: "in", Value: x}
			query = append(query, q)
		} else {
			if _, ok := x.(*search.Filter); ks == "bool" || (strings.Contains(ks, "int") && x != 0) || (strings.Contains(ks, "float") && x != 0) || (!ok && ks == "ptr" && field.Pointer() != 0) {
				v := value.Type().Field(i).Name
				columnName := getFirestoreName(resultModelType, v)
				if len(columnName) > 0 {
					oper := "=="
					if key, ok := value.Type().Field(i).Tag.Lookup("operator"); ok && len(key) > 0 {
						oper = key
					}
					q := f.Query{Path: columnName, Operator: oper, Value: x}
					query = append(query, q)
				}
			}
		}
	}
	return query, fields
}

func getFieldByJson(modelType reflect.Type, jsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 && strings.Split(tag1, ",")[0] == jsonName {
			if tag2, ok2 := field.Tag.Lookup("firestore"); ok2 {
				return i, field.Name, strings.Split(tag2, ",")[0]
			}
			return i, field.Name, ""
		}
	}
	return -1, jsonName, jsonName
}
func getFirestoreName(modelType reflect.Type, fieldName string) string {
	field, _ := modelType.FieldByName(fieldName)
	bsonTag := field.Tag.Get("firestore")
	tags := strings.Split(bsonTag, ",")
	if len(tags) > 0 {
		return tags[0]
	}
	return fieldName
}
