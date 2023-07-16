package convert

import (
	s "github.com/core-go/core/search"
	"reflect"
	"strings"
)

const (
	desc = "desc"
	asc  = "asc"
)
func ToMap(in interface{}, modelType *reflect.Type, opts...func(sortString string, modelType reflect.Type) string) map[string]interface{} {
	return ToMapWithFields(in, "", modelType, opts...)
}
func ToMapWithFields(in interface{}, sfields string, modelType *reflect.Type, opts...func(sortString string, modelType reflect.Type) string) map[string]interface{} {
	var buildSort func(string, reflect.Type) string
	if len(opts) > 0 && opts[0] != nil {
		buildSort = opts[0]
	} else {
		buildSort = BuildSort
	}
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fv := f.Interface()
		k := f.Kind()
		if k == reflect.Ptr {
			if f.IsNil() {
				continue
			} else {
				fv = reflect.Indirect(reflect.ValueOf(fv)).Interface()
			}
		} else if k == reflect.Slice {
			if f.IsNil() {
				continue
			}
		}
		sv, ok := fv.(string)
		if ok {
			if len(sv) > 0 {
				n := getTag(typ.Field(i), "json")
				m := getTag(typ.Field(i), "operator")
				w := sv
				if m == "like" {
					w = Q(sv)
				} else if m != "=" {
					w = Prefix(sv)
				}
				out[n] = w
			}
			continue
		}
		if v, ok := fv.(s.Filter); ok {
			if modelType != nil {
				if len(v.Fields) > 0 {
					fields := make([]string, 0)
					for _, key := range v.Fields {
						i, _, columnName := getFieldByJson(*modelType, key)
						if len(columnName) < 0 {
							fields = fields[len(fields):]
							break
						} else if i > -1 {
							fields = append(fields, columnName)
						}
					}
					if len(fields) > 0 {
						out["fields"] = strings.Join(fields, ",")
					} else {
						if len(sfields) > 0 {
							out["fields"] = sfields
						}
					}
				} else if len(sfields) > 0 {
					out["fields"] = sfields
				}
			}
			if len(v.Sort) > 0 {
				t := typ
				if *modelType != nil {
					t = *modelType
				}
				sortString := buildSort(v.Sort, t)
				if len(sortString) > 0 {
					out["sort"] = sortString
				}
			}
			if v.Excluding != nil && len(v.Excluding) > 0 {
				out["excluding"] = v.Excluding
			}
			if len(v.Q) > 0 {
				out["q"] = strings.TrimSpace(v.Q)
			}
		} else {
			n := getTag(typ.Field(i), "json")
			if dateTime, ok := fv.(s.TimeRange); ok {
				if dateTime.Min != nil || dateTime.Max != nil || dateTime.Top != nil {
					sub := make(map[string]interface{})
					if dateTime.Min != nil {
						sub["min"] = *dateTime.Min
					}
					if dateTime.Max != nil {
						sub["max"] = *dateTime.Max
					} else if dateTime.Top != nil {
						sub["top"] = *dateTime.Top
					}
					out[n] = sub
				}
			} else if numberRange, ok := fv.(s.NumberRange); ok {
				if numberRange.Min != nil || numberRange.Max != nil || numberRange.Top != nil {
					sub := make(map[string]interface{})
					if numberRange.Min != nil {
						sub["min"] = *numberRange.Min
					}
					if numberRange.Max != nil {
						sub["max"] = *numberRange.Max
					} else if numberRange.Top != nil {
						sub["top"] = *numberRange.Top
					}
					out[n] = sub
				}
			} else if numberRange, ok := fv.(s.Int64Range); ok {
				if numberRange.Min != nil || numberRange.Max != nil || numberRange.Top != nil {
					sub := make(map[string]interface{})
					if numberRange.Min != nil {
						sub["min"] = *numberRange.Min
					}
					if numberRange.Max != nil {
						sub["max"] = *numberRange.Max
					} else if numberRange.Top != nil {
						sub["top"] = *numberRange.Top
					}
					out[n] = sub
				}
			} else if numberRange, ok := fv.(s.IntRange); ok {
				if numberRange.Min != nil || numberRange.Max != nil || numberRange.Top != nil {
					sub := make(map[string]interface{})
					if numberRange.Min != nil {
						sub["min"] = *numberRange.Min
					}
					if numberRange.Max != nil {
						sub["max"] = *numberRange.Max
					} else if numberRange.Top != nil {
						sub["top"] = *numberRange.Top
					}
					out[n] = sub
				}
			} else if numberRange, ok := fv.(s.Int32Range); ok {
				if numberRange.Min != nil || numberRange.Max != nil || numberRange.Top != nil {
					sub := make(map[string]interface{})
					if numberRange.Min != nil {
						sub["min"] = *numberRange.Min
					}
					if numberRange.Max != nil {
						sub["max"] = *numberRange.Max
					} else if numberRange.Top != nil {
						sub["top"] = *numberRange.Top
					}
					out[n] = sub
				}
			} else {
				out[n] = fv
			}
		}
	}
	return out
}
func getTag(fi reflect.StructField, tag string) string {
	if tagv := fi.Tag.Get(tag); tagv != "" {
		arrValue := strings.Split(tagv, ",")
		if len(arrValue) > 0 {
			return arrValue[0]
		} else {
			return tagv
		}
	}
	return fi.Name
}
func BuildSort(sortString string, modelType reflect.Type) string {
	var sort = make([]string, 0)
	sorts := strings.Split(sortString, ",")
	for i := 0; i < len(sorts); i++ {
		sortField := strings.TrimSpace(sorts[i])
		fieldName := sortField
		c := sortField[0:1]
		if c == "-" || c == "+" {
			fieldName = sortField[1:]
		}
		columnName := getColumnNameForSearch(modelType, fieldName)
		if len(columnName) > 0 {
			sortType := getSortType(c)
			sort = append(sort, columnName+" "+sortType)
		}
	}
	if len(sort) > 0 {
		return strings.Join(sort, ",")
	} else {
		return ""
	}
}
func getColumnNameForSearch(modelType reflect.Type, sortField string) string {
	sortField = strings.TrimSpace(sortField)
	i, _, column := getFieldByJson(modelType, sortField)
	if i > -1 {
		return column
	}
	return ""
}
func getSortType(sortType string) string {
	if sortType == "-" {
		return desc
	} else {
		return asc
	}
}
func getFieldByJson(modelType reflect.Type, jsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 && strings.Split(tag1, ",")[0] == jsonName {
			if tag2, ok2 := field.Tag.Lookup("gorm"); ok2 {
				if has := strings.Contains(tag2, "column"); has {
					str1 := strings.Split(tag2, ";")
					num := len(str1)
					for k := 0; k < num; k++ {
						str2 := strings.Split(str1[k], ":")
						for j := 0; j < len(str2); j++ {
							if str2[j] == "column" {
								return i, field.Name, str2[j+1]
							}
						}
					}
				}
			}
			return i, field.Name, ""
		}
	}
	return -1, jsonName, jsonName
}
func Q(s string) string {
	if !(strings.HasPrefix(s, "%") && strings.HasSuffix(s, "%")) {
		return "%" + s + "%"
	} else if strings.HasPrefix(s, "%") {
		return s + "%"
	} else if strings.HasSuffix(s, "%") {
		return "%" + s
	}
	return s
}
func Prefix(s string) string {
	if strings.HasSuffix(s, "%") {
		return s
	} else {
		return s + "%"
	}
}
