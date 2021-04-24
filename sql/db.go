package sql

import (
	"database/sql"
	"reflect"
	"strconv"
)

const (
	driverPostgres   = "postgres"
	driverMysql      = "mysql"
	driverMssql      = "mssql"
	driverOracle     = "oracle"
	driverSqlite3    = "sqlite3"
	driverNotSupport = "no support"
)

func buildParam(i int) string {
	return "?"
}
func buildOracleParam(i int) string {
	return ":val" + strconv.Itoa(i)
}
func buildMsSqlParam(i int) string {
	return "@p" + strconv.Itoa(i)
}
func buildDollarParam(i int) string {
	return "$" + strconv.Itoa(i)
}
func getBuild(db *sql.DB) func(i int) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return buildDollarParam
	case "*godror.drv":
		return buildOracleParam
	case "*mssql.Driver":
		return buildMsSqlParam
	default:
		return buildParam
	}
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
