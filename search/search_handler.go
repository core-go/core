package search

import (
	"context"
	"reflect"
)

type SearchHandler struct {
	search       func(ctx context.Context, filter interface{}, results interface{}, limit int64, options ...int64) (int64, string, error)
	modelType    reflect.Type
	filterType   reflect.Type
	LogError     func(context.Context, string)
	Config       SearchResultConfig
	CSV          bool
	WriteLog     func(ctx context.Context, resource string, action string, success bool, desc string) error
	ResourceName string
	Activity     string
	embedField   string
	userId       string
	// search by GET
	paramIndex  map[string]int
	filterIndex int
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
	Uid                = "uid"
	Username           = "username"
	sSearch            = "search"
)

func NewSearchHandler(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, true, options...)
}
func NewJSONSearchHandler(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, false, options...)
}
func NewSearchHandlerWithQuickSearch(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *SearchHandler {
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
	return NewSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, action, user, "")
}
func NewSearchHandlerWithUserId(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithUserIdAndQuickSearch(search, modelType, filterType, userId, logError, writeLog, true, options...)
}
func NewJSONSearchHandlerWithUserId(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithUserIdAndQuickSearch(search, modelType, filterType, userId, logError, writeLog, false, options...)
}
func NewSearchHandlerWithUserIdAndQuickSearch(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *SearchHandler {
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
	return NewSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, action, userId, "")
}
func NewDefaultSearchHandler(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, resource string, logError func(context.Context, string), userId string, quickSearch bool, writeLog func(context.Context, string, string, bool, string) error) *SearchHandler {
	return NewSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, sSearch, userId, "")
}
func NewSearchHandlerWithConfig(search func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string), config *SearchResultConfig, writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, resource string, action string, userId string, embedField string) *SearchHandler {
	var c SearchResultConfig
	if len(action) == 0 {
		action = sSearch
	}
	if config != nil {
		c = *config
	} else {
		c.LastPage = "last"
		c.Results = "list"
		c.Total = "total"
		c.NextPageToken = "nextPageToken"
	}

	paramIndex := BuildParamIndex(filterType)
	filterIndex := FindFilterIndex(filterType)

	return &SearchHandler{search: search, modelType: modelType, filterType: filterType, Config: c, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, paramIndex: paramIndex, filterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError}
}
