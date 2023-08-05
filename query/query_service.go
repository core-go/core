package query

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	DriverPostgres   = "postgres"
	DriverMysql      = "mysql"
	DriverMssql      = "mssql"
	DriverOracle     = "oracle"
	DriverSqlite3    = "sqlite3"
	DriverNotSupport = "no support"
)

type QueryService struct {
	DB         *sql.DB
	Table      string
	Field      string
	Sql        string
	Driver     string
	BuildParam func(i int) string
}
func NewStringService(db *sql.DB, table string, field string, options ...func(i int) string) *QueryService {
	return NewQueryService(db, table, field, options...)
}
func NewQueryService(db *sql.DB, table string, field string, options ...func(i int) string) *QueryService {
	driver := GetDriver(db)
	var b func(i int) string
	if len(options) > 0 {
		b = options[0]
	} else {
		b = GetBuild(db)
	}
	var sql string
	if driver == DriverPostgres {
		sql = fmt.Sprintf("select %s from %s where %s ilike %s", field, table, field, b(1)) + " fetch next %d rows only"
	} else if driver == DriverOracle {
		sql = fmt.Sprintf("select %s from %s where %s like %s", field, table, field, b(1)) + " fetch next %d rows only"
	} else {
		sql = fmt.Sprintf("select %s from %s where %s like %s", field, table, field, b(1)) + " limit %d"
	}
	return &QueryService{DB: db, Table: table, Field: field, Sql: sql, Driver: driver, BuildParam: b}
}

func (s *QueryService) Load(ctx context.Context, key string, max int64) ([]string, error) {
	re := regexp.MustCompile(`\%|\?`)
	key = re.ReplaceAllString(key, "")
	key = key + "%"
	vs := make([]string, 0)
	sql := fmt.Sprintf(s.Sql, max)
	rows, er1 := s.DB.QueryContext(ctx, sql, key)
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

func (s *QueryService) Save(ctx context.Context, values []string) (int64, error) {
	driver := s.Driver
	l := len(values)
	if l == 0 {
		return 0, nil
	}
	if driver == DriverPostgres || driver == DriverMysql {
		ps := make([]string, 0)
		p := make([]interface{}, 0)
		for _, str := range values {
			p = append(p, str)
		}
		if driver == DriverPostgres {
			for i := 1; i <= l; i++ {
				ps = append(ps, "(" + BuildDollarParam(i) + ")")
			}
		} else {
			for i := 1; i <= l; i++ {
				ps = append(ps, "(?)")
			}
		}
		var query string
		if driver == DriverPostgres {
			query = fmt.Sprintf("insert into %s (%s) values %s on conflict do nothing", s.Table, s.Field, strings.Join(ps, ","))
		} else {
			query = fmt.Sprintf("insert ignore into %s (%s) values %s", s.Table, s.Field, strings.Join(ps, ","))
		}
		tx, err := s.DB.Begin()
		if err != nil {
			return -1, err
		}
		res, err := tx.ExecContext(ctx, query, p...)
		if err != nil {
			er := tx.Rollback()
			if er != nil {
				return -1, er
			}
			return -1, err
		}
		err = tx.Commit()
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	} else if driver == DriverSqlite3 {
		tx, err := s.DB.Begin()
		if err != nil {
			return -1, err
		}
		var c int64
		c = 0
		for _, e := range values {
			query := fmt.Sprintf("insert or ignore into %s (%s) values (?)", s.Table, s.Field)
			res, err := tx.ExecContext(ctx, query, e)
			if err != nil {
				er := tx.Rollback()
				if er != nil {
					return -1, er
				}
				return -1, err
			}
			a, err := res.RowsAffected()
			if err != nil {
				return -1, err
			}
			c = c + a
		}
		err = tx.Commit()
		if err != nil {
			return -1, err
		}
		return c, nil
	} else {
		mainScope := BatchStatement{}
		for _, e := range values {
			mainScope.Values = append(mainScope.Values, e)
		}
		query := ""
		holders := BuildPlaceHolders(len(mainScope.Values), s.BuildParam)
		if driver == DriverOracle || driver == DriverMssql {
			onDupe := s.Table + "." + s.Field + " = " + "temp." + s.Field
			value := "temp." + s.Field
			query = fmt.Sprintf("merge into %s using (values %s) as temp (%s) on %s when not matched then insert (%s) values (%s);",
				s.Table,
				holders,
				s.Field,
				onDupe,
				s.Field,
				value,
			)
		} else {
			return 0, fmt.Errorf("unsupported db vendor, current vendor is %s", driver)
		}
		mainScope.Query = query
		x, err := s.DB.ExecContext(ctx, mainScope.Query, mainScope.Values...)
		if err != nil {
			return 0, err
		}
		return x.RowsAffected()
	}
}

func (s *QueryService) Delete(ctx context.Context, values []string) (int64, error) {
	var arrValue []string
	le := len(values)
	buildParam := GetBuild(s.DB)
	p := make([]interface{}, 0)
	for _, str := range values {
		p = append(p, str)
	}
	for i := 1; i <= le; i++ {
		param := buildParam(i)
		arrValue = append(arrValue, param)
	}
	query := `delete from ` + s.Table + ` where ` + s.Field + ` in (` + strings.Join(arrValue, ",") + `)`
	x, err := s.DB.ExecContext(ctx, query, p...)
	if err != nil {
		return 0, err
	}
	return x.RowsAffected()
}
func BuildPlaceHolders(n int, buildParam func(int) string) string {
	ss := make([]string, 0)
	for i := 1; i <= n; i++ {
		s := buildParam(i)
		ss = append(ss, s)
	}
	return strings.Join(ss, ",")
}
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
type BatchStatement struct {
	Query         string
	Values        []interface{}
	Keys          []string
	Columns       []string
	Attributes    map[string]interface{}
	AttributeKeys map[string]interface{}
}
