package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	DBName          = "column"
	PrimaryKey      = "primary_key"
	DriverCassandra = "cassandra"
)

type TxCache interface {
	Put(key string, value *sql.Tx, timeToLive time.Duration) error
	Expire(key string, timeToLive time.Duration) (bool, error)
	Get(key string) (*sql.Tx, error)
	Remove(key string) (bool, error)
	Clear() error
	Keys() ([]string, error)
	Count() (int64, error)
	Size() (int64, error)
}

func InitSingleResult(modelType reflect.Type) interface{} {
	return reflect.New(modelType).Interface()
}

func InitArrayResults(modelsType reflect.Type) interface{} {
	return reflect.New(modelsType).Interface()
}

func ExecStmt(ctx context.Context, stmt *sql.Stmt, values ...interface{}) (int64, error) {
	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func handleDuplicate(db *sql.DB, err error) (int64, error) {
	x := err.Error()
	driver := GetDriver(db)
	if driver == DriverPostgres && strings.Contains(x, "pq: duplicate key value violates unique constraint") {
		return 0, nil
	} else if driver == DriverMysql && strings.Contains(x, "Error 1062: Duplicate entry") {
		return 0, nil //mysql Error 1062: Duplicate entry 'a-1' for key 'PRIMARY'
	} else if driver == DriverOracle && strings.Contains(x, "ORA-00001: unique constraint") {
		return 0, nil //mysql Error 1062: Duplicate entry 'a-1' for key 'PRIMARY'
	} else if driver == DriverMssql && strings.Contains(x, "Violation of PRIMARY KEY constraint") {
		return 0, nil //Violation of PRIMARY KEY constraint 'PK_aa'. Cannot insert duplicate key in object 'dbo.aa'. The duplicate key value is (b, 2).
	} else if driver == DriverSqlite3 && strings.Contains(x, "UNIQUE constraint failed") {
		return 0, nil
	}
	return 0, err
}
func Insert(ctx context.Context, db *sql.DB, table string, model interface{}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return InsertWithVersion(ctx, db, table, model, -1, nil, schema)
}
func InsertWithArray(ctx context.Context, db *sql.DB, table string, model interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return InsertWithVersion(ctx, db, table, model, -1, toArray, schema)
}
func InsertWithVersion(ctx context.Context, db *sql.DB, table string, model interface{}, versionIndex int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, schema *Schema, options ...func(i int) string) (int64, error) {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	driver := GetDriver(db)
	boolSupport := driver == DriverPostgres
	queryInsert, values := BuildToInsertWithVersion(table, model, versionIndex, buildParam, boolSupport, toArray, schema)

	result, err := db.ExecContext(ctx, queryInsert, values...)
	if err != nil {
		return handleDuplicate(db, err)
	}
	return result.RowsAffected()
}
func InsertTx(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return InsertTxWithVersion(ctx, db, tx, table, model, -1, nil, schema)
}
func InsertTxWithArray(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return InsertTxWithVersion(ctx, db, tx, table, model, -1, toArray, schema)
}
func InsertTxWithVersion(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, versionIndex int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, schema *Schema, options ...func(i int) string) (int64, error) {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	driver := GetDriver(db)
	boolSupport := driver == DriverPostgres
	queryInsert, values := BuildToInsertWithSchema(table, model, versionIndex, buildParam, boolSupport, false, toArray, schema)

	result, err := tx.ExecContext(ctx, queryInsert, values...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func Update(ctx context.Context, db *sql.DB, table string, model interface{}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return UpdateWithVersion(ctx, db, table, model, -1, nil, schema)
}
func UpdateWithArray(ctx context.Context, db *sql.DB, table string, model interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return UpdateWithVersion(ctx, db, table, model, -1, toArray, schema)
}
func UpdateWithVersion(ctx context.Context, db *sql.DB, table string, model interface{}, versionIndex int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, schema *Schema, options ...func(i int) string) (int64, error) {
	if versionIndex < 0 {
		return 0, errors.New("version's index not found")
	}
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	driver := GetDriver(db)
	boolSupport := driver == DriverPostgres
	query, values := BuildToUpdateWithVersion(table, model, versionIndex, buildParam, boolSupport, toArray, schema)

	result, err := db.ExecContext(ctx, query, values...)

	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
func UpdateTx(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return UpdateTxWithVersion(ctx, db, tx, table, model, -1, nil, schema)
}
func UpdateTxWithArray(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return UpdateTxWithVersion(ctx, db, tx, table, model, -1, toArray, schema)
}
func UpdateTxWithVersion(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, versionIndex int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, schema *Schema, options ...func(i int) string) (int64, error) {
	if versionIndex < 0 {
		return 0, errors.New("version's index not found")
	}
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	driver := GetDriver(db)
	boolSupport := driver == DriverPostgres
	query, values := BuildToUpdateWithVersion(table, model, versionIndex, buildParam, boolSupport, toArray, schema)

	result, err := tx.ExecContext(ctx, query, values...)

	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
func InsertBatch(ctx context.Context, db *sql.DB, tableName string, models interface{}, options ...*Schema) (int64, error) {
	buildParam := GetBuild(db)
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return InsertBatchWithSchema(ctx, db, tableName, models, nil, buildParam, schema)
}
func InsertBatchWithArray(ctx context.Context, db *sql.DB, tableName string, models interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	buildParam := GetBuild(db)
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return InsertBatchWithSchema(ctx, db, tableName, models, toArray, buildParam, schema)
}
func InsertBatchWithSchema(ctx context.Context, db *sql.DB, tableName string, models interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, buildParam func(i int) string, options ...*Schema) (int64, error) {
	if buildParam == nil {
		buildParam = GetBuild(db)
	}
	driver := GetDriver(db)
	query, args, er1 := BuildToInsertBatchWithSchema(tableName, models, driver, toArray, buildParam, options...)
	if er1 != nil {
		return 0, er1
	}
	x, er2 := db.ExecContext(ctx, query, args...)
	if er2 != nil {
		return 0, er2
	}
	return x.RowsAffected()
}
func UpdateBatch(ctx context.Context, db *sql.DB, tableName string, models interface{}, options ...*Schema) (int64, error) {
	buildParam := GetBuild(db)
	driver := GetDriver(db)
	boolSupport := driver == DriverPostgres
	return UpdateBatchWithVersion(ctx, db, tableName, models, -1, nil, buildParam, boolSupport, options...)
}
func UpdateBatchWithArray(ctx context.Context, db *sql.DB, tableName string, models interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	buildParam := GetBuild(db)
	driver := GetDriver(db)
	boolSupport := driver == DriverPostgres
	return UpdateBatchWithVersion(ctx, db, tableName, models, -1, toArray, buildParam, boolSupport, options...)
}
func UpdateBatchWithVersion(ctx context.Context, db *sql.DB, tableName string, models interface{}, versionIndex int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, buildParam func(int) string, boolSupport bool, options ...*Schema) (int64, error) {
	if buildParam == nil {
		buildParam = GetBuild(db)
	}
	stmts, er1 := BuildToUpdateBatchWithVersion(tableName, models, versionIndex, buildParam, boolSupport, toArray, options...)
	if er1 != nil {
		return 0, er1
	}
	return ExecuteAll(ctx, db, stmts...)
}

func Save(ctx context.Context, db *sql.DB, table string, model interface{}, options ...*Schema) (int64, error) {
	return SaveWithArray(ctx, db, table, model, nil, options...)
}
func SaveWithArray(ctx context.Context, db *sql.DB, table string, model interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	drive := GetDriver(db)
	buildParam := GetBuild(db)
	queryString, value, err := BuildToSaveWithSchema(table, model, drive, buildParam, toArray, options...)
	if err != nil {
		return 0, err
	}
	res, err := db.ExecContext(ctx, queryString, value...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func SaveTx(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return SaveTxWithArray(ctx, db, tx, table, model, nil, schema)
}
func SaveTxWithArray(ctx context.Context, db *sql.DB, tx *sql.Tx, table string, model interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	drive := GetDriver(db)
	buildParam := GetBuild(db)
	queryString, value, err := BuildToSaveWithSchema(table, model, drive, buildParam, toArray, options...)
	if err != nil {
		return 0, err
	}
	res, err := tx.ExecContext(ctx, queryString, value...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func SaveBatch(ctx context.Context, db *sql.DB, tableName string, models interface{}, options ...*Schema) (int64, error) {
	var schema *Schema
	if len(options) > 0 {
		schema = options[0]
	}
	return SaveBatchWithArray(ctx, db, tableName, models, nil, schema)
}
func SaveBatchWithArray(ctx context.Context, db *sql.DB, tableName string, models interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...*Schema) (int64, error) {
	driver := GetDriver(db)
	stmts, er1 := BuildToSaveBatchWithArray(tableName, models, driver, toArray, options...)
	if er1 != nil {
		return 0, er1
	}
	_, err := ExecuteAll(ctx, db, stmts...)
	total := int64(len(stmts))
	return total, err
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
func ExtractBySchema(value interface{}, columns []string, schema map[string]FieldDB) (map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		value = rv.Interface()
	}
	if rv.Kind() != reflect.Struct {
		return nil, nil, nil, errors.New("value must be kind of Struct")
	}

	var attrs = map[string]interface{}{}
	var nAttrs = map[string]interface{}{}
	var attrsKey = map[string]interface{}{}

	for _, col := range columns {
		fdb, ok := schema[col]
		if ok {
			f := rv.Field(fdb.Index)
			fieldValue := f.Interface()
			isNil := false
			if f.Kind() == reflect.Ptr {
				if reflect.ValueOf(fieldValue).IsNil() {
					isNil = true
				} else {
					fieldValue = reflect.Indirect(reflect.ValueOf(fieldValue)).Interface()
				}
			}
			if !fdb.Key {
				if !isNil {
					if boolValue, ok := fieldValue.(bool); ok {
						if boolValue {
							attrs[col] = fdb.True
							nAttrs[col] = fdb.True
						} else {
							attrs[col] = fdb.False
							nAttrs[col] = fdb.False
						}
					} else {
						attrs[col] = fieldValue
						nAttrs[col] = fieldValue
					}
				} else {
					attrs[col] = fieldValue
				}
			} else {
				attrsKey[col] = fieldValue
				if !isNil {
					nAttrs[col] = fieldValue
				}
			}
		}
	}
	return attrs, attrsKey, nAttrs, nil
}

// Obtain columns and values required for insert from interface
func ExtractMapValue(value interface{}, excludeColumns *[]string, ignoreNull bool) (map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		value = rv.Interface()
	}
	if rv.Kind() != reflect.Struct {
		return nil, nil, nil, errors.New("value must be kind of Struct")
	}

	var attrs = map[string]interface{}{}
	var nAttrs = map[string]interface{}{}
	var attrsKey = map[string]interface{}{}

	for index, field := range GetMapField(value) {
		if GetTag(field, IgnoreReadWrite) == IgnoreReadWrite {
			continue
		}
		kind := field.Value.Kind()
		fieldValue := field.Value.Interface()
		isNil := false
		if kind == reflect.Ptr {
			if reflect.ValueOf(fieldValue).IsNil() {
				if ignoreNull {
					*excludeColumns = append(*excludeColumns, field.Tags["fieldName"])
				}
				isNil = true
			} else {
				fieldValue = reflect.Indirect(reflect.ValueOf(fieldValue)).Interface()
			}
		}
		if !ContainString(*excludeColumns, GetTag(field, "fieldName")) && !IsPrimary(field) {
			if dBName, ok := field.Tags[DBName]; ok {
				if !isNil {
					if boolValue, ok := fieldValue.(bool); ok {
						bv := field.Type.Field(index).Tag.Get(strconv.FormatBool(boolValue))
						attrs[dBName] = bv
						nAttrs[dBName] = bv
					} else {
						attrs[dBName] = fieldValue
						nAttrs[dBName] = fieldValue
					}
				} else {
					attrs[dBName] = fieldValue
				}

			}
		}
		if IsPrimary(field) {
			if dBName, ok := field.Tags[DBName]; ok {
				attrsKey[dBName] = fieldValue
				if !isNil {
					nAttrs[dBName] = fieldValue
				}
			}
		}
	}
	return attrs, attrsKey, nAttrs, nil
}

func GetJsonNameByIndex(ModelType reflect.Type, index int) (string, bool) {
	field := ModelType.Field(index)
	if tagJson, ok := field.Tag.Lookup("json"); ok {
		arrValue := strings.Split(tagJson, ",")
		if len(arrValue) > 0 {
			return arrValue[0], true
		}
	}

	return "", false
}

func FindFieldByName(modelType reflect.Type, fieldName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			name1 := fieldName
			name2 := fieldName
			tag1, ok1 := field.Tag.Lookup("json")
			tag2, ok2 := field.Tag.Lookup("gorm")
			if ok1 {
				name1 = strings.Split(tag1, ",")[0]
			}
			if ok2 {
				if has := strings.Contains(tag2, "column"); has {
					str1 := strings.Split(tag2, ";")
					num := len(str1)
					for k := 0; k < num; k++ {
						str2 := strings.Split(str1[k], ":")
						for j := 0; j < len(str2); j++ {
							if str2[j] == "column" {
								return i, name1, str2[j+1]
							}
						}
					}
				}
			}
			return i, name1, name2
		}
	}
	return -1, fieldName, fieldName
}

func FindIdFields(modelType reflect.Type) []string {
	numField := modelType.NumField()
	var idFields []string
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
				idFields = append(idFields, field.Name)
			}
		}
	}
	return idFields
}

func FindIdColumns(modelType reflect.Type) []string {
	numField := modelType.NumField()
	var idFields = make([]string, 0)
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
				if has := strings.Contains(ormTag, "column"); has {
					str1 := strings.Split(ormTag, ";")
					num := len(str1)
					for i := 0; i < num; i++ {
						str2 := strings.Split(str1[i], ":")
						for j := 0; j < len(str2); j++ {
							if str2[j] == "column" {
								idFields = append(idFields, str2[j+1])
							}
						}
					}
				}
			}
		}
	}
	return idFields
}

