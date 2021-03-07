package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	DriverPostgres   = "postgres"
	DriverMysql      = "mysql"
	DriverMssql      = "mssql"
	DriverOracle     = "oracle"
	DriverNotSupport = "no support"
	FormatDate     = "2006-01-02 15:04:05"
)

type HistoryWriter interface {
	Write(ctx context.Context, id interface{}, diff DiffModel, approvedBy string) error
}
type KeyBuilder interface {
	BuildKey(object interface{}) string
	BuildKeyFromMap(keyMap map[string]interface{}, idNames []string) string
}
type DiffConfig struct {
	HistoryId  string `mapstructure:"history_id" json:"historyId,omitempty" gorm:"column:historyid" bson:"_historyId,omitempty" dynamodbav:"historyId,omitempty" firestore:"historyId,omitempty"`
	Id         string `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Origin     string `mapstructure:"origin" json:"origin,omitempty" gorm:"column:origin" bson:"origin,omitempty" dynamodbav:"origin,omitempty" firestore:"origin,omitempty"`
	Value      string `mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	ChangedBy  string `mapstructure:"changedBy" json:"changedBy,omitempty" gorm:"column:changedBy" bson:"changedBy,omitempty" dynamodbav:"changedBy,omitempty" firestore:"changedBy,omitempty"`
	ApprovedBy string `mapstructure:"approvedBy" json:"approvedBy,omitempty" gorm:"column:approvedBy" bson:"approvedBy,omitempty" dynamodbav:"approvedBy,omitempty" firestore:"approvedBy,omitempty"`
	Timestamp  string `mapstructure:"timestamp" json:"timestamp,omitempty" gorm:"column:timestamp" bson:"timestamp,omitempty" dynamodbav:"timestamp,omitempty" firestore:"timestamp,omitempty"`
}
type SqlDiffReader struct {
	DB           *sql.DB
	Table        string
	Entity       string
	EntityType   string
	IdNames      []string
	Config       DiffConfig
	KeyBuilder   KeyBuilder
	Driver       string
	columnSelect string
}

type SqlDiffListReader struct {
	DB           *sql.DB
	Table        string
	Entity       string
	EntityType   string
	IdNames      []string
	Config       DiffConfig
	KeyBuilder   KeyBuilder
	Driver       string
	columnSelect string
}
type SqlHistoryWriter struct {
	// DB         *sql.DB
	Table      string
	Entity     string
	IdNames    []string
	Config     DiffConfig
	KeyBuilder KeyBuilder
	Generate   func()(string, error)
}

func NewSqlDiffReader(DB *sql.DB, table string, entity string, entityType string, idNames []string, config DiffConfig, keyBuilder KeyBuilder) *SqlDiffReader {
	columnSelect := BuildQueryColumn(config)
	driver := GetDriver(DB)
	return &SqlDiffReader{DB, table, entity, entityType, idNames, getDefaultConfig(config), keyBuilder, driver, columnSelect}
}

func NewSqlDiffListReader(DB *sql.DB, table string, tableEntity string, entityType string, idNames []string, config DiffConfig, keyBuilder KeyBuilder) *SqlDiffListReader {
	columnSelect := BuildQueryColumn(config)
	driver := GetDriver(DB)
	return &SqlDiffListReader{DB, table, tableEntity, entityType, idNames, getDefaultConfig(config), keyBuilder, driver, columnSelect}
}

func NewSqlHistoryWriter(table string, entity string, idNames []string, config DiffConfig, keyBuilder KeyBuilder, generate func()(string, error)) *SqlHistoryWriter {
	return &SqlHistoryWriter{table, entity, idNames, getDefaultConfig(config), keyBuilder, generate}
}

func getDefaultConfig(config DiffConfig) DiffConfig {
	if config.Id == "" {
		config.Id = "id"
	}
	if config.Origin == "" {
		config.Origin = "origin"
	}
	if config.Value == "" {
		config.Value = "value"
	}
	return config
}

