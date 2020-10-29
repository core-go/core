package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type CodeModel struct {
	Code     string `mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Text     string `mapstructure:"text" json:"text,omitempty" gorm:"column:text" bson:"text,omitempty" dynamodbav:"text,omitempty" firestore:"text,omitempty"`
	Name     string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Value    string `mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Sequence int32  `mapstructure:"sequence" json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence,omitempty" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
}
type CodeConfig struct {
	Master   string `mapstructure:"master" json:"master,omitempty" gorm:"column:master" bson:"master,omitempty" dynamodbav:"master,omitempty" firestore:"master,omitempty"`
	Code     string `mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Text     string `mapstructure:"text" json:"text,omitempty" gorm:"column:text" bson:"text,omitempty" dynamodbav:"text,omitempty" firestore:"text,omitempty"`
	Name     string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Value    string `mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Sequence string `mapstructure:"sequence" json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence,omitempty" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
	Active   string `mapstructure:"active" json:"active,omitempty" gorm:"column:active" bson:"active,omitempty" dynamodbav:"active,omitempty" firestore:"active,omitempty"`
}
type CodeLoader interface {
	Load(ctx context.Context, master string) ([]CodeModel, error)
}
type SqlCodeLoader struct {
	DB            *sql.DB
	Table         string
	Config        CodeConfig
	QuestionParam bool
}

func NewSqlCodeLoader(db *sql.DB, table string, config CodeConfig, questionParam bool) *SqlCodeLoader {
	return &SqlCodeLoader{DB: db, Table: table, Config: config, QuestionParam: questionParam}
}
func (l SqlCodeLoader) Load(ctx context.Context, master string) ([]CodeModel, error) {
	models := make([]CodeModel, 0)
	s := make([]string, 0)
	values := make([]interface{}, 0)
	sql2 := ""

	c := l.Config
	if len(c.Code) > 0 {
		sf := fmt.Sprintf("%s as code", c.Code)
		s = append(s, sf)
	}
	if len(c.Text) > 0 {
		sf := fmt.Sprintf("%s as text", c.Text)
		s = append(s, sf)
	}
	if len(c.Name) > 0 {
		sf := fmt.Sprintf("%s as name", c.Name)
		s = append(s, sf)
	}
	if len(c.Value) > 0 {
		sf := fmt.Sprintf("%s as value", c.Value)
		s = append(s, sf)
	}
	osequence := ""
	if len(c.Sequence) > 0 {
		sf := fmt.Sprintf("%s as sequence", c.Sequence)
		s = append(s, sf)
		osequence = fmt.Sprintf("order by %s", c.Sequence)
	}
	p1 := ""
	if l.QuestionParam {
		p1 = fmt.Sprintf("%s = ?", c.Master)
	} else {
		p1 = fmt.Sprintf("%s = $1", c.Master)
	}
	values = append(values, master)
	cols := strings.Join(s, ",")
	if len(c.Active) > 0 {
		p2 := ""
		if !l.QuestionParam {
			p2 = fmt.Sprintf("and %s = $2", c.Active)
		} else {
			p2 = fmt.Sprintf("and %s = ?", c.Active)
		}
		values = append(values, true)
		if cols == "" {
			cols = "*"
		}
		sql2 = fmt.Sprintf("select %s from %s where %s %s %s", cols, l.Table, p1, p2, osequence)
	} else {
		if cols == "" {
			cols = "*"
		}
		sql2 = fmt.Sprintf("select %s from %s where %s %s", cols, l.Table, p1, osequence)
	}
	if len(sql2) > 0 {
		rows, err1 := l.DB.Query(sql2, values...)
		if err1 != nil {
			return nil, err1
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		// get list indexes column
		modelTypes := reflect.TypeOf(models).Elem()
		modelType := reflect.TypeOf(CodeModel{})
		indexes, _ := getColumnIndexes(modelType, columns)
		tb, err2 := ScanType(rows, modelTypes, indexes)
		if err2 != nil {
			return nil, err2
		}
		for _, v := range tb {
			if c, ok := v.(*CodeModel); ok {
				models = append(models, *c)
			}
		}
	}
	return models, nil
}

// StructScan : transfer struct to slice for scan
func StructScan(s interface{}, indexColumns []int) (r []interface{}) {
	if s != nil {
		maps := reflect.Indirect(reflect.ValueOf(s))
		for _, index := range indexColumns {
			r = append(r, maps.Field(index).Addr().Interface())
		}
	}
	return
}

func getColumnIndexes(modelType reflect.Type, columnsName []string) (indexes []int, err error) {
	if modelType.Kind() != reflect.Struct {
		return nil, errors.New("Bad Type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := findTag(ormTag, "column")
		if ok {
			if contains(columnsName, column) {
				indexes = append(indexes, i)
			}
		}
	}
	return
}

func findTag(tag string, key string) (string, bool) {
	if has := strings.Contains(tag, key); has {
		str1 := strings.Split(tag, ";")
		num := len(str1)
		for i := 0; i < num; i++ {
			str2 := strings.Split(str1[i], ":")
			for j := 0; j < len(str2); j++ {
				if str2[j] == key {
					return str2[j+1], true
				}
			}
		}
	}
	return "", false
}

func contains(array []string, v string) bool {
	for _, s := range array {
		if s == v {
			return true
		}
	}
	return false
}

func ScanType(rows *sql.Rows, modelTypes reflect.Type, indexes []int) (t []interface{}, err error) {
	for rows.Next() {
		initArray := reflect.New(modelTypes).Interface()
		if err = rows.Scan(StructScan(initArray, indexes)...); err == nil {
			t = append(t, initArray)
		}
	}
	return
}