func MapToGORM(ids map[string]interface{}, modelType reflect.Type) (query map[string]interface{}) {
	queryGen := make(map[string]interface{})
	var columnName string
	for colName, value := range ids {
		columnName, _ = GetColumnName(modelType, colName)
		queryGen[columnName] = value
	}
	return queryGen
}

// For DefaultGenericService
func BuildQueryByIdFromObject(object interface{}, modelType reflect.Type, idNames []string) (query map[string]interface{}) {
	queryGen := make(map[string]interface{})
	var value interface{}
	for _, colId := range idNames {
		value = reflect.Indirect(reflect.ValueOf(object)).FieldByName(colId).Interface()
		queryGen[colId] = value
	}
	return MapToGORM(queryGen, modelType)
}

func BuildQueryByIdFromMap(object map[string]interface{}, modelType reflect.Type, idNames []string) (query map[string]interface{}) {
	queryGen := make(map[string]interface{})
	//var value interface{}
	for _, colId := range idNames {
		queryGen[colId] = object[colId]
	}
	return MapToGORM(queryGen, modelType)
}

// For Search
func GetSqlBuilderTags(modelType reflect.Type) []QueryType {
	numField := modelType.NumField()
	//queries := make([]QueryType, 0)
	var sqlQueries []QueryType
	for i := 0; i < numField; i++ {
		sqlQuery := QueryType{}
		field := modelType.Field(i)
		sqlTags := field.Tag.Get("sql_builder")
		tags := strings.Split(sqlTags, ";")
		for _, tag := range tags {
			key := strings.Split(tag, ":")
			switch key[0] {
			case "join":
				sqlQuery.Join = key[1]
				break
			case "select":
				sqlQuery.Select = key[1]
				break
			case "select_count":
				sqlQuery.SelectCount = key[1]
			}
		}
		if sqlQuery.Select != "" || sqlQuery.Join != "" || sqlQuery.SelectCount != "" {
			sqlQueries = append(sqlQueries, sqlQuery)
		}
	}
	return sqlQueries
}