func (r SqlHistoryWriter) Write(ctx context.Context, db *sql.DB, id interface{}, diff DiffModel, approvedBy string) error {
	entityID := ""
	dt := time.Now()
	updateTime := dt.Format(FormatDate)

	if len(r.IdNames) == 1 {
		entityID = r.IdNames[0]
	} else {
		if v, ok := id.(string); ok {
			entityID = r.KeyBuilder.BuildKey(v)
		}
		if v, ok := id.(int); ok {
			entityID = r.KeyBuilder.BuildKey(v)
		}
		if v, ok := id.(map[string]interface{}); ok {
			entityID = r.KeyBuilder.BuildKeyFromMap(v, r.IdNames)
		}
	}
	var sqlVar []interface{}
	strSQL := ""
	sqlParam := ""
	strSQL += "entitytablename, "
	sqlParam += "?,"
	sqlVar = append(sqlVar, r.Entity)
	if len(r.Config.ApprovedBy) > 1 {
		strSQL += r.Config.ApprovedBy + `, `
		sqlVar = append(sqlVar, approvedBy)
		sqlParam += "?,"
	}
	if len(r.Config.HistoryId) > 1 {
		if r.Generate != nil {
			historyID, err := r.Generate()
			if err != nil {
				return err
			}
			strSQL += r.Config.HistoryId + `, `
			sqlVar = append(sqlVar, historyID)
			sqlParam += "?,"
		}
	}
	if len(r.Config.Timestamp) > 1 {
		strSQL += r.Config.Timestamp + `, `
		sqlVar = append(sqlVar, updateTime)
		sqlParam += "?,"
	}
	if len(r.Config.Value) > 1 {
		strSQL += r.Config.Value + `, `
		str := fmt.Sprintf("%v", diff.Value)
		sqlVar = append(sqlVar, str)
		sqlParam += "?,"
	}
	if len(r.Config.Origin) > 1 {
		strSQL += r.Config.Origin + `, `
		str := fmt.Sprintf("%v", diff.Origin)
		sqlVar = append(sqlVar, str)
		sqlParam += "?,"
	}
	if len(r.Config.Id) > 1 {
		strSQL += r.Config.Id + `, `
		sqlVar = append(sqlVar, entityID)
		sqlParam += "?,"
	}
	if len(r.Config.ChangedBy) > 1 {
		strSQL += r.Config.ChangedBy + `, `
		sqlVar = append(sqlVar, diff.By)
		sqlParam += "?,"
	}
	strSQL = strings.TrimRight(strSQL, ", ")
	sqlParam = strings.TrimRight(sqlParam, ", ")
	query := `insert into ` + r.Table + `(` + strSQL + `) 
		values (` + sqlParam + `)`
	_, err := db.Exec(query, sqlVar...)
	if err != nil {
		return err
	}
	return nil
}

func (r SqlDiffReader) Diff(ctx context.Context, id interface{}) (*DiffModel, error) {
	i, err := r.GetEntityById(id, r.IdNames)
	if err != nil {
		return nil, err
	}
	if result, ok := i.(*DiffModel); ok {
		return result, nil
	}
	return nil, nil
}

func (c SqlDiffListReader) Diff(ctx context.Context, ids interface{}) (*[]DiffModel, error) {
	i, err := c.getEntityByIds(c.KeyBuilder, ids, c.IdNames)
	if err != nil {
		return nil, err
	}
	if result, ok := i.(*[]DiffModel); ok {
		return result, nil
	}
	return nil, nil
}

func (r SqlDiffReader) GetEntityById(key interface{}, idNames []string) (interface{}, error) {
	var saveValueId interface{}
	if keyMap, ok := key.(map[string]interface{}); ok {
		entityId := r.KeyBuilder.BuildKeyFromMap(keyMap, idNames)
		if entityId == "" {
			return nil, errors.New("failed to build key")
		}
		key = entityId
		saveValueId = keyMap
	}
	result := DiffModel{}
	querySql := fmt.Sprintf("select %s from %s where %s = %s and %s = %s", r.columnSelect, r.Entity,
		r.Config.Id, BuildParam(1, r.Driver),
		r.EntityType, BuildParam(2, r.Driver))
	err := SqlQueryOne(r.DB, &result, querySql, key, r.Table)
	if err != nil {
		return nil, err
	}
	if saveValueId != nil && len(idNames) > 1 {
		result.Id = saveValueId
	}
	return &result, nil
}

func BuildParameters(numCol int, driver string) string {
	var arrValue []string
	for i := 0; i < numCol; i++ {
		arrValue = append(arrValue, BuildParam(i+1, driver))
	}
	return strings.Join(arrValue, ",")
}

