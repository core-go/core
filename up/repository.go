package upload

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type StorageRepository interface {
	Load(ctx context.Context, id string) (*UploadModel, error)
	Update(ctx context.Context, item UploadModel) (int64, error)
}

func NewRepository(DB *sql.DB,
	Table string,
	column  string,
	columns UploadFieldColumn, toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}) *SqlRepository {
	utype := reflect.TypeOf(UploadModel{})
	fieldIndex, _ := GetColumnIndexes(utype)
	return &SqlRepository{DB: DB, Table: Table, Column: column, Columns: &columns, toArray: toArray, utype: utype, fieldIndex: fieldIndex}
}

type UploadFieldColumn struct {
	Cover   *string
	Image   *string
	Gallery *string
	Id      string
}

type SqlRepository struct {
	DB      *sql.DB
	Table   string
	Column  string
	Columns *UploadFieldColumn
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	utype reflect.Type
	fieldIndex map[string]int
	// CoverUrl string
}

func (s *SqlRepository) Load(ctx context.Context, id string) (*UploadModel, error) {
	var models []UploadModel
	query := BuildFindById(s.Table, s.Columns.Id, *s.Columns)
	err := Query(ctx, s.DB, s.fieldIndex, &models, s.toArray, query, id)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	} else {
		return &models[0], nil
	}
}

func (s *SqlRepository) Update(ctx context.Context, item UploadModel) (int64, error) {
	// query := fmt.Sprintf("update %s set %s = $1 where %s =$2", s.Table, *s.Columns.Cover, s.IdCol)
	query, value := BuildUpdate(s.Table, *s.Columns, item, s.toArray)
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	res, err := stmt.ExecContext(ctx, value...)

	row, er2 := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if row < 0 {
		return -1, er2
	}
	return row, er2
}

func BuildFindById(table string, id string, columns UploadFieldColumn) string {
	var where = ""
	var selectQuery []string
	where = fmt.Sprintf("where %s = $1", id)
	selectQuery = append(selectQuery, id)
	if columns.Image != nil && len(*columns.Image) > 0 {
		selectQuery = append(selectQuery, *columns.Image + " as imageurl")
	}
	if columns.Cover != nil && len(*columns.Cover) > 0 {
		selectQuery = append(selectQuery, *columns.Cover + " as coverurl")
	}
	if columns.Gallery != nil && len(*columns.Gallery) > 0 {
		selectQuery = append(selectQuery, *columns.Gallery + " as gallery")
	}
	selectFields := strings.Join(selectQuery, ",")
	return fmt.Sprintf("select %s from %v %v", selectFields, table, where)
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

func BuildUpdate(table string, columns UploadFieldColumn, item UploadModel, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (string, []interface{}) {
	var setQuery = []string{}
	var value []interface{}
	var index = 1
	if *columns.Image != "" && item.ImageURL != nil {
		setData := fmt.Sprintf("%s = $%d", *columns.Image, index)
		value = append(value, item.ImageURL)
		setQuery = append(setQuery, setData)
		index++
	}
	if *columns.Cover != "" && item.CoverURL != nil {
		setData := fmt.Sprintf("%s = $%d", *columns.Cover, index)
		value = append(value, item.CoverURL)
		setQuery = append(setQuery, setData)
		index++
	}
	if *columns.Gallery != "" && item.Gallery != nil {
		var interfaceSlice []interface{} = make([]interface{}, len(item.Gallery))
		for i, d := range item.Gallery {
			interfaceSlice[i] = d
		}
		setData := fmt.Sprintf("%s = $%d", *columns.Gallery, index)

		value = append(value, toArray(item.Gallery))
		setQuery = append(setQuery, setData)
		index++
	}
	where := fmt.Sprintf("where %s = $%d", columns.Id, index)
	value = append(value, item.Id)
	sets := "set " + strings.Join(setQuery, ",")
	query := fmt.Sprintf("update %s %s %s", table, sets, where)
	return query, value
}

func Query(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, results interface{}, toArray func(interface{}) interface {
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
