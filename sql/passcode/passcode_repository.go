package passcode

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	q "github.com/core-go/core/sql"
)

type PasscodeRepository struct {
	db            *sql.DB
	tableName     string
	idName        string
	passcodeName  string
	expiredAtName string
	Tx            string
	Driver        string
	BuildParam    func(i int) string
}
func NewPasscodeAdapter(db *sql.DB, tableName string, options ...string) *PasscodeRepository {
	return NewPasscodeRepositoryWithTx(db, tableName, "", options...)
}
func NewPasscodeAdapterWithTx(db *sql.DB, tableName string, tx string, options ...string) *PasscodeRepository {
	return NewPasscodeRepositoryWithTx(db, tableName, tx, options...)
}
func NewPasscodeRepository(db *sql.DB, tableName string, options ...string) *PasscodeRepository {
	return NewPasscodeRepositoryWithTx(db, tableName, "", options...)
}
func NewPasscodeRepositoryWithTx(db *sql.DB, tableName string, tx string, options ...string) *PasscodeRepository {
	var idName, passcodeName, expiredAtName string
	if len(options) >= 1 && len(options[0]) > 0 {
		expiredAtName = options[0]
	} else {
		expiredAtName = "expiredat"
	}
	if len(options) >= 2 && len(options[1]) > 0 {
		idName = options[1]
	} else {
		idName = "id"
	}
	if len(options) >= 3 && len(options[2]) > 0 {
		passcodeName = options[2]
	} else {
		passcodeName = "passcode"
	}
	driver := q.GetDriver(db)
	buildParam := q.GetBuild(db)
	return &PasscodeRepository{
		db:            db,
		tableName:     strings.ToLower(tableName),
		idName:        strings.ToLower(idName),
		passcodeName:  strings.ToLower(passcodeName),
		expiredAtName: strings.ToLower(expiredAtName),
		Tx:            tx,
		Driver:        driver,
		BuildParam:    buildParam,
	}
}