func MapColumnToJson(query map[string]interface{}) interface{} {
	result := make(map[string]interface{})
	for k, v := range query {
		dem := strings.Count(k, "_")
		for i := 0; i < dem; i++ {
			if strings.Index(k, "_") > -1 {
				hoa := []rune(strings.ToUpper(string(k[strings.Index(k, "_")+1])))
				k = ReplaceAtIndex(k, hoa[0], strings.Index(k, "_")+1)
				k = strings.Replace(k, "_", "", 1)
			}
		}
		result[k] = v
	}
	return result
}
func ReplaceAtIndex(str string, replacement rune, index int) string {
	out := []rune(str)
	out[index] = replacement
	return string(out)
}

func GetTableName(object interface{}) string {
	vo := reflect.Indirect(reflect.ValueOf(object))
	tableName := vo.MethodByName("TableName").Call([]reflect.Value{})
	return tableName[0].String()
}

func EscapeString(value string) string {
	//replace := map[string]string{"'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}
	//if strings.Contains(value, `\\`) {
	//	value = strings.Replace(value, "\\", "\\\\", -1)
	//}
	//for b, a := range replace {
	//	if strings.Contains(value, b) {
	//		value = strings.Replace(value, b, a, -1)
	//	}
	//}
	return strings.NewReplacer("\\", "\\\\", "'", `\'`, "\\0", "\\\\0", "\n", "\\n", "\r", "\\r", `"`, `\"`, "\x1a", "\\Z" /*We have more here*/).Replace(value)
}

