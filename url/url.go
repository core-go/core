package text

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type URL struct {
	Id   string  `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name *string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Url  *string `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
}

type URLPort interface {
	Load(ctx context.Context, id string) (*URL, error)
	Query(ctx context.Context, ids []string) ([]URL, error)
}

type URLAdapter struct {
	db         *sql.DB
	Select     string
	BuildParam func(int) string
}

func NewURLAdapter(db *sql.DB, query string, opts ...func(i int) string) (*URLAdapter, error) {
	var buildParam func(i int) string
	if len(opts) > 0 && opts[0] != nil {
		buildParam = opts[0]
	} else {
		buildParam = getBuild(db)
	}
	return &URLAdapter{
		db:         db,
		Select:     query,
		BuildParam: buildParam,
	}, nil
}

func (r *URLAdapter) Load(ctx context.Context, id string) (*URL, error) {
	p := make([]string, 0)
	p = append(p, id)
	values, err := r.Query(ctx, p)
	if err != nil {
		return nil, err
	}
	if len(values) > 0 {
		return &values[0], nil
	}
	return nil, nil
}

func (r *URLAdapter) Query(ctx context.Context, ids []string) ([]URL, error) {
	var values []URL
	if len(ids) == 0 {
		return values, nil
	}
	le := len(ids)
	p := make([]interface{}, 0)
	for _, str := range ids {
		p = append(p, str)
	}
	var arrValue []string
	for i := 1; i <= le; i++ {
		param := r.BuildParam(i)
		arrValue = append(arrValue, param)
	}
	query := r.Select + fmt.Sprintf(" in (%s)", strings.Join(arrValue, ","))
	rows, err := r.db.QueryContext(ctx, query, p...)
	defer rows.Close()

	for rows.Next() {
		var row URL
		if err := rows.Scan(&row.Id, &row.Name, &row.Url); err != nil {
			return values, err
		}
		values = append(values, row)
	}
	if err = rows.Err(); err != nil {
		return values, err
	}
	SortById(values)
	return values, nil
}

func ToMap(rows []URL) map[string]*URL {
	rs := make(map[string]*URL, 0)
	for _, row := range rows {
		rs[row.Id] = &row
	}
	return rs
}

func Unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
func SortById(urls []URL) {
	sort.Slice(urls, func(i, j int) bool { return urls[i].Id < urls[j].Id })
}
func BinarySearch(id string, a []URL) (result int, searchCount int) {
	mid := len(a) / 2
	x := strings.Compare(a[mid].Id, id)
	switch {
	case len(a) == 0:
		result = -1 // not found
	case x > 0:
		result, searchCount = BinarySearch(id, a[:mid])
	case x < 0:
		result, searchCount = BinarySearch(id, a[mid+1:])
		if result >= 0 { // if anything but the -1 "not found" result
			result += mid + 1
		}
	default: // a[mid] == id
		result = mid // found
	}
	searchCount++
	return
}
func getBuild(db *sql.DB) func(i int) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return buildDollarParam
	case "*godror.drv":
		return buildOracleParam
	case "*mssql.Driver":
		return buildMsSqlParam
	default:
		return buildParam
	}
}
func buildParam(i int) string {
	return "?"
}
func buildOracleParam(i int) string {
	return ":" + strconv.Itoa(i)
}
func buildMsSqlParam(i int) string {
	return "@p" + strconv.Itoa(i)
}
func buildDollarParam(i int) string {
	return "$" + strconv.Itoa(i)
}
