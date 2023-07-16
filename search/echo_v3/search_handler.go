package echo

import (
	"context"
	s "github.com/core-go/core/search"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
)

const sSearch = "search"

type SearchHandler struct {
	Find         func(ctx context.Context, filter interface{}, results interface{}, limit int64, options int64) (int64, error)
	modelType    reflect.Type
	filterType   reflect.Type
	LogError     func(context.Context, string,...map[string]interface{})
	Config       s.SearchResultConfig
	CSV          bool
	Log          func(ctx context.Context, resource string, action string, success bool, desc string) error
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

func NewSearchHandler(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string,...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, true, options...)
}
func NewJSONSearchHandler(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string,...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithQuickSearch(search, modelType, filterType, logError, writeLog, false, options...)
}
func NewSearchHandlerWithQuickSearch(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string,...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *SearchHandler {
	var resource, action, user string
	if len(options) > 0 && len(options[0]) > 0 {
		user = options[0]
	} else {
		user = s.UserId
	}
	if len(options) > 1 && len(options[1]) > 0 {
		resource = options[1]
	} else {
		name := modelType.Name()
		resource = s.BuildResourceName(name)
	}
	if len(options) > 2 && len(options[2]) > 0 {
		action = options[2]
	} else {
		action = sSearch
	}
	return NewSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, action, user, "")
}
func NewSearchHandlerWithUserId(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string,...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithUserIdAndQuickSearch(search, modelType, filterType, userId, logError, writeLog, true, options...)
}
func NewJSONSearchHandlerWithUserId(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string,...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler {
	return NewSearchHandlerWithUserIdAndQuickSearch(search, modelType, filterType, userId, logError, writeLog, false, options...)
}
func NewSearchHandlerWithUserIdAndQuickSearch(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, userId string, logError func(context.Context, string,...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *SearchHandler {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		resource = options[0]
	} else {
		name := modelType.Name()
		resource = s.BuildResourceName(name)
	}
	if len(options) > 1 && len(options[1]) > 0 {
		action = options[1]
	} else {
		action = sSearch
	}
	return NewSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, action, userId, "")
}
func NewDefaultSearchHandler(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, resource string, logError func(context.Context, string,...map[string]interface{}), userId string, quickSearch bool, writeLog func(context.Context, string, string, bool, string) error) *SearchHandler {
	return NewSearchHandlerWithConfig(search, modelType, filterType, logError, nil, writeLog, quickSearch, resource, sSearch, userId, "")
}
func NewSearchHandlerWithConfig(search func(context.Context, interface{}, interface{}, int64, int64) (int64, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string,...map[string]interface{}), config *s.SearchResultConfig, writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, resource string, action string, userId string, embedField string) *SearchHandler {
	var c s.SearchResultConfig
	if len(action) == 0 {
		action = sSearch
	}
	if config != nil {
		c = *config
	} else {
		c.Results = "results"
		c.Total = "total"
	}
	paramIndex := s.BuildParamIndex(filterType)
	filterIndex := s.FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := s.GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := s.BuildJsonMap(model, fields, embedField)
	return &SearchHandler{Find: search, modelType: modelType, filterType: filterType, Config: c, Log: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError, JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes}
}

const internalServerError = "Internal Server Error"

func (c *SearchHandler) Search(ctx echo.Context) error {
	r := ctx.Request()
	filter, x, er0 := s.BuildFilter(r, c.filterType, c.ParamIndex, c.userId, c.FilterIndex)
	if er0 != nil {
		return ctx.String(http.StatusBadRequest, "cannot parse form: "+"cannot decode filter: "+er0.Error())
	}
	limit, offset, fs, _, _, er1 := s.Extract(filter)
	if er1 != nil {
		return respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er1, c.Log)
	}
	modelsType := reflect.Zero(reflect.SliceOf(c.modelType)).Type()
	models := reflect.New(modelsType).Interface()
	count, er2 := c.Find(r.Context(), filter, models, limit, offset)
	if er2 != nil {
		return respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er2, c.Log)
	}

	result := s.BuildResultMap(models, count, c.Config)
	if x == -1 {
		return succeed(ctx, http.StatusOK, result, c.Log, c.ResourceName, c.Activity)
	} else if c.CSV && x == 1 {
		result1, ok := s.ResultToCsv(fs, models, count, c.embedField)
		if ok {
			return succeed(ctx, http.StatusOK, result1, c.Log, c.ResourceName, c.Activity)
		} else {
			return succeed(ctx, http.StatusOK, result, c.Log, c.ResourceName, c.Activity)
		}
	} else {
		return succeed(ctx, http.StatusOK, result, c.Log, c.ResourceName, c.Activity)
	}
}
type NextSearchHandler struct {
	Find         func(ctx context.Context, filter interface{}, results interface{}, limit int64, nextPageToken string) (string, error)
	modelType    reflect.Type
	filterType   reflect.Type
	LogError     func(context.Context, string, ...map[string]interface{})
	Config       s.SearchResultConfig
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
		user = s.UserId
	}
	if len(options) > 1 && len(options[1]) > 0 {
		resource = options[1]
	} else {
		name := modelType.Name()
		resource = s.BuildResourceName(name)
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
		resource = s.BuildResourceName(name)
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
func NewNextSearchHandlerWithConfig(search func(context.Context, interface{}, interface{}, int64, string) (string, error), modelType reflect.Type, filterType reflect.Type, logError func(context.Context, string, ...map[string]interface{}), config *s.SearchResultConfig, writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, resource string, action string, userId string, embedField string) *NextSearchHandler {
	var c s.SearchResultConfig
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

	paramIndex := s.BuildParamIndex(filterType)
	filterIndex := s.FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := s.GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := s.BuildJsonMap(model, fields, embedField)
	return &NextSearchHandler{Find: search, modelType: modelType, filterType: filterType, Config: c, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
		JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes}
}
func (c *NextSearchHandler) Search(ctx echo.Context) error {
	r := ctx.Request()
	filter, x, er0 := s.BuildFilter(r, c.filterType, c.ParamIndex, c.userId, c.FilterIndex)
	if er0 != nil {
		ctx.String(http.StatusBadRequest, "cannot parse form: "+"cannot decode filter: "+er0.Error())
		return er0
	}
	limit, _, fs, _, nextPageToken, er1 := s.Extract(filter)
	if er1 != nil {
		respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er1, c.WriteLog)
		return er1
	}
	modelsType := reflect.Zero(reflect.SliceOf(c.modelType)).Type()
	models := reflect.New(modelsType).Interface()
	nx, er2 := c.Find(r.Context(), filter, models, limit, nextPageToken)
	if er2 != nil {
		respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er2, c.WriteLog)
		return er2
	}

	result := s.BuildNextResultMap(models, nx, c.Config)
	if x == -1 {
		return succeed(ctx, http.StatusOK, result, c.WriteLog, c.ResourceName, c.Activity)
	} else if c.CSV && x == 1 {
		result1, ok := s.ResultToNextCsv(fs, models, nx, c.embedField)
		if ok {
			return succeed(ctx, http.StatusOK, result1, c.WriteLog, c.ResourceName, c.Activity)
		} else {
			return succeed(ctx, http.StatusOK, result, c.WriteLog, c.ResourceName, c.Activity)
		}
	} else {
		return succeed(ctx, http.StatusOK, result, c.WriteLog, c.ResourceName, c.Activity)
	}
}

func respondError(ctx echo.Context, code int, result interface{}, logError func(context.Context, string,...map[string]interface{}), resource string, action string, err error, writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error) error {
	if logError != nil {
		logError(ctx.Request().Context(), err.Error())
	}
	return respond(ctx, code, result, writeLog, resource, action, false, err.Error())
}
func respond(ctx echo.Context, code int, result interface{}, writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error, resource string, action string, success bool, desc string) error {
	err := ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request().Context(), resource, action, success, desc)
	}
	return err
}
func succeed(ctx echo.Context, code int, result interface{}, writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error, resource string, action string) error {
	return respond(ctx, code, result, writeLog, resource, action, true, "")
}