func EscapeStringForSelect(value string) string {
	//replace := map[string]string{"'": `''`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}
	//if strings.Contains(value, `\\`) {
	//	value = strings.Replace(value, "\\", "\\\\", -1)
	//}
	//
	//for b, a := range replace {
	//	if strings.Contains(value, b) {
	//		value = strings.Replace(value, b, a, -1)
	//	}
	//}
	return strings.NewReplacer("'", `''` /*We have more here*/).Replace(value)
}

// Check if string value is contained in slice
func ContainString(s []string, value string) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

// Enable map keys to be retrieved in same order when iterating
func SortedKeys(val map[string]interface{}) []string {
	var keys []string
	for key := range val {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func GetTag(field Field, tagName string) string {
	if tag, ok := field.Tags[tagName]; ok {
		return tag
	}
	return ""
}

func IsPrimary(field Field) bool {
	return GetTag(field, PrimaryKey) != ""
}
func ReplaceQueryArgs(driver string, query string) string {
	if driver == DriverOracle || driver == DriverPostgres {
		var x string
		if driver == DriverOracle {
			x = ":"
		} else {
			x = "$"
		}
		i := 1
		k := strings.Index(query, "?")
		if k >= 0 {
			for {
				query = strings.Replace(query, "?", x+fmt.Sprintf("%v", i), 1)
				i = i + 1
				k := strings.Index(query, "?")
				if k < 0 {
					return query
				}
			}
		}
	}
	return query
}
func Exist(ctx context.Context, db *sql.DB, sql string, args ...interface{}) (bool, error) {
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}
func MapModels(ctx context.Context, models interface{}, mp func(context.Context, interface{}) (interface{}, error)) (interface{}, error) {
	vo := reflect.Indirect(reflect.ValueOf(models))
	if vo.Kind() == reflect.Ptr {
		vo = reflect.Indirect(vo)
	}
	if vo.Kind() == reflect.Slice {
		le := vo.Len()
		for i := 0; i < le; i++ {
			x := vo.Index(i)
			k := x.Kind()
			if k == reflect.Struct {
				y := x.Addr().Interface()
				mp(ctx, y)
			} else {
				y := x.Interface()
				mp(ctx, y)
			}

		}
	}
	return models, nil
}
func HandleError(tx *sql.Tx, err *error) {
	if re := recover(); re != nil {
		tx.Rollback()
	} else {
		if *err != nil {
			tx.Rollback()
		} else {
			er := tx.Commit()
			err = &er
		}
	}
}
