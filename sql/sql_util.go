package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
)

const IgnoreReadWrite = "-"
const txs = "tx"

type Executor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
func GetExec(ctx context.Context, db *sql.DB, opts...string) Executor {
	name := txs
	if len(opts) > 0 && len(opts[0]) > 0 {
		name = opts[0]
	}
	txi := ctx.Value(name)
	if txi != nil {
		txx, ok := txi.(*sql.Tx)
		if ok {
			return txx
		}
	}
	return db
}
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
func Count(ctx context.Context, db *sql.DB, sql string, values ...interface{}) (int64, error) {
	var total int64
	row := db.QueryRowContext(ctx, sql, values...)
	err2 := row.Scan(&total)
	if err2 != nil {
		return total, err2
	}
	return total, nil
}
func QueryMapWithTx(ctx context.Context, db *sql.Tx, transform func(s string) string, sql string, values ...interface{}) ([]map[string]interface{}, error) {
	rows, er1 := db.QueryContext(ctx, sql, values...)
	if er1 != nil {
		return nil, er1
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	colMaps := make([]string, len(cols))
	if transform != nil {
		for i, colName := range cols {
			colMaps[i] = transform(colName)
		}
	} else {
		for i, colName := range cols {
			colMaps[i] = colName
		}
	}
	res := make([]map[string]interface{}, 0)
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, _ := range cols {
			val := columnPointers[i].(*interface{})
			x := *val
			switch s := x.(type) {
			case *[]byte:
				x2 := *s
				s2 := string(x2)
				m[colMaps[i]] = s2
			case []byte:
				s2 := string(s)
				m[colMaps[i]] = s2
			default:
				m[colMaps[i]] = *val
			}
		}
		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		res = append(res, m)
	}
	return res, nil
}
func QueryMap(ctx context.Context, db *sql.DB, transform func(s string) string, sql string, values ...interface{}) ([]map[string]interface{}, error) {
	rows, er1 := db.QueryContext(ctx, sql, values...)
	if er1 != nil {
		return nil, er1
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	colMaps := make([]string, len(cols))
	if transform != nil {
		for i, colName := range cols {
			colMaps[i] = transform(colName)
		}
	} else {
		for i, colName := range cols {
			colMaps[i] = colName
		}
	}
	res := make([]map[string]interface{}, 0)
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, _ := range cols {
			val := columnPointers[i].(*interface{})
			x := *val
			switch s := x.(type) {
			case *[]byte:
				x2 := *s
				s2 := string(x2)
				m[colMaps[i]] = s2
			case []byte:
				s2 := string(s)
				m[colMaps[i]] = s2
			default:
				m[colMaps[i]] = *val
			}
		}
		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		res = append(res, m)
	}
	return res, nil
}
func Query(ctx context.Context, db Executor, fieldsIndex map[string]int, results interface{}, sql string, values ...interface{}) error {
	return QueryWithArray(ctx, db, fieldsIndex, results, nil, sql, values...)
}
func ExecContext(ctx context.Context, db Executor, query string, args ...interface{}) (sql.Result, error){
	tx := GetTx(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	} else {
		return db.ExecContext(ctx, query, args...)
	}
}
func Exec(ctx context.Context, db *sql.DB, query string, args ...interface{}) (int64, error){
	tx := GetTx(ctx)
	if tx != nil {
		res, err := tx.ExecContext(ctx, query, args...)
		return RowsAffected(res, err)
	} else {
		res, err := db.ExecContext(ctx, query, args...)
		return RowsAffected(res, err)
	}
}
func RowsAffected(res sql.Result, err error) (int64, error) {
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
func SelectContext(ctx context.Context, db *sql.DB, results interface{}, sql string, values ...interface{}) error {
	return SelectContextWithArray(ctx, db, results, nil, sql, values...)
}
func SelectContextWithArray(ctx context.Context, db *sql.DB, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values ...interface{}) error {
	return QueryContextWithArray(ctx, db, nil, results, toArray, sql, values...)
}
func QueryContext(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, results interface{}, query string, values ...interface{}) error {
	return QueryContextWithArray(ctx, db, fieldsIndex, results, nil, query, values...)
}
func QueryContextWithMap(ctx context.Context, db *sql.DB, results interface{}, sql string, values []interface{}, options...map[string]int) error {
	return QueryContextWithMapAndArray(ctx, db, results, nil, sql, values, options...)
}
func QueryContextWithMapAndArray(ctx context.Context, db *sql.DB, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values []interface{}, options...map[string]int) error {
	var fieldsIndex map[string]int
	if len(options) > 0 && options[0] != nil {
		fieldsIndex = options[0]
	}
	return QueryContextWithArray(ctx, db, fieldsIndex, results, toArray, sql, values...)
}
func QueryContextWithArray(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, query string, values ...interface{}) error {
	var rows *sql.Rows
	var er1 error
	tx := GetTx(ctx)
	if tx != nil {
		rows, er1 = tx.QueryContext(ctx, query, values...)
	} else {
		rows, er1 = db.QueryContext(ctx, query, values...)
	}
	if er1 != nil {
		return er1
	}
	defer rows.Close()
	modelType := reflect.TypeOf(results).Elem().Elem()
	tb, er3 := Scan(rows, modelType, fieldsIndex, toArray)
	if er3 != nil {
		return er3
	}
	for _, element := range tb {
		appendToArray(results, element)
	}
	er4 := rows.Close()
	if er4 != nil {
		return er4
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er5 := rows.Err(); er5 != nil {
		return er5
	}
	return nil
}
func QueryWithArray(ctx context.Context, db Executor, fieldsIndex map[string]int, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values ...interface{}) error {
	rows, er1 := db.QueryContext(ctx, sql, values...)
	if er1 != nil {
		return er1
	}
	defer rows.Close()
	modelType := reflect.TypeOf(results).Elem().Elem()
	tb, er3 := Scan(rows, modelType, fieldsIndex, toArray)
	if er3 != nil {
		return er3
	}
	for _, element := range tb {
		appendToArray(results, element)
	}
	er4 := rows.Close()
	if er4 != nil {
		return er4
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er5 := rows.Err(); er5 != nil {
		return er5
	}
	return nil
}
func QueryWithMap(ctx context.Context, db Executor, results interface{}, sql string, values []interface{}, options...map[string]int) error {
	return QueryWithMapAndArray(ctx, db, results, nil, sql, values, options...)
}
func QueryWithMapAndArray(ctx context.Context, db Executor, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values []interface{}, options...map[string]int) error {
	var fieldsIndex map[string]int
	if len(options) > 0 && options[0] != nil {
		fieldsIndex = options[0]
	}
	return QueryWithArray(ctx, db, fieldsIndex, results, toArray, sql, values...)
}
func Select(ctx context.Context, db Executor, results interface{}, sql string, values ...interface{}) error {
	return SelectWithArray(ctx, db, results, nil, sql, values...)
}
func SelectWithArray(ctx context.Context, db Executor, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values ...interface{}) error {
	return QueryWithArray(ctx, db, nil, results, toArray, sql, values...)
}
func QueryTx(ctx context.Context, tx *sql.Tx, fieldsIndex map[string]int, results interface{}, sql string, values ...interface{}) error {
	return QueryTxWithArray(ctx, tx, fieldsIndex, results, nil, sql, values...)
}
func QueryTxWithArray(ctx context.Context, tx *sql.Tx, fieldsIndex map[string]int, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values ...interface{}) error {
	rows, er1 := tx.QueryContext(ctx, sql, values...)
	if er1 != nil {
		return er1
	}
	defer rows.Close()

	modelType := reflect.TypeOf(results).Elem().Elem()
	tb, er3 := Scan(rows, modelType, fieldsIndex, toArray)
	if er3 != nil {
		return er3
	}
	for _, element := range tb {
		appendToArray(results, element)
	}
	er4 := rows.Close()
	if er4 != nil {
		return er4
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er5 := rows.Err(); er5 != nil {
		return er5
	}
	return nil
}
func QueryAndCount(ctx context.Context, db Executor, fieldsIndex map[string]int, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, count *int64, sql string, values ...interface{}) error {
	rows, er1 := db.QueryContext(ctx, sql, values...)
	if er1 != nil {
		return er1
	}
	defer rows.Close()
	modelType := reflect.TypeOf(results).Elem().Elem()

	if fieldsIndex == nil {
		fieldsIndex, er1 = GetColumnIndexes(modelType)
		if er1 != nil {
			return er1
		}
	}

	tb, c, er3 := ScanAndCount(rows, modelType, fieldsIndex, toArray)
	*count = c
	if er3 != nil {
		return er3
	}
	for _, element := range tb {
		appendToArray(results, element)
	}
	er4 := rows.Close()
	if er4 != nil {
		return er4
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er5 := rows.Err(); er5 != nil {
		return er5
	}
	return nil
}
func QueryRow(ctx context.Context, db *sql.DB, modelType reflect.Type, fieldsIndex map[string]int, sql string, values ...interface{}) (interface{}, error) {
	return QueryRowWithArray(ctx, db, modelType, fieldsIndex, nil, sql, values...)
}
func QueryRowWithArray(ctx context.Context, db *sql.DB, modelType reflect.Type, fieldsIndex map[string]int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values ...interface{}) (interface{}, error) {
	strSQL := "limit 1"
	driver := GetDriver(db)
	if driver == DriverOracle {
		strSQL = "AND ROWNUM = 1"
	}
	s := sql + " " + strSQL
	rows, er1 := db.QueryContext(ctx, s, values...)
	if er1 != nil {
		return nil, er1
	}
	tb, er2 := Scan(rows, modelType, fieldsIndex, toArray)
	if er2 != nil {
		return nil, er2
	}
	er3 := rows.Close()
	if er3 != nil {
		return nil, er3
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er4 := rows.Err(); er4 != nil {
		return nil, er3
	}
	if len(tb) == 0 {
		return nil, nil
	} else {
		return tb[0], nil
	}
}
func QueryRowTx(ctx context.Context, tx *sql.Tx, modelType reflect.Type, fieldsIndex map[string]int, sql string, values ...interface{}) (interface{}, error) {
	return QueryRowTxWithArray(ctx, tx, modelType, fieldsIndex, nil, sql, values...)
}
func QueryRowTxWithArray(ctx context.Context, tx *sql.Tx, modelType reflect.Type, fieldsIndex map[string]int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values ...interface{}) (interface{}, error) {
	rows, er1 := tx.QueryContext(ctx, sql, values...)
	if er1 != nil {
		return nil, er1
	}
	tb, er2 := Scan(rows, modelType, fieldsIndex, toArray)
	if er2 != nil {
		return nil, er2
	}
	er3 := rows.Close()
	if er3 != nil {
		return nil, er3
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er4 := rows.Err(); er4 != nil {
		return nil, er3
	}
	return tb, nil
}
func QueryRowByStatement(ctx context.Context, stm *sql.Stmt, modelType reflect.Type, fieldsIndex map[string]int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, values ...interface{}) (interface{}, error) {
	rows, er1 := stm.QueryContext(ctx, values...)
	// rows, er1 := db.Query(s, values...)
	if er1 != nil {
		return nil, er1
	}
	tb, er2 := Scan(rows, modelType, fieldsIndex, toArray)
	if er2 != nil {
		return nil, er2
	}
	er3 := rows.Close()
	if er3 != nil {
		return nil, er3
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er4 := rows.Err(); er4 != nil {
		return nil, er3
	}
	return tb, nil
}
func appendToArray(arr interface{}, item interface{}) interface{} {
	arrValue := reflect.ValueOf(arr)
	elemValue := reflect.Indirect(arrValue)

	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = reflect.Indirect(itemValue)
	}
	elemValue.Set(reflect.Append(elemValue, itemValue))
	return arr
}
func GetColumnIndexes(modelType reflect.Type) (map[string]int, error) {
	ma := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := FindTag(ormTag, "column")
		column = strings.ToLower(column)
		if ok {
			ma[column] = i
		}
	}
	return ma, nil
}

func GetIndexesByTagJson(modelType reflect.Type) (map[string]int, error) {
	ma := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagJson := field.Tag.Get("json")
		if len(tagJson) > 0 {
			ma[tagJson] = i
		}
	}
	return ma, nil
}

func FindTag(tag string, key string) (string, bool) {
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

func GetColumnsSelect(modelType reflect.Type) []string {
	numField := modelType.NumField()
	columnNameKeys := make([]string, 0)
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		if has := strings.Contains(ormTag, "column"); has {
			str1 := strings.Split(ormTag, ";")
			num := len(str1)
			for i := 0; i < num; i++ {
				str2 := strings.Split(str1[i], ":")
				for j := 0; j < len(str2); j++ {
					if str2[j] == "column" {
						columnName := strings.ToLower(str2[j+1])
						columnNameKeys = append(columnNameKeys, columnName)
					}
				}
			}
		}
	}
	return columnNameKeys
}
func GetColumns(cols []string, err error) ([]string, error) {
	if cols == nil || err != nil {
		return cols, err
	}
	c2 := make([]string, 0)
	for _, c := range cols {
		s := strings.ToLower(c)
		c2 = append(c2, s)
	}
	return c2, nil
}
func Scan(rows *sql.Rows, modelType reflect.Type, fieldsIndex map[string]int, options... func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (t []interface{}, err error) {
	if fieldsIndex == nil {
		fieldsIndex, err = GetColumnIndexes(modelType)
		if err != nil {
			return
		}
	}
	var toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	if len(options) > 0 {
		toArray = options[0]
	}
	columns, er0 := GetColumns(rows.Columns())
	if er0 != nil {
		return nil, er0
	}
	for rows.Next() {
		initModel := reflect.New(modelType).Interface()
		r, swapValues := StructScan(initModel, columns, fieldsIndex, toArray)
		if err = rows.Scan(r...); err == nil {
			SwapValuesToBool(initModel, &swapValues)
			t = append(t, initModel)
		}
	}
	return
}
func ScanRow(rows *sql.Rows, s interface{}, columns []string, fieldsIndex map[string]int, options...func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) error {
	r, swapValues := StructScan(s, columns, fieldsIndex, options...)
	err := rows.Scan(r...)
	if err == nil {
		SwapValuesToBool(s, &swapValues)
	}
	return err
}
func StructScan(s interface{}, columns []string, fieldsIndex map[string]int, options...func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (r []interface{}, swapValues map[int]interface{}) {
	var toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	if len(options) > 0 {
		toArray = options[0]
	}
	return StructScanAndIgnore(s, columns, fieldsIndex, toArray, -1)
}
func StructScanAndIgnore(s interface{}, columns []string, fieldsIndex map[string]int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, indexIgnore int) (r []interface{}, swapValues map[int]interface{}) {
	if s != nil {
		modelType := reflect.TypeOf(s).Elem()
		swapValues = make(map[int]interface{}, 0)
		maps := reflect.Indirect(reflect.ValueOf(s))

		if columns == nil {
			for i := 0; i < maps.NumField(); i++ {
				tagBool := modelType.Field(i).Tag.Get("true")
				if tagBool == "" {
					r = append(r, maps.Field(i).Addr().Interface())
				} else {
					var str string
					swapValues[i] = reflect.New(reflect.TypeOf(str)).Elem().Addr().Interface()
					r = append(r, swapValues[i])
				}
			}
			return
		}

		for i, columnsName := range columns {
			if i == indexIgnore {
				continue
			}
			var index int
			var ok bool
			var modelField reflect.StructField
			var valueField reflect.Value
			if fieldsIndex == nil {
				if modelField, ok = modelType.FieldByName(columnsName); !ok {
					var t interface{}
					r = append(r, &t)
					continue
				}
				valueField = maps.FieldByName(columnsName)
			} else {
				if index, ok = fieldsIndex[columnsName]; !ok {
					var t interface{}
					r = append(r, &t)
					continue
				}
				modelField = modelType.Field(index)
				valueField = maps.Field(index)
			}
			x := valueField.Addr().Interface()
			tagBool := modelField.Tag.Get("true")
			if tagBool == "" {
				if toArray != nil && valueField.Kind() == reflect.Slice {
					x = toArray(x)
				}
				r = append(r, x)
			} else {
				var str string
				y := reflect.New(reflect.TypeOf(str))
				swapValues[index] = y.Elem().Addr().Interface()
				r = append(r, swapValues[index])
			}
		}
	}
	return
}
func SwapValuesToBool(s interface{}, swap *map[int]interface{}) {
	if s != nil {
		modelType := reflect.TypeOf(s).Elem()
		maps := reflect.Indirect(reflect.ValueOf(s))
		for index, element := range *swap {
			dbValue2, ok2 := element.(*bool)
			if ok2 {
				if maps.Field(index).Kind() == reflect.Ptr {
					maps.Field(index).Set(reflect.ValueOf(dbValue2))
				} else {
					maps.Field(index).SetBool(*dbValue2)
				}
			} else {
				dbValue, ok := element.(*string)
				if ok {
					var isBool bool
					if *dbValue == "true" {
						isBool = true
					} else if *dbValue == "false" {
						isBool = false
					} else {
						boolStr := modelType.Field(index).Tag.Get("true")
						isBool = *dbValue == boolStr
					}
					if maps.Field(index).Kind() == reflect.Ptr {
						maps.Field(index).Set(reflect.ValueOf(&isBool))
					} else {
						maps.Field(index).SetBool(isBool)
					}
				}
			}
		}
	}
}
func ScanAndCount(rows *sql.Rows, modelType reflect.Type, fieldsIndex map[string]int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) ([]interface{}, int64, error) {
	var t []interface{}
	columns, er0 := GetColumns(rows.Columns())
	if er0 != nil {
		return nil, 0, er0
	}
	if fieldsIndex == nil {
		fieldsIndex, er0 = GetColumnIndexes(modelType)
		if er0 != nil {
			return nil, 0, er0
		}
	}
	var count int64
	for rows.Next() {
		initModel := reflect.New(modelType).Interface()
		var c []interface{}
		c = append(c, &count)
		r, swapValues := StructScanAndIgnore(initModel, columns, fieldsIndex, toArray, 0)
		c = append(c, r...)
		if err := rows.Scan(c...); err == nil {
			SwapValuesToBool(initModel, &swapValues)
			t = append(t, initModel)
		}
	}
	return t, count, nil
}

func ScanRowsWithArray(rows *sql.Rows, structType reflect.Type, fieldsIndex map[string]int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (t interface{}, err error) {
	columns, er0 := GetColumns(rows.Columns())
	err = er0
	if err != nil {
		return
	}
	if fieldsIndex == nil {
		fieldsIndex, er0 = GetColumnIndexes(structType)
		if er0 != nil {
			err = er0
			return
		}
	}
	for rows.Next() {
		gTb := reflect.New(structType).Interface()
		r, swapValues := StructScanAndIgnore(gTb, columns, fieldsIndex, toArray, -1)
		if err = rows.Scan(r...); err == nil {
			SwapValuesToBool(gTb, &swapValues)
			t = gTb
			break
		}
	}
	return
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
//Row
func ScanRowWithArray(row *sql.Row, structType reflect.Type, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (t interface{}, err error) {
	t = reflect.New(structType).Interface()
	r, swapValues := StructScan(t, nil, nil, toArray)
	err = row.Scan(r...)
	SwapValuesToBool(t, &swapValues)
	return
}
func ToCamelCase(s string) string {
	s2 := strings.ToLower(s)
	s1 := string(s2[0])
	for i := 1; i < len(s); i++ {
		if string(s2[i-1]) == "_" {
			s1 = s1[:len(s1)-1]
			s1 += strings.ToUpper(string(s2[i]))
		} else {
			s1 += string(s2[i])
		}
	}
	return s1
}
