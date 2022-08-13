package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var formatDateUpdate = "2006-01-02 15:04:05"

func newStatement(value interface{}, excludeColumns ...string) BatchStatement {
	attribute, attributeKey, _, _ := ExtractMapValue(value, &excludeColumns, false)
	attrSize := len(attribute)
	modelType := reflect.TypeOf(value)
	keys := FindDBColumNames(modelType)
	// Replace with database column name
	dbColumns := make([]string, 0, attrSize)
	for _, key := range SortedKeys(attribute) {
		dbColumns = append(dbColumns, QuoteColumnName(key))
	}
	// Scope to eventually run SQL
	statement := BatchStatement{Keys: keys, Columns: dbColumns, Attributes: attribute, AttributeKeys: attributeKey}
	return statement
}
func FindDBColumNames(modelType reflect.Type) []string {
	numField := modelType.NumField()
	var idFields []string
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		tags := strings.Split(ormTag, ";")
		for _, tag := range tags {
			if strings.Compare(strings.TrimSpace(tag), "primary_key") == 0 {
				k, ok := findTag(ormTag, "column")
				if ok {
					idFields = append(idFields, k)
				}
			}
		}
	}
	return idFields
}
// Field model field definition
type Field struct {
	Tags  map[string]string
	Value reflect.Value
	Type  reflect.Type
}

func GetMapField(object interface{}) []Field {
	modelType := reflect.TypeOf(object)
	value := reflect.Indirect(reflect.ValueOf(object))
	var result []Field
	numField := modelType.NumField()

	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		selectField := Field{Value: value.Field(i), Type: modelType}
		gormTag, ok := field.Tag.Lookup("gorm")
		tag := make(map[string]string)
		tag["fieldName"] = field.Name
		if ok {
			str1 := strings.Split(gormTag, ";")
			for k := 0; k < len(str1); k++ {
				str2 := strings.Split(str1[k], ":")
				if len(str2) == 1 {
					tag[str2[0]] = str2[0]
					selectField.Tags = tag
				} else {
					tag[str2[0]] = str2[1]
					selectField.Tags = tag
				}
			}
			result = append(result, selectField)
		}
	}
	return result
}

