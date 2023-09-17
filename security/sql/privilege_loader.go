package sql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
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

type PrivilegeLoader struct {
	DB    *sql.DB
	Query string
}

func NewPrivilegeLoader(db *sql.DB, query string, options ...bool) *PrivilegeLoader {
	var handleDriver bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	if handleDriver {
		driver := getDriver(db)
		query = replaceQueryArgs(driver, query)
	}
	return &PrivilegeLoader{DB: db, Query: query}
}

func (l PrivilegeLoader) Privilege(ctx context.Context, userId string, privilegeId string) int32 {
	var permissions int32 = 0
	rows, er0 := l.DB.QueryContext(ctx, l.Query, userId, privilegeId)
	if er0 != nil {
		return actionNone
	}
	defer rows.Close()
	exist := false
	for rows.Next() {
		exist = true
		var action int32
		er1 := rows.Scan(&action)
		if er1 != nil {
			return actionNone
		}
		permissions = permissions | action
	}
	if !exist {
		return actionNone
	}
	if permissions == actionNone {
		return actionAll
	}
	return permissions
}

func replaceQueryArgs(driver string, query string) string {
	if driver == driverOracle || driver == driverPostgres || driver == driverMssql {
		var x string
		if driver == driverOracle {
			x = ":val"
		} else if driver == driverPostgres {
			x = "$"
		} else if driver == driverMssql {
			x = "@p"
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
