package search

import (
	"reflect"
	"strconv"
	"strings"
)

func ToCsv(fields []string, r interface{}, total int64, embedField string, opts ...map[string]int) (out string) {
	val := reflect.ValueOf(r)
	models := reflect.Indirect(val)

	if models.Len() == 0 {
		return "0"
	}
	var rows []string
	rows = append(rows, strconv.FormatInt(total, 10))
	rows = BuildCsv(rows, fields, models, embedField, opts...)
	return strings.Join(rows, "\n")
}
func ToNextCsv(fields []string, r interface{}, nextPageToken string, embedField string, opts ...map[string]int) (out string) {
	val := reflect.ValueOf(r)
	models := reflect.Indirect(val)

	if models.Len() == 0 {
		return "0"
	}
	var rows []string
	rows = append(rows, nextPageToken)
	rows = BuildCsv(rows, fields, models, embedField, opts...)
	return strings.Join(rows, "\n")
}
