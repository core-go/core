package sql

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

type batchStatement struct {
	Query         string
	Values        []interface{}
	Keys          []string
	Columns       []string
	Attributes    map[string]interface{}
	AttributeKeys map[string]interface{}
}

type StringService struct {
	DB         *sql.DB
	Table      string
	Field      string
	Sql        string
	Driver     string
	BuildParam func(i int) string
}

func NewStringService(db *sql.DB, table string, field string, options ...func(i int) string) *StringService {
	driver := getDriver(db)
	var b func(i int) string
	if len(options) > 0 {
		b = options[0]
	} else {
		b = getBuild(db)
	}
	var sql string
	if driver == driverPostgres {
		sql = fmt.Sprintf("select %s from %s where %s ilike %s", field, table, field, b(1)) + " fetch next %d rows only"
	} else if driver == driverOracle {
		sql = fmt.Sprintf("select %s from %s where %s like %s", field, table, field, b(1)) + " fetch next %d rows only"
	} else {
		sql = fmt.Sprintf("select %s from %s where %s like %s", field, table, field, b(1)) + " limit %d"
	}
	return &StringService{DB: db, Table: table, Field: field, Sql: sql, Driver: driver, BuildParam: b}
}

func (s *StringService) Load(ctx context.Context, key string, max int64) ([]string, error) {
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

func (s *StringService) Save(ctx context.Context, values []string) (int64, error) {
	mainScope := batchStatement{}
	driver := s.Driver
	for _, e := range values {
		mainScope.Values = append(mainScope.Values, e)
	}
	query := ""
	holders := buildPlaceHolders(len(mainScope.Values), s.BuildParam)
	if driver == driverPostgres {
		query = fmt.Sprintf("insert into %s (%s) values %s on conflict do nothing",
			s.Table,
			s.Field,
			holders,
		)
	} else if driver == driverSqlite3 {
		query = fmt.Sprintf("insert or ignore into %s (%s) values %s",
			s.Table,
			s.Field,
			holders,
		)
	} else if driver == driverMysql {
		qKey := s.Field + " = " + s.Field
		query = fmt.Sprintf("insert into %s (%s) values %s on duplicate key update %s",
			s.Table,
			s.Field,
			holders,
			qKey,
		)
	} else if driver == driverOracle || driver == driverMssql {
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

func (s *StringService) Delete(ctx context.Context, values []string) (int64, error) {
	var arrValue []string
	le := len(values)
	for i := 1; i <= le; i++ {
		param := buildParam(i)
		arrValue = append(arrValue, param)
	}
	query := `delete from ` + s.Table + ` where ` + s.Field + ` in (` + strings.Join(arrValue, ",") + `)`
	x, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return x.RowsAffected()
}
func buildPlaceHolders(n int, buildParam func(int) string) string {
	ss := make([]string, 0)
	for i := 1; i <= n; i++ {
		s := buildParam(i)
		ss = append(ss, s)
	}
	return strings.Join(ss, ",")
}
