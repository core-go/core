package search

import (
	"context"
	"reflect"
)

type NextSearchHandler struct {
	Find         func(ctx context.Context, filter interface{}, results interface{}, limit int64, nextPageToken string) (string, error)
	modelType    reflect.Type
	filterType   reflect.Type
	LogError     func(context.Context, string, ...map[string]interface{})
	List         string
	Next         string
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

func NewCSVNextSearchHandler(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler {
	return NewNextSearchHandlerWithLog(search, modelType, filterType, logError, writeLog, true, options...)
}
func NewNextSearchHandler(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler {
	return NewNextSearchHandlerWithLog(search, modelType, filterType, logError, writeLog, false, options...)
}
func NewNextSearchHandlerWithLog(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *NextSearchHandler {
	var list, next, resource, action, user string
	if len(options) > 0 && len(options[0]) > 0 {
		list = options[0]
	} else {
		list = "list"
	}
	if len(options) > 1 && len(options[1]) > 0 {
		next = options[1]
	} else {
		next = "next"
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
	return NewNextSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, quickSearch, list, next, resource, action, user, "")
}
func NewNextSearchHandlerWithQuickSearch(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, list string, total string, resource string, action string, userId string, embedField string) *NextSearchHandler {
	if len(action) == 0 {
		action = sSearch
	}
	paramIndex := BuildParamIndex(filterType)
	filterIndex := FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := BuildJsonMap(model, fields, embedField)
	return &NextSearchHandler{Find: search, modelType: modelType, filterType: filterType, List: list, Next: total, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
		JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes}
}
