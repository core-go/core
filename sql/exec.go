package sql

import (
	"context"
	"database/sql"
	"strings"
)

func ExecuteStatements(ctx context.Context, tx *sql.Tx, commit bool, stmts ...Statement) (int64, error) {
	if stmts == nil || len(stmts) == 0 {
		return 0, nil
	}
	var count int64
	count = 0
	for _, stmt := range stmts {
		r2, er3 := tx.ExecContext(ctx, stmt.Query, stmt.Params...)
		if er3 != nil {
			er4 := tx.Rollback()
			if er4 != nil {
				return count, er4
			}
			return count, er3
		}
		a2, er5 := r2.RowsAffected()
		if er5 != nil {
			tx.Rollback()
			return count, er5
		}
		count = count + a2
	}
	if commit {
		er6 := tx.Commit()
		return count, er6
	} else {
		return count, nil
	}
}
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
func ExecuteWithBatchSize(ctx context.Context, db *sql.DB, size int, stmts ...Statement) (int64, error) {
	l := len(stmts)
	if len(stmts) == 0 {
		return -1, nil
	}
	min := Min(l, size)
	max := Max(l, size)
	if min == l {
		return ExecuteAll(ctx, db, stmts...)
	} else {
		i := 0
		k := 0
		s := make([]Statement, 0)
		for {
			for j := 0; j < min; j++ {
				s = append(s, stmts[k])
				k = k + 1
				if k >= l {
					break
				}
			}
			_, err := ExecuteAll(ctx, db, s...)
			if err != nil {
				return int64(i), err
			}
			i += min
			if i >= max {
				break
			}
		}
	}
	return int64(l), nil
}
func ExecuteAll(ctx context.Context, db *sql.DB, stmts ...Statement) (int64, error) {
	if stmts == nil || len(stmts) == 0 {
		return 0, nil
	}
	tx, er1 := db.Begin()
	if er1 != nil {
		return 0, er1
	}
	var count int64
	count = 0
	for _, stmt := range stmts {
		r2, er3 := tx.ExecContext(ctx, stmt.Query, stmt.Params...)
		if er3 != nil {
			er4 := tx.Rollback()
			if er4 != nil {
				return count, er4
			}
			return count, er3
		}
		a2, er5 := r2.RowsAffected()
		if er5 != nil {
			tx.Rollback()
			return count, er5
		}
		count = count + a2
	}
	er6 := tx.Commit()
	return count, er6
}
func ExecuteBatch(ctx context.Context, db *sql.DB, sts []Statement, firstRowSuccess bool, countAll bool) (int64, error) {
	if sts == nil || len(sts) == 0 {
		return 0, nil
	}
	driver := GetDriver(db)
	tx, er0 := db.Begin()
	if er0 != nil {
		return 0, er0
	}
	result, er1 := tx.ExecContext(ctx, sts[0].Query, sts[0].Params...)
	if er1 != nil {
		_ = tx.Rollback()
		str := er1.Error()
		if driver == DriverPostgres && strings.Contains(str, "pq: duplicate key value violates unique constraint") {
			return 0, nil
		} else if driver == DriverMysql && strings.Contains(str, "Error 1062: Duplicate entry") {
			return 0, nil //mysql Error 1062: Duplicate entry 'a-1' for key 'PRIMARY'
		} else if driver == DriverOracle && strings.Contains(str, "ORA-00001: unique constraint") {
			return 0, nil //mysql Error 1062: Duplicate entry 'a-1' for key 'PRIMARY'
		} else if driver == DriverMssql && strings.Contains(str, "Violation of PRIMARY KEY constraint") {
			return 0, nil //Violation of PRIMARY KEY constraint 'PK_aa'. Cannot insert duplicate key in object 'dbo.aa'. The duplicate key value is (b, 2).
		} else if driver == DriverSqlite3 && strings.Contains(str, "UNIQUE constraint failed") {
			return 0, nil
		} else {
			return 0, er1
		}
	}
	rowAffected, er2 := result.RowsAffected()
	if er2 != nil {
		tx.Rollback()
		return 0, er2
	}
	if firstRowSuccess {
		if rowAffected == 0 {
			return 0, nil
		}
	}
	count := rowAffected
	for i := 1; i < len(sts); i++ {
		r2, er3 := tx.ExecContext(ctx, sts[i].Query, sts[i].Params...)
		if er3 != nil {
			er4 := tx.Rollback()
			if er4 != nil {
				return count, er4
			}
			return count, er3
		}
		a2, er5 := r2.RowsAffected()
		if er5 != nil {
			tx.Rollback()
			return count, er5
		}
		count = count + a2
	}
	er6 := tx.Commit()
	if er6 != nil {
		return count, er6
	}
	if countAll {
		return count, nil
	}
	return 1, nil
}

type Statements interface {
	Exec(ctx context.Context, db *sql.DB) (int64, error)
	Add(sql string, args []interface{}) Statements
	Clear() Statements
}

func NewDefaultStatements(successFirst bool) *DefaultStatements {
	stmts := make([]Statement, 0)
	s := &DefaultStatements{Statements: stmts, SuccessFirst: successFirst}
	return s
}
func NewStatements(successFirst bool) Statements {
	return NewDefaultStatements(successFirst)
}

type DefaultStatements struct {
	Statements   []Statement
	SuccessFirst bool
}

func (s *DefaultStatements) Exec(ctx context.Context, db *sql.DB) (int64, error) {
	if s.SuccessFirst {
		return ExecuteBatch(ctx, db, s.Statements, true, false)
	} else {
		return ExecuteAll(ctx, db, s.Statements...)
	}
}
func (s *DefaultStatements) Add(sql string, args []interface{}) Statements {
	var stm = Statement{Query: sql, Params: args}
	s.Statements = append(s.Statements, stm)
	return s
}
func (s *DefaultStatements) Clear() Statements {
	s.Statements = s.Statements[:0]
	return s
}
