package search

import (
	"errors"
	"reflect"
)

type Extractor struct {
	Page       string
	Limit      string
	FirstLimit string
}

func NewExtractor(options ...string) *Extractor {
	var page, limit, firstLimit string
	if len(options) >= 1 && len(options[0]) > 0 {
		page = options[0]
	} else {
		page = "Page"
	}
	if len(options) >= 2 && len(options[1]) > 0 {
		limit = options[1]
	} else {
		limit = "Limit"
	}
	if len(options) >= 3 && len(options[2]) > 0 {
		firstLimit = options[2]
	} else {
		firstLimit = "FirstLimit"
	}
	return &Extractor{Page: page, Limit: limit, FirstLimit: firstLimit}
}

func (e *Extractor) Extract(m interface{}) (int64, int64, int64, error) {
	if sModel, ok0 := m.(*Filter); ok0 {
		return sModel.Page, sModel.Limit, sModel.FirstLimit, nil
	}
	var page, limit, firstLimit int64
	page = -1
	limit = -1
	firstLimit = -1
	value := reflect.Indirect(reflect.ValueOf(m))
	t := value.Type()
	numField := t.NumField()
	// numField := value.NumField()
	for i := 0; i < numField; i++ {
		if sModel1, ok1 := value.Field(i).Interface().(*Filter); ok1 {
			return sModel1.Page, sModel1.Limit, sModel1.FirstLimit, nil
		} else {
			n := t.Field(i).Name
			if n == e.Page {
				if p, ok := value.Field(i).Interface().(int64); ok {
					page = p
				}
			} else if n == e.Limit {
				if p, ok := value.Field(i).Interface().(int64); ok {
					limit = p
				}
			} else if n == e.FirstLimit {
				if p, ok := value.Field(i).Interface().(int64); ok {
					firstLimit = p
				}
			}
			if page >= 0 && limit >= 0 && firstLimit >= 0 {
				return page, limit, firstLimit, nil
			}
		}
	}
	return page, limit, firstLimit, nil
}

func Extract(m interface{}) (int64, int64, []string, string, string, error) {
	if sModel, ok := m.(*Filter); ok {
		var limit, offset int64
		if sModel.FirstLimit > 0 {
			if sModel.Page == 1 {
				limit = sModel.FirstLimit
				offset = 0
			} else {
				limit = sModel.Limit
				offset = sModel.Limit*(sModel.Page-2) + sModel.FirstLimit
			}
		} else {
			limit = sModel.Limit
			offset = sModel.Limit * (sModel.Page - 1)
		}
		return limit, offset, sModel.Fields, sModel.Sort, sModel.RefId, nil
	} else {
		value := reflect.Indirect(reflect.ValueOf(m))
		numField := value.NumField()
		for i := 0; i < numField; i++ {
			if sModel1, ok := value.Field(i).Interface().(*Filter); ok {
				var limit1, offset1 int64
				if sModel1.FirstLimit > 0 {
					if sModel1.Page == 1 {
						limit1 = sModel1.FirstLimit
						offset1 = 0
					} else {
						limit1 = sModel1.Limit
						offset1 = sModel1.Limit*(sModel1.Page-2) + sModel1.FirstLimit
					}
				} else {
					limit1 = sModel1.Limit
					offset1 = sModel1.Limit * (sModel1.Page - 1)
				}
				return limit1, offset1, sModel1.Fields, sModel1.Sort, sModel1.RefId, nil
			}
		}
		return 0, 0, nil, "", "", errors.New("cannot extract sort, pageIndex, pageSize, firstPageSize from model")
	}
}
func GetFieldsAndSort(m interface{}) ([]string, string) {
	f, s, _ := GetFieldsAndSortAndRefId(m)
	return f, s
}
func GetFieldsAndRefId(m interface{}) ([]string, string) {
	f, _, r := GetFieldsAndSortAndRefId(m)
	return f, r
}
func GetSortAndRefId(m interface{}) (string, string) {
	_, s, r := GetFieldsAndSortAndRefId(m)
	return s, r
}
func GetFields(m interface{}) []string {
	f, _, _ := GetFieldsAndSortAndRefId(m)
	return f
}
func GetSort(m interface{}) string {
	_, s, _ := GetFieldsAndSortAndRefId(m)
	return s
}
func GetRefId(m interface{}) string {
	_, _, r := GetFieldsAndSortAndRefId(m)
	return r
}
func GetFieldsAndSortAndRefId(m interface{}) ([]string, string, string) {
	var fields []string
	var sort, refId string
	if sModel, ok := m.(*Filter); ok {
		return sModel.Fields, sModel.Sort, sModel.RefId
	} else {
		value := reflect.Indirect(reflect.ValueOf(m))
		numField := value.NumField()
		for i := 0; i < numField; i++ {
			fn := value.Type().Field(i).Name
			if fn == "Sort" {
				if s, ok := value.Field(i).Interface().(string); ok {
					sort = s
				}
			} else if fn == "Fields" {
				if s, ok := value.Field(i).Interface().([]string); ok {
					fields = s
				}
			} else if fn == "RefId" {
				if s, ok := value.Field(i).Interface().(string); ok {
					refId = s
				}
			}
			if sModel1, ok := value.Field(i).Interface().(*Filter); ok {
				return sModel1.Fields, sModel1.Sort, sModel1.RefId
			}
		}
		return fields, sort, refId
	}

}