func InsertManyWithSize(ctx context.Context, db *sql.DB, tableName string, objects []interface{}, chunkSize int, buildParam func(i int) string, excludeColumns ...string) (int64, error) {
	// Split records with specified size not to exceed Database parameter limit
	if chunkSize <= 0 {
		chunkSize = len(objects)
	}
	var c int64 = 0
	for _, objSet := range splitObjects(objects, chunkSize) {
		count, err := InsertManyRaw(ctx, db, tableName, objSet, false, buildParam, excludeColumns...)
		c = c + count
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func TransactionInsertMany(ctx context.Context, db *sql.DB, tableName string, objects []interface{}, chunkSize int, buildParam func(i int) string, excludeColumns ...string) (int64, error) {
	// Split records with specified size not to exceed Database parameter limit
	if chunkSize <= 0 {
		chunkSize = len(objects)
	}
	var c int64 = 0
	for _, objSet := range splitObjects(objects, chunkSize) {
		count, err := InsertInTransaction(ctx, db, tableName, objSet, false, buildParam, excludeColumns...)
		c = c + count
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func InsertInTransactionTx(ctx context.Context, db *sql.DB, tx *sql.Tx, tableName string, objects interface{}, skipDuplicate bool, buildParam func(i int) string, excludeColumns ...string) (int64, error) {
	objectValues := reflect.Indirect(reflect.ValueOf(objects))
	if objectValues.Kind() == reflect.Slice {
		if objectValues.Len() == 0 {
			return 0, nil
		}
		driver := GetDriver(db)
		firstObj := objectValues.Index(0).Interface()
		columns, primaryKeys, _, err := ExtractMapValue(firstObj, &excludeColumns, true)
		if err != nil {
			return 0, err
		}
		// append primaryKey
		for key, value := range primaryKeys {
			columns[key] = value
		}
		modelType := reflect.TypeOf(objectValues.Index(0).Interface())
		pkey := FindIdFields(modelType)
		attrSize := len(columns)
		// Replace with database column name
		dbColumns := make([]string, 0, attrSize)
		for _, key := range SortedKeys(columns) {
			dbColumns = append(dbColumns, QuoteColumnName(key))
		}
		var start int
		for i := 0; i < objectValues.Len(); i++ {
			obj := objectValues.Index(i).Interface()
			// Store placeholders for embedding variables
			placeholders := make([]string, 0, attrSize)
			objAttrs, primaryKeys, _, err := ExtractMapValue(obj, &excludeColumns, true)
			if err != nil {
				return 0, err
			}
			// append primaryKey
			for key, value := range primaryKeys {
				objAttrs[key] = value
			}
			// If object sizes are different, SQL statement loses consistency
			if len(objAttrs) != attrSize {
				return 0, errors.New("attribute sizes are inconsistent")
			}
			scope := BatchStatement{}
			// Append variables
			variables := make([]string, 0, attrSize)
			for _, key := range SortedKeys(objAttrs) {
				scope.Values = append(scope.Values, objAttrs[key])
				variables = append(variables, BuildParametersFrom(start, 1, buildParam))
				start++
			}

			valueQuery := "(" + strings.Join(variables, ", ") + ")"
			placeholders = append(placeholders, valueQuery)
			var query string
			if skipDuplicate {
				if driver == DriverPostgres {
					query = fmt.Sprintf("insert into %s (%s) values %s on conflict do nothing",
						tableName,
						strings.Join(dbColumns, ", "),
						strings.Join(placeholders, ", "),
					)
				} else if driver == DriverSqlite3 {
					query = fmt.Sprintf("insert or ignore into %s (%s) values %s",
						tableName,
						strings.Join(dbColumns, ", "),
						strings.Join(placeholders, ", "),
					)
				} else if driver == DriverOracle || driver == DriverMysql {
					var qKey []string
					for _, i2 := range pkey {
						key := i2 + " = " + i2
						qKey = append(qKey, key)
					}
					query = fmt.Sprintf("insert into %s (%s) values %s on duplicate key update %s",
						tableName,
						strings.Join(dbColumns, ", "),
						strings.Join(placeholders, ", "),
						strings.Join(qKey, ", "),
					)
				} else {
					return 0, fmt.Errorf("only support skip duplicate on mysql and postgresql, current vendor is %s", driver)
				}
			} else {
				query = fmt.Sprintf("insert into %s (%s) values %s",
					tableName,
					strings.Join(dbColumns, ", "),
					strings.Join(placeholders, ", "),
				)
			}
			_, execErr := tx.ExecContext(ctx, query, scope.Values...)
			if execErr != nil {
				_ = tx.Rollback()
				return 0, execErr
			}
		}
		count := objectValues.Len()
		return int64(count), err
	}
	return 0, fmt.Errorf("objects must be slice")
}

// Separate objects into several size
func splitObjects(objArr []interface{}, size int) [][]interface{} {
	var chunkSet [][]interface{}
	var chunk []interface{}

	for len(objArr) > size {
		chunk, objArr = objArr[:size], objArr[size:]
		chunkSet = append(chunkSet, chunk)
	}
	if len(objArr) > 0 {
		chunkSet = append(chunkSet, objArr[:])
	}

	return chunkSet
}

func InsertManySkipErrors(ctx context.Context, db *sql.DB, tableName string, objects []interface{}, chunkSize int, buildParam func(i int) string, excludeColumns ...string) (int64, error) {
	// Split records with specified size not to exceed Database parameter limit
	if chunkSize <= 0 {
		chunkSize = len(objects)
	}
	var c int64 = 0
	for _, objSet := range splitObjects(objects, chunkSize) {
		count, err := InsertManyRaw(ctx, db, tableName, objSet, true, buildParam, excludeColumns...)
		c = c + count
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func InsertManyRaw(ctx context.Context, db *sql.DB, tableName string, objects []interface{}, skipDuplicate bool, buildParam func(i int) string, excludeColumns ...string) (int64, error) {
	if len(objects) == 0 {
		return 0, nil
	}
	driver := GetDriver(db)
	firstAttrs, primaryKeys, _, err := ExtractMapValue(objects[0], &excludeColumns, true)
	if err != nil {
		return 0, err
	}
	// append primaryKey
	for key, value := range primaryKeys {
		firstAttrs[key] = value
	}

	attrSize := len(firstAttrs)
	modelType := reflect.TypeOf(objects[0])
	pkey := FindIdFields(modelType)
	// Scope to eventually run SQL
	mainScope := BatchStatement{}
	// Store placeholders for embedding variables
	placeholders := make([]string, 0, attrSize)

	// Replace with database column name
	dbColumns := make([]string, 0, attrSize)
	for _, key := range SortedKeys(firstAttrs) {
		dbColumns = append(dbColumns, QuoteColumnName(key))
	}
	var start int
	for _, obj := range objects {
		objAttrs, keys, _, err := ExtractMapValue(obj, &excludeColumns, true)
		if err != nil {
			return 0, err
		}
		for key, value := range keys {
			objAttrs[key] = value
		}

		// If object sizes are different, SQL statement loses consistency
		if len(objAttrs) != attrSize {
			return 0, errors.New("attribute sizes are inconsistent")
		}

		scope := BatchStatement{}

		// Append variables
		variables := make([]string, 0, attrSize)
		for _, key := range SortedKeys(objAttrs) {
			scope.Values = append(scope.Values, objAttrs[key])
			variables = append(variables, BuildParametersFrom(start, 1, buildParam))
			start++
		}

		valueQuery := "(" + strings.Join(variables, ", ") + ")"
		placeholders = append(placeholders, valueQuery)

		// Also append variables to mainScope
		mainScope.Values = append(mainScope.Values, scope.Values...)
	}
	var query string
	if skipDuplicate {
		if driver == DriverPostgres {
			query = fmt.Sprintf("insert into %s (%s) values %s on conflict do nothing",
				tableName,
				strings.Join(dbColumns, ", "),
				strings.Join(placeholders, ", "),
			)
		} else if driver == DriverSqlite3 {
			query = fmt.Sprintf("insert or ignore into %s (%s) values %s",
				tableName,
				strings.Join(dbColumns, ", "),
				strings.Join(placeholders, ", "),
			)
		} else if driver == DriverOracle || driver == DriverMysql {
			var qKey []string
			for _, i2 := range pkey {
				key := i2 + " = " + i2
				qKey = append(qKey, key)
			}
			query = fmt.Sprintf("insert into %s (%s) values %s on duplicate key update %s",
				tableName,
				strings.Join(dbColumns, ", "),
				strings.Join(placeholders, ", "),
				strings.Join(qKey, ", "),
			)
		} else {
			return 0, fmt.Errorf("only support skip duplicate on mysql and postgresql, current vendor is %s", driver)
		}
	} else {
		if driver != DriverOracle {
			query = fmt.Sprintf(fmt.Sprintf("insert into %s (%s) values %s",
				tableName,
				strings.Join(dbColumns, ","),
				strings.Join(placeholders, ","),
			))
		} else {
			all := make([]string, 0)
			colNames := "(" + strings.Join(dbColumns, ",") + ")"
			for _, s0 := range placeholders {
				s1 := fmt.Sprintf(" into %s %s values %s ", tableName, colNames, s0)
				all = append(all, s1)
			}
			query = fmt.Sprintf(" insert all %s select * from dual", strings.Join(all, " "))
		}
	}
	mainScope.Query = query

	x, err := db.ExecContext(ctx, mainScope.Query, mainScope.Values...)
	if err != nil {
		return -1, err
	}
	return x.RowsAffected()
}

func InsertInTransaction(ctx context.Context, db *sql.DB, tableName string, objects []interface{}, skipDuplicate bool, buildParam func(i int) string, excludeColumns ...string) (int64, error) {
	if len(objects) == 0 {
		return 0, nil
	}
	driver := GetDriver(db)
	firstAttrs, _, _, err := ExtractMapValue(objects[0], &excludeColumns, true)
	if err != nil {
		return 0, err
	}

	attrSize := len(firstAttrs)
	modelType := reflect.TypeOf(objects[0])
	pkey := FindIdFields(modelType)
	// Replace with database column name
	dbColumns := make([]string, 0, attrSize)
	for _, key := range SortedKeys(firstAttrs) {
		dbColumns = append(dbColumns, QuoteColumnName(key))
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	var start int
	for _, obj := range objects {
		// Scope to eventually run SQL
		mainScope := BatchStatement{}
		// Store placeholders for embedding variables
		placeholders := make([]string, 0, attrSize)
		objAttrs, _, _, err := ExtractMapValue(obj, &excludeColumns, true)
		if err != nil {
			return 0, err
		}

		// If object sizes are different, SQL statement loses consistency
		if len(objAttrs) != attrSize {
			return 0, errors.New("attribute sizes are inconsistent")
		}

		scope := BatchStatement{}

		// Append variables
		variables := make([]string, 0, attrSize)
		for _, key := range SortedKeys(objAttrs) {
			scope.Values = append(scope.Values, objAttrs[key])
			variables = append(variables, BuildParametersFrom(start, 1, buildParam))
			start++
		}

		valueQuery := "(" + strings.Join(variables, ", ") + ")"
		placeholders = append(placeholders, valueQuery)

		// Also append variables to mainScope
		mainScope.Values = append(mainScope.Values, scope.Values...)

		var query string
		if skipDuplicate {
			if driver == DriverPostgres {
				query = fmt.Sprintf("insert into %s (%s) values %s on conflict do nothing",
					tableName,
					strings.Join(dbColumns, ", "),
					strings.Join(placeholders, ", "),
				)
			} else if driver == DriverSqlite3 {
				query = fmt.Sprintf("insert or ignore into %s (%s) values %s",
					tableName,
					strings.Join(dbColumns, ", "),
					strings.Join(placeholders, ", "),
				)
			} else if driver == DriverOracle || driver == DriverMysql {
				var qKey []string
				for _, i2 := range pkey {
					key := i2 + " = " + i2
					qKey = append(qKey, key)
				}
				query = fmt.Sprintf("insert into %s (%s) values %s on duplicate key update %s",
					tableName,
					strings.Join(dbColumns, ", "),
					strings.Join(placeholders, ", "),
					strings.Join(qKey, ", "),
				)
			} else {
				return 0, fmt.Errorf("only support skip duplicate on mysql and postgresql, current vendor is %s", driver)
			}
		} else {
			query = fmt.Sprintf("insert into %s (%s) values %s",
				tableName,
				strings.Join(dbColumns, ", "),
				strings.Join(placeholders, ", "),
			)
		}
		query = ReplaceParameters(driver, query, len(mainScope.Values))
		mainScope.Query = query

		_, execErr := tx.ExecContext(ctx, mainScope.Query, mainScope.Values...)
		if execErr != nil {
			_ = tx.Rollback()
			return 0, execErr
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	count := len(objects)
	return int64(count), err
}

func InterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.Indirect(reflect.ValueOf(slice))
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("InterfaceSlice() given a non-slice type")
	}
	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret, nil
}

func UpdateInTransaction(ctx context.Context, db *sql.DB, tableName string, objects []interface{}, options ...func(i int) string) (int64, error) {
	var placeholder []string
	driver := GetDriver(db)
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	var query []string
	var value [][]interface{}
	if len(objects) == 0 {
		return 0, nil
	}
	valueDefault := objects[0]
	statement := newStatement(valueDefault, placeholder...)
	columns := make([]string, 0) // columns set value for update
	for _, key := range SortedKeys(statement.Attributes) {
		columns = append(columns, QuoteByDriver(key, driver))
	}
	for _, obj := range objects {
		scope := newStatement(obj, placeholder...)
		// Append variables set column
		for _, key := range SortedKeys(scope.Attributes) {
			scope.Values = append(scope.Values, scope.Attributes[key])
		}
		// Append variables where
		for _, key := range scope.Keys {
			scope.Values = append(scope.Values, scope.AttributeKeys[key])
		}
		// Also append variables to mainScope
		//statement.Values = append(statement.Values, scope.Values...)

		n := len(scope.Columns)
		sets, setVal, err1 := BuildSqlParametersAndValues(scope.Columns, scope.Values, &n, 0, ", ", buildParam)
		if err1 != nil {
			return 0, err1
		}
		numKeys := len(scope.Keys)
		where, whereVal, err2 := BuildSqlParametersAndValues(scope.Keys, scope.Values, &numKeys, n, " and ", buildParam)
		if err2 != nil {
			return 0, err2
		}
		setVal = append(setVal, whereVal...)
		value = append(value, setVal)
		query = append(query, fmt.Sprintf(fmt.Sprintf("update %s set %s where %s",
			tableName,
			sets,
			where,
		)))
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(query); i++ {
		_, execErr := tx.ExecContext(ctx, query[i], value[i]...)
		if execErr != nil {
			_ = tx.Rollback()
			return 0, execErr
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	total := int64(len(query))
	return total, err
}
func PatchInTransaction(ctx context.Context, db *sql.DB, tableName string, objects []map[string]interface{}, idTagJsonNames []string, idColumNames []string, options ...func(i int) string) (int64, error) {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	var query []string
	var values [][]interface{}
	if len(objects) == 0 {
		return 0, nil
	}
	for _, obj := range objects {
		scope := statement()
		valueRow := make([]interface{}, 0)
		// Append variables set column
		for key, _ := range obj {
			if _, ok := Find(idTagJsonNames, key); !ok {
				scope.Columns = append(scope.Columns, key)
				scope.Values = append(scope.Values, obj[key])
			}
		}
		// Append variables where
		for i, key := range idTagJsonNames {
			scope.Values = append(scope.Values, obj[key])
			scope.Keys = append(scope.Keys, idColumNames[i])
		}

		n := len(scope.Columns)
		sets, setVal, err1 := BuildSqlParametersAndValues(scope.Columns, scope.Values, &n, 0, ", ", buildParam)
		if err1 != nil {
			return 0, err1
		}
		valueRow = append(valueRow, setVal...)
		numKeys := len(scope.Keys)
		where, whereVal, err2 := BuildSqlParametersAndValues(scope.Keys, scope.Values, &numKeys, n, " and ", buildParam)
		if err2 != nil {
			return 0, err2
		}
		valueRow = append(valueRow, whereVal...)
		query = append(query, fmt.Sprintf("update %s set %s where %s",
			tableName,
			sets,
			where,
		))
		values = append(values, valueRow)
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(query); i++ {
		_, execErr := tx.ExecContext(ctx, query[i], values[i]...)
		if execErr != nil {
			_ = tx.Rollback()
			return 0, execErr
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	total := int64(len(query))
	return total, err
}
func PatchMaps(ctx context.Context, db *sql.DB, tableName string, objects []map[string]interface{}, idTagJsonNames []string, idColumNames []string, options ...func(i int) string) (int64, error) {
	var buildParam func(i int) string
	if len(options) > 0 && options[0] != nil {
		buildParam = options[0]
	} else {
		buildParam = GetBuild(db)
	}
	var query []string
	var value [][]interface{}
	if len(objects) == 0 {
		return 0, nil
	}
	for _, obj := range objects {
		scope := statement()
		// Append variables set column
		for key, _ := range obj {
			if _, ok := Find(idTagJsonNames, key); !ok {
				scope.Columns = append(scope.Columns, key)
				scope.Values = append(scope.Values, obj[key])
			}
		}
		// Append variables where
		for i, key := range idTagJsonNames {
			scope.Values = append(scope.Values, obj[key])
			scope.Keys = append(scope.Keys, idColumNames[i])
		}

		n := len(scope.Columns)
		sets, setVal, err1 := BuildSqlParametersAndValues(scope.Columns, scope.Values, &n, 0, ", ", buildParam)
		if err1 != nil {
			return 0, err1
		}
		value = append(value, setVal)
		numKeys := len(scope.Keys)
		where, whereVal, err2 := BuildSqlParametersAndValues(scope.Keys, scope.Values, &numKeys, n, " and ", buildParam)
		if err2 != nil {
			return 0, err2
		}
		value = append(value, whereVal)
		query = append(query, fmt.Sprintf("update %s set %s where %s",
			tableName,
			sets,
			where,
		))
	}

	var count int64
	for i := 0; i < len(query); i++ {
		x, execErr := db.ExecContext(ctx, query[i], value[i]...)
		if execErr != nil {
			return 0, execErr
		}
		rowsAffected, _ := x.RowsAffected()
		count += rowsAffected
	}
	return count, nil
}

func GetValueColumn(value interface{}, driver string) (string, error) {
	str := ""
	switch v := value.(type) {
	case int:
		str = strconv.Itoa(v)
	case int64:
		str = strconv.Itoa(int(v))
	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		str = strconv.FormatBool(v)
	case time.Time:
		str = formatStringByDriver(v.Format(formatDateUpdate), driver)
	case *time.Time:
		str = formatStringByDriver(v.Format(formatDateUpdate), driver)
	case string:
		str = formatStringByDriver(v, driver)
	default:
		return "", errors.New("unsupported type")
	}
	return str, nil
}

func formatStringByDriver(v, driver string) string {
	if driver == DriverPostgres {
		return `E'` + EscapeString(v) + `'`
	} else if driver == DriverMssql {
		return `'` + EscapeStringForSelect(v) + `'`
	} else {
		return `'` + EscapeString(v) + `'`
	}
	return v
}

func BuildSqlParametersByColumns(columns []string, values []interface{}, n int, start int, driver string, joinStr string) (string, error) {
	arr := make([]string, n)
	j := start
	for i, _ := range arr {
		columnName := columns[i]
		value, err := GetValueColumn(values[j], driver)
		if err == nil {
			arr[i] = fmt.Sprintf("%s = %s", columnName, value)
		} else {
			return "", err
		}
		j++
	}
	return strings.Join(arr, joinStr), nil
}

func ReplaceParameters(driver string, query string, n int) string {
	if driver == DriverOracle || driver == DriverPostgres || driver == DriverSqlite3 {
		var x string
		if driver == DriverOracle {
			x = ":"
		} else {
			x = "$"
		}
		for i := 0; i < n; i++ {
			count := i + 1
			query = strings.Replace(query, "?", x+fmt.Sprintf("%v", count), 1)
		}
	}
	return query
}
