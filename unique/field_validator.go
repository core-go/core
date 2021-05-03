package unique

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	sv "github.com/core-go/service"
)

const (
	defaultPagingFormat = " limit %s offset %s "
	oraclePagingFormat  = " offset %s rows fetch next %s rows only "

	driverPostgres   = "postgres"
	driverMysql      = "mysql"
	driverMssql      = "mssql"
	driverOracle     = "oracle"
	driverSqlite3    = "sqlite3"
	driverNotSupport = "no support"
)

type FieldValidator struct {
	db             *sql.DB
	driver         string
	validate       func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	modelType      reflect.Type
	tableName      string
	fieldName      string
	fieldIndex     int
	columnName     string
	idColumnFields []string
	keyIndexes     map[string]int
	query          string
}

func NewUniqueFieldValidator(db *sql.DB, validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error), tableName string, fieldName string, modelType reflect.Type) *FieldValidator {
	driver := getDriver(db)
	keyIndexes, _ := getColumnIndexes(modelType)
	idColumnFieldsName, _ := findNames(modelType)

	return &FieldValidator{
		db:             db,
		driver:         driver,
		validate:       validate,
		modelType:      modelType,
		tableName:      tableName,
		fieldName:      fieldName,
		columnName:     "",
		idColumnFields: idColumnFieldsName,
		keyIndexes:     keyIndexes,
		query:          "",
	}
}
func (v *FieldValidator) Validate(ctx context.Context, model interface{}) (errors []sv.ErrorMessage, err error) {
	errors, err = v.validate(ctx, model)
	if err != nil {
		return errors, err
	}

	vo := reflect.Indirect(reflect.ValueOf(model))
	if vo.Kind() == reflect.Ptr {
		vo = reflect.Indirect(vo)
	}
	updateStatus, valuesId := IsIdValid(v.keyIndexes, v.idColumnFields, vo)
	values := BuildParameters(v.keyIndexes, vo, v.fieldName, valuesId, updateStatus)
	syntax := GetDriverParam(v.driver, values)
	query := BuildQuery(v.tableName, v.fieldName, v.idColumnFields, syntax, v.driver, updateStatus)

	rows, err := v.db.Query(query, values...)
	if err != nil {
		return errors, err
	}

	for rows.Next() {
		er := sv.ErrorMessage{Field: v.fieldName, Code: "duplicate"}
		return append(errors, er), nil
	}
	return errors, err
}

func IsIdValid(keyIndexes map[string]int, idColumnFields []string, modelType reflect.Value) (bool, []interface{}) {
	var valuesID []interface{}
	for _, field := range idColumnFields {
		index := keyIndexes[field]
		fieldId := modelType.Field(index)
		if fieldId.IsValid() {
			valuesID = append(valuesID, fieldId.Interface())
		} else {
			return false, nil
		}
	}
	return true, valuesID
}

func BuildQuery(tableName string, fieldsName string, idColumnFields []string, syntax []string, driver string, updateStatus bool) string {
	var update string
	query := fmt.Sprintf("select %s from %s", fieldsName, tableName) + " where " + fieldsName + fmt.Sprintf(" = %s", syntax[0])
	n := len(idColumnFields) - 1
	if updateStatus {
		for i, id := range idColumnFields {
			var u string
			if i < n {
				u = fmt.Sprintf(" and %s = %s ", id, syntax[i+1])
			} else if i == n {
				u = fmt.Sprintf(" and %s != %s ", id, syntax[i+1])
			}
			update += u
		}
	}

	var limit string
	if driver == driverOracle {
		limit = fmt.Sprintf(oraclePagingFormat, "0", "1")
	} else {
		limit = fmt.Sprintf(defaultPagingFormat, "1", "0")
	}

	return query + update + limit
}

func BuildParameters(keyIndexes map[string]int, modelType reflect.Value, fieldsName string, valuesID []interface{}, updateStatus bool) []interface{} {
	var values []interface{}
	index := keyIndexes[fieldsName]
	if updateStatus {
		values = append(values, modelType.Field(index).Interface())
		for _, id := range valuesID {
			values = append(values, id)
		}
	} else {
		values = append(values, modelType.Field(index).Interface())
	}
	return values
}
func GetDriverParam(driver string, values []interface{}) []string {
	var syntax []string
	for i := 0; i < len(values); i++ {
		var s string
		if driver == "postgres" {
			s = "$" + strconv.Itoa(i+1)
		} else if driver == "oracle" {
			s = ":val" + strconv.Itoa(i+1)
		} else if driver == "mssql" {
			s = "@p" + strconv.Itoa(i+1)
		} else {
			s = "?"
		}
		syntax = append(syntax, s)
	}
	return syntax
}
func getColumnIndexes(modelType reflect.Type) (map[string]int, error) {
	ma := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := findTag(ormTag, "column")
		column = strings.ToLower(column)
		if ok {
			ma[column] = i
		}
	}
	return ma, nil
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
func findNames(modelType reflect.Type) ([]string, []string) {
	numField := modelType.NumField()
	var idColumnFields []string
	var idJsons []string
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
				k, ok := findPrivateTag(ormTag, "column")
				if ok {
					idColumnFields = append(idColumnFields, k)
					tag1, ok1 := field.Tag.Lookup("json")
					tagJsons := strings.Split(tag1, ",")
					if ok1 && len(tagJsons) > 0 {
						idJsons = append(idJsons, tagJsons[0])
					}
				}
			}
		}
	}
	return idColumnFields, idJsons
}
func findPrivateTag(tag string, key string) (string, bool) {
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
