package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const txs = "tx"

func GetTx(ctx context.Context) *sql.Tx {
	txi := ctx.Value(txs)
	if txi != nil {
		txx, ok := txi.(*sql.Tx)
		if ok {
			return txx
		}
	}
	return nil
}
func GetTxId(ctx context.Context) *string {
	txi := ctx.Value("txId")
	if txi != nil {
		txx, ok := txi.(*string)
		if ok {
			return txx
		}
	}
	return nil
}

type Loader struct {
	Database          *sql.DB
	BuildParam        func(i int) string
	Map               func(ctx context.Context, model interface{}) (interface{}, error)
	modelType         reflect.Type
	modelsType        reflect.Type
	keys              []string
	mapJsonColumnKeys map[string]string
	fieldsIndex       map[string]int
	table             string
	IsRollback        bool
	toArray           func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func UseLoadWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (func(context.Context, interface{}, interface{}) (bool, error), error) {
	l, err := NewLoaderWithArray(db, tableName, modelType, toArray, options...)
	if err != nil {
		return nil, err
	}
	return l.LoadAndDecode, nil
}
func UseGetWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (func(context.Context, interface{}, interface{}) (bool, error), error) {
	return UseLoadWithArray(db, tableName, modelType, toArray, options...)
}
func Load(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (func(context.Context, interface{}, interface{}) (bool, error), error) {
	l, err := NewLoader(db, tableName, modelType, options...)
	if err != nil {
		return nil, err
	}
	return l.LoadAndDecode, nil
}
func UseLoad(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (func(context.Context, interface{}, interface{}) (bool, error), error) {
	return Load(db, tableName, modelType, options...)
}
func Get(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (func(context.Context, interface{}, interface{}) (bool, error), error) {
	return Load(db, tableName, modelType, options...)
}
func UseGet(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (func(context.Context, interface{}, interface{}) (bool, error), error) {
	return Load(db, tableName, modelType, options...)
}
func NewLoader(db *sql.DB, tableName string, modelType reflect.Type, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	return NewLoaderWithArray(db, tableName, modelType, nil, options...)
}
func NewLoaderWithArray(db *sql.DB, tableName string, modelType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(context.Context, interface{}) (interface{}, error)) (*Loader, error) {
	var mp func(ctx context.Context, model interface{}) (interface{}, error)
	if len(options) >= 1 {
		mp = options[0]
	}
	return NewSqlLoader(db, tableName, modelType, mp, toArray)
}
func NewSqlLoader(db *sql.DB, tableName string, modelType reflect.Type, mp func(context.Context, interface{}) (interface{}, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, options ...func(i int) string) (*Loader, error) {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	_, idNames := FindPrimaryKeys(modelType)
	mapJsonColumnKeys := MapJsonColumn(modelType)
	modelsType := reflect.Zero(reflect.SliceOf(modelType)).Type()

	fieldsIndex, er0 := GetColumnIndexes(modelType)
	if er0 != nil {
		return nil, er0
	}
	return &Loader{Database: db, IsRollback: true, BuildParam: buildParam, Map: mp, modelType: modelType, modelsType: modelsType, keys: idNames, mapJsonColumnKeys: mapJsonColumnKeys, fieldsIndex: fieldsIndex, table: tableName, toArray: toArray}, nil
}

func (s *Loader) Keys() []string {
	return s.keys
}

func (s *Loader) All(ctx context.Context) (interface{}, error) {
	query := BuildSelectAllQuery(s.table)
	result := reflect.New(s.modelsType).Interface()
	var err error
	tx := GetTx(ctx)
	if tx == nil {
		err = QueryWithArray(ctx, s.Database, s.fieldsIndex, result, s.toArray, query)
	} else {
		err = QueryTxWithArray(ctx, tx, s.fieldsIndex, result, s.toArray, query)
		if err != nil {
			if s.IsRollback {
				tx.Rollback()
			}
			return result, err
		}
	}
	if err == nil {
		if s.Map != nil {
			return MapModels(ctx, result, s.Map)
		}
		return result, err
	}
	return result, err
}

func (s *Loader) Load(ctx context.Context, id interface{}) (interface{}, error) {
	queryFindById, values := BuildFindByIdWithDB(s.Database, s.table, id, s.mapJsonColumnKeys, s.keys, s.BuildParam)
	tx := GetTx(ctx)
	var r interface{}
	var er1 error
	if tx == nil {
		r, er1 = QueryRowWithArray(ctx, s.Database, s.modelType, s.fieldsIndex, s.toArray, queryFindById, values...)
	} else {
		r, er1 = QueryRowTxWithArray(ctx, tx, s.modelType, s.fieldsIndex, s.toArray, queryFindById, values...)
		if er1 != nil {

			return r, er1
		}
	}
	if er1 != nil {
		if s.IsRollback && tx != nil {
			tx.Rollback()
		}
		return r, er1
	}
	if s.Map != nil {
		_, er2 := s.Map(ctx, &r)
		if er2 != nil {
			return r, er2
		}
		return r, er2
	}
	return r, er1
}

func (s *Loader) Exist(ctx context.Context, id interface{}) (bool, error) {
	var count int32
	var where string
	var values []interface{}
	colNumber := 1
	if len(s.keys) == 1 {
		where = fmt.Sprintf("where %s = %s", s.mapJsonColumnKeys[s.keys[0]], s.BuildParam(colNumber))
		values = append(values, id)
		colNumber++
	} else {
		conditions := make([]string, 0)
		var ids = id.(map[string]interface{})
		for k, idk := range ids {
			columnName := s.mapJsonColumnKeys[k]
			conditions = append(conditions, fmt.Sprintf("%s = %s", columnName, s.BuildParam(colNumber)))
			values = append(values, idk)
			colNumber++
		}
		where = "where " + strings.Join(conditions, " and ")
	}
	var row *sql.Row
	tx := GetTx(ctx)
	if tx == nil {
		row = s.Database.QueryRowContext(ctx, fmt.Sprintf("select count(*) from %s %s", s.table, where), values...)
	} else {
		row = tx.QueryRowContext(ctx, fmt.Sprintf("select count(*) from %s %s", s.table, where), values...)
	}
	if err := row.Scan(&count); err != nil {
		if s.IsRollback && tx != nil {
			tx.Rollback()
		}
		return false, err
	} else {
		if count >= 1 {
			return true, nil
		}
		return false, nil
	}
}

func (s *Loader) LoadAndDecode(ctx context.Context, id interface{}, result interface{}) (bool, error) {
	var values []interface{}
	sql, values := BuildFindByIdWithDB(s.Database, s.table, id, s.mapJsonColumnKeys, s.keys, s.BuildParam)
	var rowData interface{}
	var er1 error
	tx := GetTx(ctx)
	if tx == nil {
		rowData, er1 = QueryRowWithArray(ctx, s.Database, s.modelType, s.fieldsIndex, s.toArray, sql, values...)
	} else {
		rowData, er1 = QueryRowTxWithArray(ctx, tx, s.modelType, s.fieldsIndex, s.toArray, sql, values...)
	}
	if er1 != nil && s.IsRollback && tx != nil {
		tx.Rollback()
		return false, er1
	}
	if er1 != nil || rowData == nil {
		return false, er1
	}
	byteData, _ := json.Marshal(rowData)
	er2 := json.Unmarshal(byteData, &result)
	if er2 != nil {
		return false, er2
	}
	//reflect.ValueOf(result).Elem().Set(reflect.ValueOf(rowData).Elem())
	if s.Map != nil {
		_, er3 := s.Map(ctx, result)
		if er3 != nil {
			return true, er3
		}
	}
	return true, nil
}

func MapJsonColumn(modelType reflect.Type) map[string]string {
	numField := modelType.NumField()
	columnNameKeys := make(map[string]string)
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
								tagj, ok1 := field.Tag.Lookup("json")
								t := strings.Split(tagj, ",")
								if ok1 && len(t) > 0 {
									json := t[0]
									columnNameKeys[json] = str2[j+1]
								}
							}
						}
					}
				}
			}
		}
	}
	return columnNameKeys
}
