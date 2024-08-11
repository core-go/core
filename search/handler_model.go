package search

import "reflect"

func BuildResultMap(models interface{}, count int64, list string, total string) map[string]interface{} {
	result := make(map[string]interface{})
	result[total] = count
	result[list] = models
	return result
}
func BuildNextResultMap(models interface{}, nextPageToken string, list string, next string) map[string]interface{} {
	result := make(map[string]interface{})
	result[list] = models
	if len(nextPageToken) > 0 {
		result[next] = nextPageToken
	}
	return result
}
func SetUserId(sm interface{}, currentUserId string) {
	if s, ok := sm.(*Filter); ok { // Is Filter struct
		RepairFilter(s, currentUserId)
	} else { // Is extended from Filter struct
		value := reflect.Indirect(reflect.ValueOf(sm))
		numField := value.NumField()
		for i := 0; i < numField; i++ {
			// Find Filter field of extended struct
			if s, ok := value.Field(i).Interface().(*Filter); ok {
				RepairFilter(s, currentUserId)
				break
			}
		}
	}
}
func CreateFilter(filterType reflect.Type, options ...int) interface{} {
	filterIndex := -1
	if len(options) > 0 && options[0] >= 0 {
		filterIndex = options[0]
	}
	var filter = reflect.New(filterType).Interface()
	if filterIndex >= 0 {
		value := reflect.Indirect(reflect.ValueOf(filter))
		if filterIndex < value.NumField() {
			f := value.Field(filterIndex)
			if _, ok := f.Interface().(*Filter); ok {
				f.Set(reflect.ValueOf(&Filter{}))
			} else if _, ok := f.Interface().(Filter); ok {
				f.Set(reflect.ValueOf(Filter{}))
			}
		}
	}
	return filter
}

func FindFilterIndex(filterType reflect.Type) int {
	numField := filterType.NumField()
	t := reflect.TypeOf(&Filter{})
	for i := 0; i < numField; i++ {
		if filterType.Field(i).Type == t {
			return i
		}
	}
	return -1
}

// Check valid and change value of pagination to correct
func RepairFilter(filter *Filter, options ...string) {
	if len(options) > 0 {
		filter.CurrentUserId = options[0]
	}

	if filter.PageIndex != 0 && filter.Page == 0 {
		filter.Page = filter.PageIndex
	}
	if filter.PageSize != 0 && filter.Limit == 0 {
		filter.Limit = filter.PageSize
	}
	if filter.FirstPageSize != 0 && filter.FirstLimit == 0 {
		filter.FirstLimit = filter.FirstPageSize
	}

	pageSize := filter.Limit
	if pageSize > MaxPageSizeDefault {
		pageSize = PageSizeDefault
	}

	pageIndex := filter.Page
	if filter.Page < 1 {
		pageIndex = 1
	}

	if filter.Limit != pageSize {
		filter.Limit = pageSize
	}

	if filter.Page != pageIndex {
		filter.Page = pageIndex
	}
}
