package search

import (
	"context"
	"net/http"
	"reflect"

	s "github.com/core-go/core/search"
	"github.com/labstack/echo"
)

const sSearch = "search"

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
		user = s.UserId
	}
	if len(options) > 3 && len(options[3]) > 0 {
		resource = options[3]
	} else {
		name := modelType.Name()
		resource = s.BuildResourceName(name)
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
	paramIndex := s.BuildParamIndex(filterType)
	filterIndex := s.FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := s.GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := s.BuildJsonMap(model, fields, embedField)
	return &SearchHandler{Find: search, modelType: modelType, filterType: filterType, List: list, Total: total, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
		JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes}
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
		return respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er1, c.WriteLog)
	}
	modelsType := reflect.Zero(reflect.SliceOf(c.modelType)).Type()
	models := reflect.New(modelsType).Interface()
	count, er2 := c.Find(r.Context(), filter, models, limit, offset)
	if er2 != nil {
		return respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er2, c.WriteLog)
	}

	result := s.BuildResultMap(models, count, c.List, c.Total)
	if x == -1 {
		return succeed(ctx, http.StatusOK, result, c.WriteLog, c.ResourceName, c.Activity)
	} else if c.CSV && x == 1 {
		result1, ok := s.ResultToCsv(fs, models, count, c.embedField)
		if ok {
			return succeed(ctx, http.StatusOK, result1, c.WriteLog, c.ResourceName, c.Activity)
		} else {
			return succeed(ctx, http.StatusOK, result, c.WriteLog, c.ResourceName, c.Activity)
		}
	} else {
		return succeed(ctx, http.StatusOK, result, c.WriteLog, c.ResourceName, c.Activity)
	}
}

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
		user = s.UserId
	}
	if len(options) > 3 && len(options[3]) > 0 {
		resource = options[3]
	} else {
		name := modelType.Name()
		resource = s.BuildResourceName(name)
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
	paramIndex := s.BuildParamIndex(filterType)
	filterIndex := s.FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := s.GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := s.BuildJsonMap(model, fields, embedField)
	return &NextSearchHandler{Find: search, modelType: modelType, filterType: filterType, List: list, Next: total, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
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

	result := s.BuildNextResultMap(models, nx, c.List, c.Next)
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

func respondError(ctx echo.Context, code int, result interface{}, logError func(context.Context, string, ...map[string]interface{}), resource string, action string, err error, writeLog func(ctx context.Context, resource string, action string, success bool, desc string) error) error {
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
