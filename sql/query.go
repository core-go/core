package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	DefaultPagingFormat = " limit %s offset %s "
	OraclePagingFormat  = " offset %s rows fetch next %s rows only "
	desc                = "desc"
	asc                 = "asc"
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
func BuildFromQuery(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, models interface{}, query string, params []interface{}, limit int64, offset int64, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options...func(context.Context, interface{}) (interface{}, error)) (int64, error) {
	var mp func(context.Context, interface{}) (interface{}, error)
	if len(options) > 0 && options[0] != nil {
		mp = options[0]
	}
	var total int64
	driver := GetDriver(db)
	if limit <= 0 {
		er1 := QueryWithArray(ctx, db, fieldsIndex, models, toArray, query, params...)
		if er1 != nil {
			return -1, er1
		}
		objectValues := reflect.Indirect(reflect.ValueOf(models))
		if objectValues.Kind() == reflect.Slice {
			i := objectValues.Len()
			total = int64(i)
		}
		er2 := BuildSearchResult(ctx, models, mp)
		return total, er2
	} else {
		if driver == DriverOracle {
			queryPaging := BuildPagingQueryByDriver(query, limit, offset, driver)
			er1 := QueryAndCount(ctx, db, fieldsIndex, models, toArray, &total, queryPaging, params...)
			if er1 != nil {
				return -1, er1
			}
			er2 := BuildSearchResult(ctx, models, mp)
			return total, er2
		} else {
			queryPaging := BuildPagingQuery(query, limit, offset, driver)
			queryCount := BuildCountQuery(query)
			er1 := QueryWithArray(ctx, db, fieldsIndex, models, toArray, queryPaging, params...)
			if er1 != nil {
				return -1, er1
			}
			total, er2 := Count(ctx, db, queryCount, params...)
			if er2 != nil {
				total = 0
			}
			er3 := BuildSearchResult(ctx, models, mp)
			return total, er3
		}
	}
}
func BuildPagingQueryByDriver(sql string, limit int64, offset int64, driver string) string {
	s2 := BuildPagingQuery(sql, limit, offset, driver)
	if driver != DriverOracle {
		return s2
	} else {
		l := len(" distinct ")
		i := strings.Index(sql, " distinct ")
		if i < 0 {
			i = strings.Index(sql, " DISTINCT ")
		}
		if i < 0 {
			l = len("select") + 1
			i = strings.Index(s2, "select")
		}
		if i < 0 {
			i = strings.Index(s2, "SELECT")
		}
		if i >= 0 {
			return s2[0:l] + " count(*) over() as total, " + s2[l:]
		}
		return s2
	}
}
func BuildPagingQuery(sql string, limit int64, offset int64, opts...string) string {
	if offset < 0 {
		offset = 0
	}
	if limit > 0 {
		var pagingQuery string
		if len(opts) > 0 && opts[0] == DriverOracle {
			pagingQuery = fmt.Sprintf(OraclePagingFormat, strconv.FormatInt(offset, 10), strconv.FormatInt(limit, 10))
		} else {
			pagingQuery = fmt.Sprintf(DefaultPagingFormat, strconv.FormatInt(limit, 10), strconv.FormatInt(offset, 10))
		}
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

func BuildSearchResult(ctx context.Context, models interface{}, mp func(context.Context, interface{}) (interface{}, error)) error {
	if mp == nil {
		return nil
	}
	_, err := MapModels(ctx, models, mp)
	return err
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
func BuildSort(sortString string, modelType reflect.Type) string {
	var sort = GetSort(sortString, modelType)
	if len(sort) > 0 {
		return ` order by ` + sort
	} else {
		return ""
	}
}
func ExtractArray(values []interface{}, field interface{}) []interface{} {
	s := reflect.Indirect(reflect.ValueOf(field))
	for i := 0; i < s.Len(); i++ {
		values = append(values, s.Index(i).Interface())
	}
	return values
}
