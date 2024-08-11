package search

import (
	"context"
	"reflect"
	"strings"
)

type SearchHandler struct {
	Find         func(ctx context.Context, filter interface{}, results interface{}, limit int64, offset int64) (int64, error)
	modelType    reflect.Type
	filterType   reflect.Type
	LogError     func(context.Context, string, ...map[string]interface{})
	List         string
	Total        string
	CSV          bool
	WriteLog     func(ctx context.Context, resource string, action string, success bool, desc string) error
	ResourceName string
	Activity     string
	embedField   string
	userId       string
	// search by GET
	ParamIndex       map[string]int
	FilterIndex      int
	JsonMap          map[string]int
	SecondaryJsonMap map[string]int
}

var filterParamIndex map[string]int

func GetFilterParamIndex() map[string]int {
	if filterParamIndex == nil || len(filterParamIndex) == 0 {
		f := BuildParamIndex(reflect.TypeOf(Filter{}))
		filterParamIndex = f
	}
	return filterParamIndex
}

const (
	PageSizeDefault    = 10
	MaxPageSizeDefault = 10000
	UserId             = "userId"
	sSearch            = "search"
)

func NewCSVSearchHandler(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithLog(search, modelType, filterType, logError, writeLog, true, options...)
}
func NewSearchHandler(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithLog(search, modelType, filterType, logError, writeLog, false, options...)
}
func NewSearchHandlerWithLog(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *SearchHandler {
	var list, total, resource, action, user string
	if len(options) > 0 && len(options[0]) > 0 {
		list = options[0]
	} else {
		list = "list"
	}
	if len(options) > 1 && len(options[1]) > 0 {
		total = options[1]
	} else {
		total = "total"
	}
	if len(options) > 2 && len(options[2]) > 0 {
		user = options[2]
	} else {
		user = UserId
	}
	if len(options) > 3 && len(options[3]) > 0 {
		resource = options[3]
	} else {
		name := modelType.Name()
		resource = BuildResourceName(name)
	}
	if len(options) > 4 && len(options[4]) > 0 {
		action = options[4]
	} else {
		action = sSearch
	}
	return NewSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, quickSearch, list, total, resource, action, user, "")
}
func NewSearchHandlerWithQuickSearch(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, list string, total string, resource string, action string, userId string, embedField string) *SearchHandler {
	if len(action) == 0 {
		action = sSearch
	}
	paramIndex := BuildParamIndex(filterType)
	filterIndex := FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := BuildJsonMap(model, fields, embedField)
	return &SearchHandler{Find: search, modelType: modelType, filterType: filterType, List: list, Total: total, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
		JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes}
}
func GetJSONFields(modelType reflect.Type) []string {
	numField := modelType.NumField()
	fields := make([]string, 0)
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag, ok := field.Tag.Lookup("json")
		if ok {
			name := strings.Split(tag, ",")[0]
			fields = append(fields, name)
		}
	}
	return fields
}
