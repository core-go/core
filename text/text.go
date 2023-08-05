package text

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Text struct {
	Id   string  `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Text *string `yaml:"text" mapstructure:"text" json:"text,omitempty" gorm:"column:text" bson:"text,omitempty" dynamodbav:"text,omitempty" firestore:"text,omitempty"`
}

type TextAdapter struct {
	db         *sql.DB
	Table      string
	Id         string
	Text       string
	BuildParam func(int) string
}

func NewTextAdapter(db *sql.DB, table string, id string, text string, opts...func(i int) string) (*TextAdapter, error) {
	var buildParam func(i int) string
	if len(opts) > 0 && opts[0] != nil {
		buildParam = opts[0]
	} else {
		buildParam = getBuild(db)
	}
	return &TextAdapter{
		db:         db,
		Table:      table,
		Id:         id,
		Text:       text,
		BuildParam: buildParam,
	}, nil
}

func (r *TextAdapter) Load(ctx context.Context, ids []string) ([]Text, error) {
	var values []Text
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
	query := fmt.Sprintf("select %s, %s from %s where %s in (%s)", r.Id, r.Text, r.Table, r.Id, strings.Join(arrValue, ","))
	rows, err := r.db.QueryContext(ctx, query, p...)
	defer rows.Close()

	for rows.Next() {
		var row Text
		if err := rows.Scan(&row.Id, &row.Text); err != nil {
			return values, err
		}
		values = append(values, row)
	}
	if err = rows.Err(); err != nil {
		return values, err
	}
	return values, nil
}
func ToMap(rows []Text) map[string]*Text {
	rs := make(map[string]*Text, 0)
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
