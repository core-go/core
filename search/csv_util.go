package search

import (
	"reflect"
	"strconv"
	"strings"
)

func ToCsv(fields []string, r interface{}, total int64, nextPageToken string, embedField string) (out string) {
	val := reflect.ValueOf(r)
	models := reflect.Indirect(val)

	if models.Len() == 0 {
		return "0"
	}
	var rows []string
	rows = append(rows, strconv.FormatInt(total, 10)+","+nextPageToken)
	rows = BuildCsv(rows, fields, models, embedField)
	return strings.Join(rows, "\n")
	return out
}

func IsLastPage(models interface{}, count int64, pageIndex int64, pageSize int64, initPageSize int64) bool {
	lengthModels := int64(reflect.Indirect(reflect.ValueOf(models)).Len())
	var receivedItems int64

	if initPageSize > 0 {
		if pageIndex == 1 {
			receivedItems = initPageSize
		} else if pageIndex > 1 {
			receivedItems = pageSize*(pageIndex-2) + initPageSize + lengthModels
		}
	} else {
		receivedItems = pageSize*(pageIndex-1) + lengthModels
	}
	return receivedItems >= count
}
