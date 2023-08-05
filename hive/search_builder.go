package hive

import (
	"context"
	"fmt"
	hv "github.com/beltran/gohive"
	"reflect"
	"strconv"
	"strings"
)

const (
	desc                = "desc"
	asc                 = "asc"
	DefaultPagingFormat = " limit %s offset %s "
)

func GetOffset(limit int64, page int64, opts...int64) int64 {
	var firstLimit int64 = 0
	if len(opts) > 0 && opts[0] > 0 {
		firstLimit = opts[0]
	}
	if firstLimit > 0 {
		if page <= 1 {
			return 0
		} else {
			offset := limit*(page-2) + firstLimit
			if offset < 0 {
				return 0
			}
			return offset
		}
	} else {
		offset := limit * (page - 1)
		if offset < 0 {
			return 0
		}
		return offset
	}
}

type SearchBuilder struct {
	Connection  *hv.Connection
	BuildQuery  func(sm interface{}) string
	ModelType   reflect.Type
	Map         func(ctx context.Context, model interface{}) (interface{}, error)
	fieldsIndex map[string]int
}
func NewSearchBuilder(connection *hv.Connection, modelType reflect.Type, buildQuery func(interface{}) string, options ...func(context.Context, interface{}) (interface{}, error)) (*SearchBuilder, error) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	fieldsIndex, err := GetColumnIndexes(modelType)
	if err != nil {
		return nil, err
	}
	builder := &SearchBuilder{Connection: connection, fieldsIndex: fieldsIndex, BuildQuery: buildQuery, ModelType: modelType, Map: mp}
	return builder, nil
}

func (b *SearchBuilder) Search(ctx context.Context, m interface{}, results interface{}, limit int64, offset int64) (int64, error) {
	sql := b.BuildQuery(m)
	query := BuildPagingQuery(sql, limit, offset)
	cursor := b.Connection.Cursor()
	defer cursor.Close()
	cursor.Exec(ctx, sql)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	err := Query(ctx, cursor, b.fieldsIndex, results, query)
	if err != nil {
		return -1, err
	}
	countQuery := BuildCountQuery(sql)
	cursor.Exec(ctx, countQuery)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	var count int64
	for cursor.HasMore(ctx) {
		cursor.FetchOne(ctx, &count)
		if cursor.Err != nil {
			return count, cursor.Err
		}
	}
	if b.Map != nil {
		_, err := MapModels(ctx, results, b.Map)
		return count, err
	}
	return count, err
}
func Count(ctx context.Context, cursor *hv.Cursor, query string) (int64, error) {
	var count int64
	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return -1, cursor.Err
	}
	for cursor.HasMore(ctx) {
		cursor.FetchOne(ctx, &count)
		if cursor.Err != nil {
			return count, cursor.Err
		}
	}
	return 0, nil
}
func BuildPagingQuery(sql string, limit int64, offset int64) string {
	if offset < 0 {
		offset = 0
	}
	if limit > 0 {
		pagingQuery := fmt.Sprintf(DefaultPagingFormat, strconv.FormatInt(limit, 10), strconv.FormatInt(offset, 10))
		sql += pagingQuery
	}
	return sql
}
func BuildCountQuery(sql string) string {
	i := strings.Index(sql, "select ")
	if i < 0 {
		return sql
	}
	j := strings.Index(sql, " from ")
	if j < 0 {
		return sql
	}
	k := strings.Index(sql, " order by ")
	h := strings.Index(sql, " distinct ")
	if h > 0 {
		sql3 := `select count(*) as total from (` + sql[i:] + `) as main`
		return sql3
	}
	if k > 0 {
		sql3 := `select count(*) as total ` + sql[j:k]
		return sql3
	} else {
		sql3 := `select count(*) as total ` + sql[j:]
		return sql3
	}
}
func GetSort(sortString string, modelType reflect.Type) string {
	var sort = make([]string, 0)
	sorts := strings.Split(sortString, ",")
	for i := 0; i < len(sorts); i++ {
		sortField := strings.TrimSpace(sorts[i])
		fieldName := sortField
		c := sortField[0:1]
		if c == "-" || c == "+" {
			fieldName = sortField[1:]
		}
		columnName := GetColumnNameForSearch(modelType, fieldName)
		if len(columnName) > 0 {
			sortType := GetSortType(c)
			sort = append(sort, columnName+" "+sortType)
		}
	}
	if len(sort) > 0 {
		return strings.Join(sort, ",")
	} else {
		return ""
	}
}
func BuildSort(sortString string, modelType reflect.Type) string {
	sort := GetSort(sortString, modelType)
	if len(sort) > 0 {
		return ` order by ` + sort
	} else {
		return ""
	}
}
func GetColumnNameForSearch(modelType reflect.Type, sortField string) string {
	sortField = strings.TrimSpace(sortField)
	i, _, column := GetFieldByJson(modelType, sortField)
	if i > -1 {
		return column
	}
	return ""
}
func GetSortType(sortType string) string {
	if sortType == "-" {
		return desc
	} else {
		return asc
	}
}
func GetFieldByJson(modelType reflect.Type, jsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 && strings.Split(tag1, ",")[0] == jsonName {
			if tag2, ok2 := field.Tag.Lookup("gorm"); ok2 {
				if has := strings.Contains(tag2, "column"); has {
					str1 := strings.Split(tag2, ";")
					num := len(str1)
					for k := 0; k < num; k++ {
						str2 := strings.Split(str1[k], ":")
						for j := 0; j < len(str2); j++ {
							if str2[j] == "column" {
								return i, field.Name, str2[j+1]
							}
						}
					}
				}
			}
			return i, field.Name, ""
		}
	}
	return -1, jsonName, jsonName
}