func (s *PasscodeRepository) Save(ctx context.Context, id string, passcode string, expireAt time.Time) (int64, error) {
	var placeholder []string
	columns := []string{s.idName, s.passcodeName, s.expiredAtName}
	var queryString string
	driver := q.GetDriver(s.db)
	for i := 0; i < 3; i++ {
		placeholder = append(placeholder, s.BuildParam(i+1))
	}
	if driver == q.DriverPostgres {
		setColumns := make([]string, 0)
		for i, key := range columns {
			setColumns = append(setColumns, key+" = "+s.BuildParam(i+4))
		}
		queryString = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON CONFLICT (%s) DO UPDATE SET %s",
			s.tableName,
			strings.Join(columns, ", "),
			"("+strings.Join(placeholder, ", ")+")",
			s.idName,
			strings.Join(setColumns, ", "),
		)
	} else if driver == q.DriverMysql {
		setColumns := make([]string, 0)
		for i, key := range columns {
			setColumns = append(setColumns, key+" = "+s.BuildParam(i+3))
		}

		queryString = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON DUPLICATE KEY UPDATE %s",
			s.tableName,
			strings.Join(columns, ", "),
			"("+strings.Join(placeholder, ", ")+")",
			strings.Join(setColumns, ", "),
		)
	} else if driver == q.DriverOracle {
		var placeholderOracle []string
		for i := 0; i < 3; i++ {
			placeholderOracle = append(placeholderOracle, s.BuildParam(i+4))
		}
		setColumns := make([]string, 0)
		onDupe := s.tableName + "." + s.idName + " = " + "temp." + s.idName
		for _, key := range columns {
			if key == s.idName {
				continue
			}
			setColumns = append(setColumns, key+" = temp."+key)
		}
		queryString = fmt.Sprintf("MERGE INTO %s USING (SELECT %s as %s, %s as %s, %s as %s  FROM dual) temp ON (%s) WHEN MATCHED THEN UPDATE SET %s WHEN NOT MATCHED THEN INSERT (%s) VALUES (%s)",
			s.tableName,
			s.BuildParam(1), s.idName,
			s.BuildParam(2), s.passcodeName,
			s.BuildParam(3), s.expiredAtName,
			onDupe,
			strings.Join(setColumns, ", "),
			strings.Join(columns, ", "),
			strings.Join(placeholderOracle, ", "),
		)
	} else if driver == q.DriverMssql {
		setColumns := make([]string, 0)
		onDupe := s.tableName + "." + s.idName + " = " + "temp." + s.idName
		for _, key := range columns {
			setColumns = append(setColumns, key+" = temp."+key)
		}
		queryString = fmt.Sprintf("MERGE INTO %s USING (VALUES %s) AS temp (%s) ON %s WHEN MATCHED THEN UPDATE SET %s WHEN NOT MATCHED THEN INSERT (%s) VALUES %s;",
			s.tableName,
			strings.Join(placeholder, ", "),
			strings.Join(columns, ", "),
			onDupe,
			strings.Join(setColumns, ", "),
			strings.Join(columns, ", "),
			strings.Join(placeholder, ", "),
		)
	} else if driver == q.DriverSqlite3 {
		setColumns := make([]string, 0)
		for i, key := range columns {
			setColumns = append(setColumns, key+" = "+s.BuildParam(i+3))
		}
		queryString = fmt.Sprintf("insert or replace into %s (%s) values %s",
			s.tableName,
			strings.Join(columns, ", "),
			"("+strings.Join(placeholder, ", ")+")",
		)
	} else {
		return 0, fmt.Errorf("unsupported db vendor, current vendor is %s", driver)
	}
	if len(s.Tx) > 0 {
		txv := ctx.Value(s.Tx)
		if txv != nil {
			tx, ok := txv.(*sql.Tx)
			if ok {
				x0, er0 := tx.ExecContext(ctx, queryString, id, passcode, expireAt, id, passcode, expireAt)
				if er0 != nil {
					return 0, er0
				}
				return x0.RowsAffected()
			}
		}
	}
	x, err := s.db.ExecContext(ctx, queryString, id, passcode, expireAt, id, passcode, expireAt)
	if err != nil {
		return 0, err
	}
	return x.RowsAffected()
}

func (s *PasscodeRepository) Load(ctx context.Context, id string) (string, time.Time, error) {
	driverName := q.GetDriver(s.db)
	arr := make(map[string]interface{})
	strSql := fmt.Sprintf(`SELECT %s, %s FROM `, s.passcodeName, s.expiredAtName) + s.tableName + ` WHERE ` + s.idName + ` = ` + s.BuildParam(1)
	rows, err := s.db.QueryContext(ctx, strSql, id)
	if err != nil {
		return "", time.Now().Add(-24 * time.Hour), err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return "", time.Now().Add(-24 * time.Hour), err
		}

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			arr[colName] = *val
		}
	}

	err2 := rows.Err()
	if err2 != nil {
		return "", time.Now().Add(-24 * time.Hour), err2
	}

	if len(arr) == 0 {
		return "", time.Now().Add(-24 * time.Hour), nil
	}

	var code string
	var expiredAt time.Time
	if driverName == q.DriverPostgres {
		code = arr[s.passcodeName].(string)
	} else if driverName == q.DriverOracle {
		code = arr[strings.ToUpper(s.passcodeName)].(string)
	} else {
		code = string(arr[s.passcodeName].([]byte))
	}
	if driverName == q.DriverOracle {
		expiredAt = arr[strings.ToUpper(s.expiredAtName)].(time.Time)
	} else {
		expiredAt = arr[s.expiredAtName].(time.Time)
	}
	return code, expiredAt, nil
}

func (s *PasscodeRepository) Delete(ctx context.Context, id string) (int64, error) {
	strSQL := `DELETE FROM ` + s.tableName + ` WHERE ` + s.idName + ` =  ` + s.BuildParam(1)
	x, err := s.db.ExecContext(ctx, strSQL, id)
	if err != nil {
		return 0, err
	}
	return x.RowsAffected()
}