func BuildParam(index int, driver string) string {
	switch driver {
	case DriverPostgres:
		return "$" + strconv.Itoa(index)
	case DriverOracle:
		return ":val" + strconv.Itoa(index)
	default:
		return "?"
	}
}

func BuildQueryColumn(config DiffConfig) string {
	sqlsel := make([]string, 0)
	colDiffModel := GetColumnNameDiffModel()
	if config.Id != "" {
		sqlsel = append(sqlsel, config.Id+" as "+colDiffModel[0])
	}
	if config.Origin != "" {
		sqlsel = append(sqlsel, config.Origin+" as "+colDiffModel[1])
	}
	if config.Value != "" {
		sqlsel = append(sqlsel, config.Value+" as "+colDiffModel[2])
	}
	if config.ApprovedBy != "" {
		sqlsel = append(sqlsel, config.ApprovedBy+" as "+colDiffModel[3])
	}
	return strings.Join(sqlsel, ",")
}

func GetColumnNameDiffModel() []string {
	ids := make([]string, 0)
	objectValue := reflect.Indirect(reflect.ValueOf(DiffModel{}))
	for i := 0; i < objectValue.NumField(); i++ {
		if colName, ok := GetColumnNameByIndex(objectValue.Type(), i); ok {
			ids = append(ids, colName)
		}
	}
	return ids
}

func (c SqlDiffListReader) getEntityByIds(keyBuilder KeyBuilder, keys interface{}, idNames []string) (interface{}, error) {
	arrayKeys := make([]interface{}, 0)
	args := make([]interface{}, 0)
	listIds := make(map[string]interface{}, 0)
	if keys != nil {
		keysInterface := reflect.Indirect(reflect.ValueOf(keys))
		n := keysInterface.Len()
		if len(idNames) > 1 {
			for i := 0; i < n; i++ {
				itemStruct := keysInterface.Index(i).Interface()
				entityId := keyBuilder.BuildKey(itemStruct)
				listIds[entityId] = itemStruct
				if entityId == "" {
					return nil, errors.New("failed to build key")
				}
				arrayKeys = append(arrayKeys, entityId)
			}
		} else {
			for i := 0; i < n; i++ {
				entityId := keysInterface.Index(i).Interface()
				if entityId == "" {
					return nil, errors.New("failed to build key")
				}
				arrayKeys = append(arrayKeys, entityId)
			}
		}
	} else {
		return nil, errors.New("failed keys nil")
	}
	n := len(arrayKeys)
	args = append(args, arrayKeys...)
	args = append(args, c.Table)
	results := make([]DiffModel, 0)
	querySql := fmt.Sprintf("select %s from %s where %s IN (%s) and %s = %s", c.columnSelect, c.Entity, c.Config.Id, BuildParameters(n, c.Driver), c.EntityType, BuildParam(n+1, c.Driver))
	err := SqlQuery(c.DB, &results, querySql, args...)
	// map object id
	for i, result := range results {
		id := result.Id.(*string)
		if idObject, ok := listIds[*id]; ok {
			results[i].Id = idObject
		}
	}
	if err != nil {
		return nil, err
	}
	return &results, nil
}

type DefaultKeyBuilder struct {
	PositionPrimaryKeysMap map[reflect.Type][]int
}

func NewDefaultKeyBuilder() *DefaultKeyBuilder {
	return &DefaultKeyBuilder{PositionPrimaryKeysMap: make(map[reflect.Type][]int)}
}

func (b *DefaultKeyBuilder) getPositionPrimaryKeys(modelType reflect.Type) []int {
	if b.PositionPrimaryKeysMap[modelType] == nil {
		var positions []int

		numField := modelType.NumField()
		for i := 0; i < numField; i++ {
			gorm := strings.Split(modelType.Field(i).Tag.Get("gorm"), ";")
			for _, value := range gorm {
				if value == "primary_key" {
					positions = append(positions, i)
					break
				}
			}
		}

		b.PositionPrimaryKeysMap[modelType] = positions
	}

	return b.PositionPrimaryKeysMap[modelType]
}

