package search

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func BuildResourceName(s string) string {
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if s2[i] != s[i] {
			s3 += "-" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}
	if string(s3[0]) == "-" || string(s3[0]) == "_" {
		return s3[1:]
	}
	return s3
}
func UrlToModel(filter interface{}, params url.Values, paramIndex map[string]int, options...int) interface{} {
	value := reflect.Indirect(reflect.ValueOf(filter))
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}

	//TimeRange
	timeType := reflect.TypeOf(TimeRange{})
	tRange := reflect.New(timeType)
	psTime := tRange.Elem()

	//Int64Range
	int64Type := reflect.TypeOf(Int64Range{})
	i64 := reflect.New(int64Type)
	psI64 := i64.Elem()

	//NumRange
	nRangeType := reflect.TypeOf(NumberRange{})
	nRange := reflect.New(nRangeType)
	psNRange := nRange.Elem()

	for paramKey, valueArr := range params {
		paramValue := ""
		if len(valueArr) > 0 {
			paramValue = valueArr[0]
		}
		if field, err := FindField(value, paramKey, paramIndex, options...); err == nil {
			kind := field.Kind()

			var v interface{}
			// Need handle more case of kind
			if kind == reflect.String {
				v = paramValue
			} else if kind == reflect.Int {
				v, _ = strconv.Atoi(paramValue)
			} else if kind == reflect.Int64 {
				v, _ = strconv.ParseInt(paramValue, 10, 64)
			} else if kind == reflect.Float64 {
				v, _ = strconv.ParseFloat(paramValue, 64)
			} else if kind == reflect.Bool{
				v, _ = strconv.ParseBool(paramValue)
			} else if kind == reflect.Slice {
				sliceKind := reflect.TypeOf(field.Interface()).Elem().Kind()
				if sliceKind == reflect.String {
					v = strings.Split(paramValue, ",")
				} else {
					log.Println("Unhandled slice kind:", kind)
					continue
				}
			} else if kind == reflect.Struct {
				newModel := reflect.New(reflect.Indirect(field).Type()).Interface()
				if errDecode := json.Unmarshal([]byte(paramValue), newModel); errDecode != nil {
					panic(errDecode)
				}
				v = newModel
			} else if kind == reflect.Ptr {
				ptrKind := field.Interface()
				switch ptrKind.(type) {
				case *string:
					field.Set(reflect.ValueOf(&paramValue))
					continue
				case *int:
					iv, er := strconv.Atoi(paramValue)
					if er == nil {
						field.Set(reflect.ValueOf(&iv))
					}
					continue
				case *int64:
					iv, er := strconv.ParseInt(paramValue, 10, 64)
					if er == nil {
						field.Set(reflect.ValueOf(&iv))
					}
					continue
				case *float64:
					iv, er := strconv.ParseFloat(paramValue, 64)
					if er == nil {
						field.Set(reflect.ValueOf(&iv))
					}
					continue
				case *bool:
					iv, er := strconv.ParseBool(paramValue)
					if er == nil {
						field.Set(reflect.ValueOf(&iv))
					}
					continue
				case *TimeRange:
					keys := strings.Split(paramKey,".")
					f := psTime.FieldByName(strings.Title(keys[1]))
					tValue, _ := time.Parse(time.RFC3339, paramValue)
					f.Set(reflect.ValueOf(&tValue))
					field.Set(reflect.ValueOf(tRange.Interface()))
					continue
				case *Int64Range:
					keys := strings.Split(paramKey, ".")
					f := psI64.FieldByName(strings.Title(keys[1]))
					i64Value, _ := strconv.ParseInt(paramValue, 10, 64)
					f.Set(reflect.ValueOf(&i64Value))
					field.Set(reflect.ValueOf(i64.Interface()))
					continue
				case *NumberRange:
					keys := strings.Split(paramKey, ".")
					f := psNRange.FieldByName(strings.Title(keys[1]))
					nValue, _ := strconv.ParseFloat(paramValue, 64)
					f.Set(reflect.ValueOf(&nValue))
					field.Set(reflect.ValueOf(nRange.Interface()))
					continue
				}
			} else {
				log.Println("Unhandled kind:", kind)
				continue
			}
			field.Set(reflect.Indirect(reflect.ValueOf(v)))
		} else {
			log.Println(err)
		}
	}
	return filter
}
func FindField(value reflect.Value, paramKey string, paramIndex map[string]int, options...int) (reflect.Value, error) {
	if keys :=strings.Split(paramKey,"."); len(keys) > 0 {
		paramKey = keys[0]
	}
	if index, ok := paramIndex[paramKey]; ok {
		return value.Field(index), nil
	}
	filterIndex := -1
	if len(options) > 0 && options[0] >= 0 {
		filterIndex = options[0]
	}
	if filterIndex >= 0 {
		filterParamIndex := GetFilterParamIndex()
		if index, ok := filterParamIndex[paramKey]; ok {
			filterField := value.Field(filterIndex)
			if filterField.Kind() == reflect.Ptr {
				filterField = reflect.Indirect(filterField)
			}
			return filterField.Field(index), nil
		}
	}
	return value, errors.New("can't find field " + paramKey)
}
func BuildParamIndex(filterType reflect.Type) map[string]int {
	params := map[string]int{}
	numField := filterType.NumField()
	for i := 0; i < numField; i++ {
		field := filterType.Field(i)
		fullJsonTag := field.Tag.Get("json")
		tagDetails := strings.Split(fullJsonTag, ",")
		if len(tagDetails) > 0 && len(tagDetails[0]) > 0 {
			params[tagDetails[0]] = i
		}
	}
	return params
}
func CreateParams(filterType reflect.Type, modelType reflect.Type, opts...string) (map[string]int, int, map[string]int, map[string]int) {
	embedField := ""
	if len(opts) > 0 {
		embedField = opts[0]
	}
	paramIndex := BuildParamIndex(filterType)
	filterIndex := FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := BuildJsonMap(model, fields, embedField)
	return paramIndex, filterIndex, firstLayerIndexes, secondLayerIndexes
}
func BuildParams(filterType reflect.Type) (map[string]int, int) {
	paramIndex := BuildParamIndex(filterType)
	filterIndex := FindFilterIndex(filterType)
	return paramIndex, filterIndex
}
func BuildFilter(r *http.Request, filterType reflect.Type, paramIndex map[string]int, userIdName string, options...int) (interface{}, int, error) {
	var filter = CreateFilter(filterType, options...)
	method := r.Method
	x := 1
	if method == http.MethodGet {
		ps := r.URL.Query()
		fs := ps.Get("fields")
		if len(fs) == 0 {
			x = -1
		}
		UrlToModel(filter, ps, paramIndex, options...)
	} else if method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			return nil, x, err
		}
	}
	userId := ""
	if len(userId) == 0 {
		u := r.Context().Value(userIdName)
		if u != nil {
			u2, ok2 := u.(string)
			if ok2 {
				userId = u2
			}
		}
	}
	SetUserId(filter, userId)
	return filter, x, nil
}
var userId = "userId"
func ApplyUserId(str string) {
	userId = str
}
func GetUser(ctx context.Context, opt...string) string {
	user := userId
	if len(opt) > 0 && len(opt[0]) > 0 {
		user = opt[0]
	}
	u := ctx.Value(user)
	if u != nil {
		u2, ok2 := u.(string)
		if ok2 {
			return u2
		}
	}
	return ""
}
func Decode(r *http.Request, filter interface{}, paramIndex map[string]int, options...int) error {
	method := r.Method
	if method == http.MethodGet {
		ps := r.URL.Query()
		UrlToModel(filter, ps, paramIndex, options...)
		return nil
	} else {
		err := json.NewDecoder(r.Body).Decode(&filter)
		return err
	}
}
func ToFilter(w http.ResponseWriter, r *http.Request, filter interface{}, paramIndex map[string]int, options...int) error {
	err := Decode(r, &filter, paramIndex, options...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return err
}
func DecodeAndCheck(w http.ResponseWriter, r *http.Request, filter interface{}, paramIndex map[string]int, options...int) error {
	return ToFilter(w, r, filter, paramIndex, options...)
}
func ResultCsv(fields []string, models interface{}, count int64, opts...map[string]int) (string, bool) {
	if len(fields) > 0 {
		result1 := ToCsv(fields, models, count, "", opts...)
		return result1, true
	} else {
		return "", false
	}
}
func ResultToCsv(fields []string, models interface{}, count int64, embedField string, opts...map[string]int) (string, bool) {
	if len(fields) > 0 {
		result1 := ToCsv(fields, models, count, embedField, opts...)
		return result1, true
	} else {
		return "", false
	}
}
func ResultNextCsv(fields []string, models interface{}, nextPageToken string, opts...map[string]int) (string, bool) {
	if len(fields) > 0 {
		result1 := ToNextCsv(fields, models, nextPageToken, "", opts...)
		return result1, true
	} else {
		return "", false
	}
}
func ResultToNextCsv(fields []string, models interface{}, nextPageToken string, embedField string, opts...map[string]int) (string, bool) {
	if len(fields) > 0 {
		result1 := ToNextCsv(fields, models, nextPageToken, embedField, opts...)
		return result1, true
	} else {
		return "", false
	}
}
func CSV(w http.ResponseWriter, code int, out string)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(out))
}
