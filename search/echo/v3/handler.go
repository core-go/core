package search

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	s "github.com/core-go/core/search"
	"github.com/labstack/echo"
)

type Search[T any, F any] func(ctx context.Context, filter F, limit int64, offset int64) ([]T, int64, error)
type SearchFn[T any, F any] func(ctx context.Context, filter F, limit int64, nextPageToken string) ([]T, string, error)

type SearchHandler[T any, F any] struct {
	Find         func(ctx context.Context, filter F, limit int64, offset int64) ([]T, int64, error)
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
	isPtr            bool
}

func NewCSVSearchHandler[T any, F any](search func(context.Context, F, int64, int64) ([]T, int64, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler[T, F] {
	return NewSearchHandlerWithLog[T, F](search, logError, writeLog, true, options...)
}
func NewSearchHandler[T any, F any](search func(context.Context, F, int64, int64) ([]T, int64, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SearchHandler[T, F] {
	return NewSearchHandlerWithLog[T, F](search, logError, writeLog, false, options...)
}
func NewSearchHandlerWithLog[T any, F any](search func(context.Context, F, int64, int64) ([]T, int64, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *SearchHandler[T, F] {
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
		var t T
		modelType := reflect.TypeOf(t)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		name := modelType.Name()
		resource = s.BuildResourceName(name)
	}
	if len(options) > 4 && len(options[4]) > 0 {
		action = options[4]
	} else {
		action = "search"
	}
	return NewSearchHandlerWithQuickSearch[T, F](search, logError, writeLog, quickSearch, list, total, resource, action, user, "")
}
func NewSearchHandlerWithQuickSearch[T any, F any](search func(context.Context, F, int64, int64) ([]T, int64, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, list string, total string, resource string, action string, userId string, embedField string) *SearchHandler[T, F] {
	if len(action) == 0 {
		action = "search"
	}
	var t T
	modelType := reflect.TypeOf(t)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	isPtr := false
	var f F
	filterType := reflect.TypeOf(f)
	if filterType.Kind() == reflect.Ptr {
		filterType = filterType.Elem()
		isPtr = true
	}
	paramIndex := s.BuildParamIndex(filterType)
	filterIndex := s.FindFilterIndex(filterType)
	model := reflect.New(modelType).Interface()
	fields := s.GetJSONFields(modelType)
	firstLayerIndexes, secondLayerIndexes := s.BuildJsonMap(model, fields, embedField)
	return &SearchHandler[T, F]{Find: search, filterType: filterType, List: list, Total: total, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
		JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes, isPtr: isPtr}
}

const internalServerError = "Internal Server Error"

func (c *SearchHandler[T, F]) Search(ctx echo.Context) error {
	r := ctx.Request()
	filter, x, er0 := s.BuildFilter(r, c.filterType, c.ParamIndex, c.userId, c.FilterIndex)
	if er0 != nil {
		return ctx.String(http.StatusBadRequest, "cannot decode filter: "+er0.Error())
	}
	limit, offset, fs, _, _, er1 := s.Extract(filter)
	if er1 != nil {
		return respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er1, c.WriteLog)
	}
	var ft F
	var ok bool
	if c.isPtr {
		ft, ok = filter.(F)
		if !ok {
			return ctx.String(http.StatusBadRequest, fmt.Sprintf("cannot cast filter %v", filter))
		}
	} else {
		mv := reflect.ValueOf(filter)
		pt := reflect.Indirect(mv).Interface()
		ft, ok = pt.(F)
		if !ok {
			return ctx.String(http.StatusBadRequest, fmt.Sprintf("cannot cast filter %v", filter))
		}
	}
	models, count, er2 := c.Find(r.Context(), ft, limit, offset)
	if er2 != nil {
		return respondError(ctx, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er2, c.WriteLog)
	}
	res := s.BuildResultMap(models, count, c.List, c.Total)
	if x == -1 {
		return respond(ctx, http.StatusOK, res, c.WriteLog, c.ResourceName, c.Activity, true, "")
	} else if c.CSV && x == 1 {
		resCSV, ok := s.ResultToCsv(fs, models, count, c.embedField, c.JsonMap, c.SecondaryJsonMap)
		if ok {
			return respond(ctx, http.StatusOK, resCSV, c.WriteLog, c.ResourceName, c.Activity, true, "")
		} else {
			return respond(ctx, http.StatusOK, res, c.WriteLog, c.ResourceName, c.Activity, true, "")
		}
	} else {
		return respond(ctx, http.StatusOK, res, c.WriteLog, c.ResourceName, c.Activity, true, "")
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
