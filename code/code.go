package code

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	driverPostgres   = "postgres"
	driverMysql      = "mysql"
	driverMssql      = "mssql"
	driverOracle     = "oracle"
	driverSqlite3    = "sqlite3"
	driverNotSupport = "no support"
)

type Model struct {
	Id       string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Code     string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Value    string `yaml:"value" mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Name     string `yaml:"name"" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Text     string `yaml:"text" mapstructure:"text" json:"text,omitempty" gorm:"column:text" bson:"text,omitempty" dynamodbav:"text,omitempty" firestore:"text,omitempty"`
	Sequence int32  `yaml:"sequence mapstructure:"sequence" json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence,omitempty" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
}
type StructureConfig struct {
	Master   string      `yaml:"master" mapstructure:"master" json:"master,omitempty" gorm:"column:master" bson:"master,omitempty" dynamodbav:"master,omitempty" firestore:"master,omitempty"`
	Id       string      `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Code     string      `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Text     string      `yaml:"text" mapstructure:"text" json:"text,omitempty" gorm:"column:text" bson:"text,omitempty" dynamodbav:"text,omitempty" firestore:"text,omitempty"`
	Name     string      `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Value    string      `yaml:"value" mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Sequence string      `yaml:"sequence" mapstructure:"sequence" json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence,omitempty" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
	Status   string      `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Active   interface{} `yaml:"active" mapstructure:"active" json:"active,omitempty" gorm:"column:active" bson:"active,omitempty" dynamodbav:"active,omitempty" firestore:"active,omitempty"`
}
type Loader interface {
	Load(ctx context.Context, master string) ([]Model, error)
}
type SqlLoader struct {
	DB     *sql.DB
	Table  string
	Config StructureConfig
	Build  func(i int) string
	Map    func(col string) string
}
type DynamicSqlLoader struct {
	DB             *sql.DB
	Query          string
	ParameterCount int
	Map            func(col string) string
	driver         string
}

func NewDefaultDynamicSqlCodeLoader(db *sql.DB, query string, options ...int) *DynamicSqlLoader {
	var parameterCount int
	if len(options) > 0 {
		parameterCount = options[0]
	} else {
		parameterCount = 1
	}
	return NewDynamicSqlCodeLoader(db, query, parameterCount, true)
}
func NewDynamicSqlCodeLoader(db *sql.DB, query string, parameterCount int, options ...bool) *DynamicSqlLoader {
	driver := getDriver(db)
	var mp func(string) string
	if driver == driverOracle {
		mp = strings.ToUpper
	} else {
		mp = strings.ToLower
	}
	if parameterCount < 0 {
		parameterCount = 1
	}
	var handleDriver bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	if handleDriver {
		if driver == driverOracle || driver == driverPostgres || driver == driverMssql {
			var x string
			if driver == driverOracle {
				x = ":val"
			} else if driver == driverPostgres {
				x = "$"
			} else if driver == driverMssql {
				x = "@p"
			}
			for i := 0; i < parameterCount; i++ {
				count := i + 1
				query = strings.Replace(query, "?", x+strconv.Itoa(count), 1)
			}
		}
	}
	return &DynamicSqlLoader{DB: db, Query: query, ParameterCount: parameterCount, Map: mp}
}
func (l DynamicSqlLoader) Load(ctx context.Context, master string) ([]Model, error) {
	models := make([]Model, 0)

	var rows *sql.Rows
	var er1 error
	if l.ParameterCount > 0 {
		params := make([]interface{}, 0)
		for i := 1; i <= l.ParameterCount; i++ {
			params = append(params, master)
		}
		rows, er1 = l.DB.QueryContext(ctx, l.Query, params...)
	} else {
		rows, er1 = l.DB.QueryContext(ctx, l.Query)
	}

	if er1 != nil {
		return models, er1
	}
	defer rows.Close()
	columns, er2 := rows.Columns()
	if er2 != nil {
		return models, er2
	}
	// get list indexes column
	modelType := reflect.TypeOf(Model{})

	fieldsIndexSelected := make([]int, 0)
	fieldsIndex, er3 := getColumnIndexes(modelType, l.Map)
	if er3 != nil {
		return models, er3
	}
	for _, columnsName := range columns {
		if index, ok := fieldsIndex[columnsName]; ok {
			fieldsIndexSelected = append(fieldsIndexSelected, index)
		}
	}
	tb, er4 := scanType(rows, modelType, fieldsIndexSelected)
	if er4 != nil {
		return models, er4
	}
	for _, v := range tb {
		if c, ok := v.(*Model); ok {
			models = append(models, *c)
		}
	}
	return models, nil
}
func NewSqlCodeLoader(db *sql.DB, table string, config StructureConfig, options ...func(i int) string) *SqlLoader {
	var build func(i int) string
	if len(options) > 0 && options[0] != nil {
		build = options[0]
	} else {
		build = getBuild(db)
	}
	driver := getDriver(db)
	var mp func(string) string
	if driver == driverOracle {
		mp = strings.ToUpper
	}
	return &SqlLoader{DB: db, Table: table, Config: config, Build: build, Map: mp}
}
func (l SqlLoader) Load(ctx context.Context, master string) ([]Model, error) {
	models := make([]Model, 0)
	s := make([]string, 0)
	values := make([]interface{}, 0)
	sql2 := ""

	c := l.Config
	if len(c.Id) > 0 {
		sf := fmt.Sprintf("%s as id", c.Id)
		s = append(s, sf)
	}
	if len(c.Code) > 0 {
		sf := fmt.Sprintf("%s as code", c.Code)
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
	if len(c.Text) > 0 {
		sf := fmt.Sprintf("%s as text", c.Text)
		s = append(s, sf)
	}
	osequence := ""
	if len(c.Sequence) > 0 {
		osequence = fmt.Sprintf("order by %s", c.Sequence)
	}
	p1 := ""
	i := 1
	if len(c.Master) > 0 {
		p1 = fmt.Sprintf("%s = %s", c.Master, l.Build(i))
		i = i + 1
		values = append(values, master)
	}
	cols := strings.Join(s, ",")
	if len(c.Status) > 0 && c.Active != nil {
		p2 := fmt.Sprintf("%s = %s", c.Status, l.Build(i))
		values = append(values, c.Active)
		if cols == "" {
			cols = "*"
		}
		if len(p1) > 0 {
			sql2 = fmt.Sprintf("select %s from %s where %s and %s %s", cols, l.Table, p1, p2, osequence)
		} else {
			sql2 = fmt.Sprintf("select %s from %s where %s %s", cols, l.Table, p2, osequence)
		}
	} else {
		if cols == "" {
			cols = "*"
		}
		if len(p1) > 0 {
			sql2 = fmt.Sprintf("select %s from %s where %s %s", cols, l.Table, p1, osequence)
		} else {
			sql2 = fmt.Sprintf("select %s from %s %s", cols, l.Table, osequence)
		}
	}
	if len(sql2) > 0 {
		rows, err1 := l.DB.QueryContext(ctx, sql2, values...)
		if err1 != nil {
			return nil, err1
		}
		defer rows.Close()
		columns, er1 := rows.Columns()
		if er1 != nil {
			return nil, er1
		}
		fieldsIndexSelected := make([]int, 0)
		modelType := reflect.TypeOf(Model{})
		// get list indexes column
		fieldsIndex, er3 := getColumnIndexes(modelType, l.Map)
		if er3 != nil {
			return models, er3
		}
		for _, columnsName := range columns {
			if index, ok := fieldsIndex[columnsName]; ok {
				fieldsIndexSelected = append(fieldsIndexSelected, index)
			}
		}
		tb, er3 := scanType(rows, modelType, fieldsIndexSelected)
		if er3 != nil {
			return nil, er3
		}
		for _, v := range tb {
			if c, ok := v.(*Model); ok {
				models = append(models, *c)
			}
		}
	}
	return models, nil
}

func scanType(rows *sql.Rows, modelType reflect.Type, indexes []int) (t []interface{}, err error) {
	for rows.Next() {
		initModel := reflect.New(modelType).Interface()
		if err = rows.Scan(structScan(initModel, indexes)...); err == nil {
			t = append(t, initModel)
		}
	}
	return
}
func structScan(s interface{}, indexColumns []int) (r []interface{}) {
	if s != nil {
		maps := reflect.Indirect(reflect.ValueOf(s))
		for _, index := range indexColumns {
			r = append(r, maps.Field(index).Addr().Interface())
		}
	}
	return
}
func getColumnIndexes(modelType reflect.Type, mp func(col string) string) (map[string]int, error) {
	mapp := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return mapp, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := findTag(ormTag, "column")
		if ok {
			if mp != nil {
				column = mp(column)
			}
			mapp[column] = i
		}
	}
	return mapp, nil
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

func buildParam(i int) string {
	return "?"
}
func buildOracleParam(i int) string {
	return ":val" + strconv.Itoa(i)
}
func buildMsSqlParam(i int) string {
	return "@p" + strconv.Itoa(i)
}
func buildDollarParam(i int) string {
	return "$" + strconv.Itoa(i)
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
func getDriver(db *sql.DB) string {
	if db == nil {
		return driverNotSupport
	}
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return driverPostgres
	case "*godror.drv":
		return driverOracle
	case "*mysql.MySQLDriver":
		return driverMysql
	case "*mssql.Driver":
		return driverMssql
	case "*sqlite3.SQLiteDriver":
		return driverSqlite3
	default:
		return driverNotSupport
	}
}
