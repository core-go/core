package sequence

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type SequenceAdapter struct {
	DB         *sql.DB
	Tables     string
	Table      string
	Sequence   string
	BuildParam func(i int) string
}
func NewSequenceRepository(db *sql.DB, buildParam func(int) string, options ...string) *SequenceAdapter {
	return NewSequenceAdapter(db, buildParam, options...)
}
func NewSequenceAdapter(db *sql.DB, buildParam func(int) string, options ...string) *SequenceAdapter {
	var tables, table, sequence string
	if len(options) > 0 && len(options[0]) > 0 {
		tables = options[0]
	} else {
		tables = "sequences"
	}
	if len(options) > 1 && len(options[1]) > 0 {
		table = options[1]
	} else {
		table = "table"
	}
	if len(options) > 2 && len(options[2]) > 0 {
		sequence = options[2]
	} else {
		sequence = "sequence"
	}
	return &SequenceAdapter{
		DB:         db,
		Tables:     strings.ToLower(tables),
		Table:      strings.ToLower(table),
		Sequence:   strings.ToLower(sequence),
		BuildParam: buildParam,
	}
}
func (s *SequenceAdapter) Next(ctx context.Context, seqName string) (int64, error) {
	seq, err := s.next(ctx, seqName)
	if err != nil {
		return seq, err
	}
	for {
		if seq == -2 {
			seq, err := s.next(ctx, seqName)
			if err != nil {
				return seq, err
			}
		} else {
			return seq, nil
		}
	}
}
func (s *SequenceAdapter) next(ctx context.Context, seqName string) (int64, error) {
	query := fmt.Sprintf(`select %s from %s where %s = %s`, s.Sequence, s.Tables, s.Table, s.BuildParam(1))
	rows, err := s.DB.QueryContext(ctx, query, seqName)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		var seq int64
		if err := rows.Scan(&seq); err != nil {
			return -1, err
		}
		updateSql := fmt.Sprintf(`update %s set %s = %s + 1 where %s = %s and %s = %d`, s.Tables, s.Sequence, s.Sequence, s.Table, s.BuildParam(1), s.Sequence, seq)
		res, err := s.DB.ExecContext(ctx, updateSql, seqName)
		if err != nil {
			return -1, err
		}
		c, err := res.RowsAffected()
		if err != nil {
			return -1, err
		}
		if c == 0 {
			return -2, nil
		}
		return seq, nil
	} else {
		insertSql := fmt.Sprintf(`insert into %s (%s, %s) values (%s, 2)`, s.Tables, s.Table, s.Sequence, s.BuildParam(1))
		_, err = s.DB.ExecContext(ctx, insertSql, seqName)
		if err != nil {
			x := strings.ToLower(err.Error())
			if strings.Index(x, "unique constraint") >= 0 || strings.Index(x, "Violation of PRIMARY KEY constraint") >= 0 {
				return -2, nil
			}
			return -1, err
		}
		return 1, nil
	}
}
func (s *SequenceAdapter) Reset(ctx context.Context, id string) (int64, error) {
	updateSql := fmt.Sprintf(`update %s set %s = 1 where %s = %s`, s.Tables, s.Sequence, s.Table, s.BuildParam(1))
	res, err := s.DB.ExecContext(ctx, updateSql, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
