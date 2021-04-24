package id

import (
	"context"
	"database/sql"
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
	driverNotSupport = "no support"
)

type FieldLoader struct {
	DB        *sql.DB
	TableName string
	Name      string
	Driver    string
}

func NewFieldLoader(db *sql.DB, tableName string, name string) *FieldLoader {
	driver := getDriver(db)
	return &FieldLoader{
		DB:        db,
		TableName: tableName,
		Name:      name,
		Driver:    driver,
	}
}

func (l *FieldLoader) Values(ctx context.Context, ids []string) ([]string, error) {
	ss := make([]string, 0)
	if ids == nil || len(ids) == 0 {
		return ss, nil
	}
	vs := make([]string, 0)
	params := make([]interface{}, 0)
	if l.Driver == driverPostgres {
		for i, s := range ids {
			ss = append(ss, "$"+strconv.Itoa(i+1))
			params = append(params, s)
		}
	} else if l.Driver == driverOracle {
		for i, s := range ids {
			ss = append(ss, ":val"+strconv.Itoa(i+1))
			params = append(params, s)
		}
	} else {
		for _, s := range ids {
			ss = append(ss, "?")
			params = append(params, s)
		}
	}
	sql := fmt.Sprintf("select distinct %s from %s where %s in (%s)", l.Name, l.TableName, l.Name, strings.Join(ss, ","))
	rows, er1 := l.DB.Query(sql, params...)
	if er1 != nil {
		return vs, er1
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if er2 := rows.Scan(&id); er2 == nil {
			vs = append(vs, id)
		} else {
			return vs, er2
		}
	}
	return vs, nil
}
func getDriver(db *sql.DB) string {
	if db == nil {
		return driverNotSupport
	}
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return driverPostgres
	case "*mysql.MySQLDriver":
		return driverMysql
	case "*mssql.Driver":
		return driverMssql
	case "*godror.drv":
		return driverOracle
	default:
		return driverNotSupport
	}
}
