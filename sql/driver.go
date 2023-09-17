package sql

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	DriverPostgres   = "postgres"
	DriverMysql      = "mysql"
	DriverMssql      = "mssql"
	DriverOracle     = "oracle"
	DriverSqlite3    = "sqlite3"
	DriverNotSupport = "no support"
)

func GetDriver(db *sql.DB) string {
	if db == nil {
		return DriverNotSupport
	}
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return DriverPostgres
	case "*godror.drv":
		return DriverOracle
	case "*mysql.MySQLDriver":
		return DriverMysql
	case "*mssql.Driver":
		return DriverMssql
	case "*sqlite3.SQLiteDriver":
		return DriverSqlite3
	default:
		return DriverNotSupport
	}
}
type BuildParamFn func(i int) string
func BuildParam(i int) string {
	return "?"
}
func BuildOracleParam(i int) string {
	return ":" + strconv.Itoa(i)
}
func BuildMsSqlParam(i int) string {
	return "@p" + strconv.Itoa(i)
}
func BuildDollarParam(i int) string {
	return "$" + strconv.Itoa(i)
}
func GetBuildByDriver(driver string) func(i int) string {
	switch driver {
	case DriverPostgres:
		return BuildDollarParam
	case DriverOracle:
		return BuildOracleParam
	case DriverMssql:
		return BuildMsSqlParam
	default:
		return BuildParam
	}
}
func GetBuild(db *sql.DB) func(i int) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return BuildDollarParam
	case "*godror.drv":
		return BuildOracleParam
	case "*mssql.Driver":
		return BuildMsSqlParam
	default:
		return BuildParam
	}
}
func Detect(s string) string {
	if strings.Index(s, "sqlserver:") == 0 {
		return "mssql"
	} else if strings.Index(s, "postgres:") == 0 {
		return "postgres"
	} else {
		if strings.Index(s, "user=") >= 0 && strings.Index(s, "password=") >= 0 {
			if strings.Index(s, "dbname=") >= 0 || strings.Index(s, "host=") >= 0 || strings.Index(s, "port=") >= 0 {
				return "postgres"
			} else {
				return "godror"
			}
		} else {
			_, err := filepath.Abs(s)
			if (strings.Index(s, "@tcp(") >= 0 || strings.Index(s, "charset=") > 0 || strings.Index(s, "parseTime=") > 0 || strings.Index(s, "loc=") > 0 || strings.Index(s, "@") >= 0 || strings.Index(s, ":") >= 0) && err != nil {
				return "mysql"
			} else {
				return "sqlite3"
			}
		}
	}
}
func OpenByConfig(c Config) (*sql.DB, error) {
	if c.Mock {
		return nil, nil
	}
	if c.Retry.Retry1 <= 0 {
		return open(c)
	} else {
		durations := DurationsFromValue(c.Retry, "Retry", 9)
		return Open(c, durations...)
	}
}
func open(c Config) (*sql.DB, error) {
	dsn := c.DataSourceName
	if len(dsn) > 0 {
		if len(strings.TrimSpace(c.Driver)) == 0 {
			c.Driver = Detect(dsn)
		}
	} else {
		dsn = BuildDataSourceName(c)
	}
	db, err := sql.Open(c.Driver, dsn)
	if err != nil {
		return db, err
	}
	if c.ConnMaxLifetime != nil {
		db.SetConnMaxLifetime(*c.ConnMaxLifetime)
	}
	if c.MaxIdleConns > 0 {
		db.SetMaxIdleConns(c.MaxIdleConns)
	}
	if c.MaxOpenConns > 0 {
		db.SetMaxOpenConns(c.MaxOpenConns)
	}
	return db, err
}
func Open(c Config, retries ...time.Duration) (*sql.DB, error) {
	if c.Mock {
		return nil, nil
	}
	if len(retries) == 0 {
		return open(c)
	} else {
		db, er1 := open(c)
		if er1 == nil {
			return db, er1
		}
		i := 0
		err := Retry(retries, func() (err error) {
			i = i + 1
			db2, er2 := open(c)
			if er2 == nil {
				db = db2
			}
			return er2
		})
		if err != nil {
			log.Printf("Cannot conect to database: %s.", err.Error())
		}
		return db, err
	}
}
func BuildDataSourceName(c Config) string {
	if c.Driver == "postgres" {
		uri := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=disable", c.User, c.Database, c.Password, c.Host, c.Port)
		return uri
	} else if c.Driver == "mysql" {
		uri := ""
		if c.MultiStatements {
			uri = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=True", c.User, c.Password, c.Host, c.Port, c.Database)
			return uri
		}
		uri = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
		return uri
	} else if c.Driver == "mssql" { // mssql
		uri := fmt.Sprintf("sqlserver://%s:%s@%s:%d?Database=%s", c.User, c.Password, c.Host, c.Port, c.Database)
		return uri
	} else if c.Driver == "godror" || c.Driver == "oracle" {
		return fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s:%d/%s\"", c.User, c.Password, c.Host, c.Port, c.Database)
	} else { //sqlite
		return c.Host // return sql.Open("sqlite3", c.Host)
	}
}
