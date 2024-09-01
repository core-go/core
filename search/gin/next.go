package search

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	s "github.com/core-go/core/search"
)

type NextSearchHandler[T any, F any] struct {
	Find         func(ctx context.Context, filter F, limit int64, offset string) ([]T, string, error)
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
	isPtr            bool
}

func NewCSVNextSearchHandler[T any, F any](search func(context.Context, F, int64, string) ([]T, string, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler[T, F] {
	return NewNextSearchHandlerWithLog[T, F](search, logError, writeLog, true, options...)
}
func NewNextSearchHandler[T any, F any](search func(context.Context, F, int64, string) ([]T, string, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *NextSearchHandler[T, F] {
	return NewNextSearchHandlerWithLog[T, F](search, logError, writeLog, false, options...)
}
func NewNextSearchHandlerWithLog[T any, F any](search func(context.Context, F, int64, string) ([]T, string, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, options ...string) *NextSearchHandler[T, F] {
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
	return NewNextSearchHandlerWithQuickSearch[T, F](search, logError, writeLog, quickSearch, list, next, resource, action, user, "")
}
func NewNextSearchHandlerWithQuickSearch[T any, F any](search func(context.Context, F, int64, string) ([]T, string, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, quickSearch bool, list string, total string, resource string, action string, userId string, embedField string) *NextSearchHandler[T, F] {
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
	return &NextSearchHandler[T, F]{Find: search, filterType: filterType, List: list, Next: total, WriteLog: writeLog, CSV: quickSearch, ResourceName: resource, Activity: action, ParamIndex: paramIndex, FilterIndex: filterIndex, userId: userId, embedField: embedField, LogError: logError,
		JsonMap: firstLayerIndexes, SecondaryJsonMap: secondLayerIndexes, isPtr: isPtr}
}

func (c *NextSearchHandler[T, F]) Search(w http.ResponseWriter, r *http.Request) {
	filter, x, er0 := s.BuildFilter(r, c.filterType, c.ParamIndex, c.userId, c.FilterIndex)
	if er0 != nil {
		http.Error(w, "cannot decode filter: "+er0.Error(), http.StatusBadRequest)
		return
	}
	limit, _, fs, _, nextPageToken, er1 := s.Extract(filter)
	if er1 != nil {
		s.RespondError(w, r, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er1, c.WriteLog)
		return
	}
	var ft F
	var ok bool
	if c.isPtr {
		ft, ok = filter.(F)
		if !ok {
			http.Error(w, fmt.Sprintf("cannot cast filter %v", filter), http.StatusBadRequest)
			return
		}
	} else {
		mv := reflect.ValueOf(filter)
		pt := reflect.Indirect(mv).Interface()
		ft, ok = pt.(F)
		if !ok {
			http.Error(w, fmt.Sprintf("cannot cast filter %v", filter), http.StatusBadRequest)
			return
		}
	}
	models, next, er2 := c.Find(r.Context(), ft, limit, nextPageToken)
	if er2 != nil {
		s.RespondError(w, r, http.StatusInternalServerError, internalServerError, c.LogError, c.ResourceName, c.Activity, er2, c.WriteLog)
		return
	}
	res := s.BuildNextResultMap(models, next, c.List, c.Next)
	if x == -1 {
		s.Respond(w, r, http.StatusOK, res, c.WriteLog, c.ResourceName, c.Activity, true, "")
	} else if c.CSV && x == 1 {
		resCSV, ok := s.ResultToNextCsv(fs, models, next, c.embedField, c.JsonMap, c.SecondaryJsonMap)
		if ok {
			s.Respond(w, r, http.StatusOK, resCSV, c.WriteLog, c.ResourceName, c.Activity, true, "")
		} else {
			s.Respond(w, r, http.StatusOK, res, c.WriteLog, c.ResourceName, c.Activity, true, "")
		}
	} else {
		s.Respond(w, r, http.StatusOK, res, c.WriteLog, c.ResourceName, c.Activity, true, "")
	}
}
