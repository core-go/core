package search

import (
	"context"
	"reflect"
)

type Search func(ctx context.Context, filter interface{}, results interface{}, limit int64, offset int64) (int64, error)
type SearchFn func(ctx context.Context, filter interface{}, results interface{}, limit int64, nextPageToken string) (string, error)

type NextSearchHandler struct {
	Find         func(ctx context.Context, filter interface{}, results interface{}, limit int64, nextPageToken string) (string, error)
	modelType    reflect.Type
	filterType   reflect.Type
	LogError     func(context.Context, string, ...map[string]interface{})
	Config       SearchResultConfig
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

func NewNextSearchHandler(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler {
	return NewNextSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, true, options...)
}
func NewJSONNextSearchHandler(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler {
	return NewNextSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, false, options...)
}
func NewNextSearchHandlerWithQuickSearch(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *NextSearchHandler {
	var resource, action, user string
	if len(options) > 0 && len(options[0]) > 0 {
		user = options[0]
	} else {
		user = UserId
	}
	if len(options) > 1 && len(options[1]) > 0 {
		resource = options[1]
	} else {
		name := modelType.Name()
		resource = BuildResourceName(name)
	}
	if len(options) > 2 && len(options[2]) > 0 {
		action = options[2]
	} else {
		action = sSearch
	}
	return NewNextSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, action, user, "")
}
func NewNextSearchHandlerWithUserId(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler {
	return NewNextSearchHandlerWithUserIdAndQuickSearch(search, modelType, filterType, userId, logError, writeLog, true, options...)
}
func NewJSONNextSearchHandlerWithUserId(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler {
	return NewNextSearchHandlerWithUserIdAndQuickSearch(search, modelType, filterType, userId, logError, writeLog, false, options...)
}
func NewNextSearchHandlerWithUserIdAndQuickSearch(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *NextSearchHandler {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	} else {
		name := modelType.Name()
		resource = BuildResourceName(name)
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	} else {
		action = sSearch
	}
	return NewNextSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, action, userId, "")
}
func NewDefaultNextSearchHandler(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, resource string, logError func(context.Context, string, ...map[string]interface{}), userId string, quickSearch bool, writeLog func(context.Context, string, string, bool, string) error) *NextSearchHandler {
	return NewNextSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, sSearch, userId, "")
}
func NewNextSearchHandlerWithConfig(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), config *SearchResultConfig, writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, resource string, action string, userId string, embedField string) *NextSearchHandler {
	var c SearchResultConfig
	if len(action) == 0 {
		action = sSearch
	}
	if config != nil {
		c = *config
	} else {
		// c.LastPage = "last"
		c.Results = "list"
		c.Total = "total"
		c.Next = "next"
	}

	paramIndex := BuildParamIndex(filterType)
	filterIndex := FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := BuildJsonMap(model, fields, embedField)
	return &NextSearchHandler{Find: search, modelType: modelType, filterType: filterType, Config: c, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
		JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes}
}
