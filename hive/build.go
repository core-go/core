package hive

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	_ "time"
)

func BuildToInsert(table string, model interface{}, options ...*Schema) string {
	return BuildToInsertWithVersion(table, model, -1, false, options...)
}
func BuildToSave(table string, model interface{}, options ...*Schema) string {
	return BuildToInsertWithVersion(table, model, -1, true, options...)
}
func BuildToInsertWithVersion(table string, model interface{}, versionIndex int, orUpdate bool, options ...*Schema) string {
	modelType := reflect.TypeOf(model)
	var cols []*FieldDB
	if len(options) > 0 && options[0] != nil {
		cols = options[0].Columns
	} else {
		m := CreateSchema(modelType)
		cols = m.Columns
	}
	mv := reflect.ValueOf(model)
	if mv.Kind() == reflect.Ptr {
		mv = mv.Elem()
	}
	values := make([]string, 0)
	icols := make([]string, 0)
	for _, fdb := range cols {
		if fdb.Index == versionIndex {
			icols = append(icols, fdb.Column)
			values = append(values, "1")
		} else {
			f := mv.Field(fdb.Index)
			fieldValue := f.Interface()
			isNil := false
			if f.Kind() == reflect.Ptr {
				if reflect.ValueOf(fieldValue).IsNil() {
					isNil = true
				} else {
					fieldValue = reflect.Indirect(reflect.ValueOf(fieldValue)).Interface()
				}
			}
			if fdb.Insert {
				if isNil {
					if orUpdate {
						icols = append(icols, fdb.Column)
						values = append(values, "null")
					}
				} else {
					icols = append(icols, fdb.Column)
					v, ok := GetDBValue(fieldValue, fdb.Scale, fdb.LayoutTime)
					if ok {
						values = append(values, v)
					}
					//TODO error here
				}
			}
		}
	}
	return fmt.Sprintf("insert into %v(%v) values (%v)", table, strings.Join(icols, ","), strings.Join(values, ","))
}
func BuildToUpdate(table string, model interface{}, options ...*Schema) string {
	return BuildToUpdateWithVersion(table, model, -1, options...)
}
func BuildToUpdateWithVersion(table string, model interface{}, versionIndex int, options ...*Schema) string {
	var cols, keys []*FieldDB
	modelType := reflect.TypeOf(model)
	if len(options) > 0 && options[0] != nil {
		m := options[0]
		cols = m.Columns
		keys = m.Keys
	} else {
		m := CreateSchema(modelType)
		cols = m.Columns
		keys = m.Keys
	}
	mv := reflect.ValueOf(model)
	if mv.Kind() == reflect.Ptr {
		mv = mv.Elem()
	}
	values := make([]string, 0)
	where := make([]string, 0)
	vw := ""
	for _, fdb := range cols {
		// fdb2 := schema[col]
		if fdb.Index == versionIndex {
			valueOfModel := reflect.Indirect(reflect.ValueOf(model))
			currentVersion := reflect.Indirect(valueOfModel.Field(versionIndex)).Int()
			nv := currentVersion + 1
			values = append(values, fdb.Column+"="+strconv.FormatInt(nv, 10))
			vw = fdb.Column + "=" + strconv.FormatInt(currentVersion, 10)
		} else if !fdb.Key && fdb.Update {
			//f := reflect.Indirect(reflect.ValueOf(model))
			f := mv.Field(fdb.Index)
			fieldValue := f.Interface()
			isNil := false
			if f.Kind() == reflect.Ptr {
				if reflect.ValueOf(fieldValue).IsNil() {
					isNil = true
				} else {
					fieldValue = reflect.Indirect(reflect.ValueOf(fieldValue)).Interface()
				}
			}
			if isNil {
				values = append(values, fdb.Column+"=null")
			} else {
				v, ok := GetDBValue(fieldValue, fdb.Scale, fdb.LayoutTime)
				if ok {
					values = append(values, fdb.Column+"="+v)
				}
			}
		}
	}
	for _, fdb := range keys {
		// fdb2 := schema[col]
		f := mv.Field(fdb.Index)
		fieldValue := f.Interface()
		if f.Kind() == reflect.Ptr {
			if !reflect.ValueOf(fieldValue).IsNil() {
				fieldValue = reflect.Indirect(reflect.ValueOf(fieldValue)).Interface()
			}
		}
		v, ok := GetDBValue(fieldValue, fdb.Scale, fdb.LayoutTime)
		if ok {
			where = append(where, fdb.Column+"="+v)
		}
	}
	if len(vw) > 0 {
		where = append(where, vw)
	}
	query := fmt.Sprintf("update %v set %v where %v", table, strings.Join(values, ","), strings.Join(where, " and "))
	return query
}
func BuildToDelete(table string, ids map[string]interface{}) string {
	var queryArr []string
	for col, value := range ids {
		v, ok := GetDBValue(value, 0, "")
		if ok {
			queryArr = append(queryArr, col+"="+v)
		}
	}
	q := strings.Join(queryArr, " and ")
	return fmt.Sprintf("delete from %v where %v", table, q)
}