func SqlQueryOne(db *sql.DB, result *DiffModel, sql string, values ...interface{}) error {
	rows, err := db.Query(sql+" LIMIT 1", values...)
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	types, _ := rows.ColumnTypes()
	for rows.Next() {
		sizeCol := len(cols)
		vals := createValuesByType(types, sizeCol)
		err := rows.Scan(vals...)
		mapToModel(vals, result)
		return err
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		return rerr
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return err
	}
	return errors.New("NotFound Item.")
}

func mapToModel(vals []interface{}, result *DiffModel) {
	result.Id = vals[0]
	n := len(vals)
	origin, _ := convertStringToMap(vals[1].(*string))
	value, _ := convertStringToMap(vals[2].(*string))
	result.Origin = origin
	result.Value = value
	if n > 3 && vals[3] != nil {
		if v, ok := vals[3].(string); ok {
			result.By = v
		}
		if v, ok := vals[3].(*string); ok {
			result.By = *v
		}
	}
}

func createValuesByType(types []*sql.ColumnType, sizeCol int) []interface{} {
	vals := make([]interface{}, sizeCol)
	for i := range types {
		//TODO check add type if any type another
		vals[i] = new(string)
	}
	return vals
}

func convertStringToMap(str *string) (*map[string]interface{}, error) {
	reader := strings.NewReader(*str)
	var p map[string]interface{}
	err := json.NewDecoder(reader).Decode(&p)
	if err != nil {
		return &p, err
	}
	return &p, err
}

func SqlQuery(db *sql.DB, results *[]DiffModel, sql string, values ...interface{}) error {
	rows, err := db.Query(sql, values...)
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, err2 := rows.Columns()
	if err2 != nil {
		return err2
	}
	types, _ := rows.ColumnTypes()
	sizeCol := len(cols)
	for rows.Next() {
		result := DiffModel{}
		vals := createValuesByType(types, sizeCol)
		err := rows.Scan(vals...)
		if err != nil {
			return err
		}
		mapToModel(vals, &result)
		*results = append(*results, result)
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		return rerr
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func (b *DefaultKeyBuilder) BuildKey(object interface{}) string {
	ids := make(map[string]interface{})
	objectValue := reflect.Indirect(reflect.ValueOf(object))
	positions := b.getPositionPrimaryKeys(objectValue.Type())
	var values []string
	for _, position := range positions {
		if colName, ok := GetColumnNameByIndex(objectValue.Type(), position); ok {
			ids[colName] = fmt.Sprint(objectValue.Field(position).Interface())
			values = append(values, fmt.Sprint(objectValue.Field(position).Interface()))
		}
	}
	return strings.Join(values, "-")
}

func (b *DefaultKeyBuilder) BuildKeyFromMap(keyMap map[string]interface{}, idNames []string) string {
	var values []string
	for _, key := range idNames {
		if keyVal, exist := keyMap[key]; !exist {
			values = append(values, "")
		} else {
			str, ok := keyVal.(string)
			if !ok {
				return ""
			}
			values = append(values, str)
		}
	}
	return strings.Join(values, "-")
}

func GetColumnNameByIndex(ModelType reflect.Type, index int) (col string, colExist bool) {
	fields := ModelType.Field(index)
	tag, _ := fields.Tag.Lookup("gorm")

	if has := strings.Contains(tag, "column"); has {
		str1 := strings.Split(tag, ";")
		num := len(str1)
		for i := 0; i < num; i++ {
			str2 := strings.Split(str1[i], ":")
			for j := 0; j < len(str2); j++ {
				if str2[j] == "column" {
					return str2[j+1], true
				}
			}
		}
	}
	return "", false
}

func GetListFieldsTagJson(modelType reflect.Type) []string {
	numField := modelType.NumField()
	var idFields []string
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
				jsonTag := field.Tag.Get("json")
				tags1 := strings.Split(jsonTag, ",")
				if len(tags1) > 0 && tags1[0] != "-" {
					idFields = append(idFields, tags1[0])
				}
			}
		}
	}
	return idFields
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

func GetDriver(db *sql.DB) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return DriverPostgres
	case "*mysql.MySQLDriver":
		return DriverMysql
	case "*mssql.Driver":
		return DriverMssql
	case "*godror.drv":
		return DriverOracle
	default:
		return DriverNotSupport
	}
}
