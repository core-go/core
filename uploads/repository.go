package uploads

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type StorageRepository interface {
	Load(ctx context.Context, id string) ([]Upload, error)
	Update(ctx context.Context, id string, attachments []Upload) (int64, error)
}

func NewRepository(DB *sql.DB,
	Table string,
	columns FieldColumn, toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}) *SqlRepository {
	return &SqlRepository{DB: DB, Table: Table, Columns: &columns, toArray: toArray}
}

type FieldColumn struct {
	Id   string
	File string
}

type SqlRepository struct {
	DB      *sql.DB
	Table   string
	Columns *FieldColumn
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *SqlRepository) Load(ctx context.Context, id string) ([]Upload, error) {
	var attachments = make([]Upload, 0)
	query := fmt.Sprintf("select %s from %s where %s= $1", s.Columns.File, s.Table, s.Columns.Id)
	rows, err := s.DB.QueryContext(ctx, query, id)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err1 := rows.Scan(s.toArray(&attachments))
		if err1 != nil {
			return nil, err1
		}
		break
	}
	if len(attachments) > 0 {
		return attachments, nil
	}
	return nil, err
}

func (s *SqlRepository) Update(ctx context.Context, id string, attachments []Upload) (int64, error) {
	query := fmt.Sprintf("update %s set %s = $1 where %s =$2", s.Table, s.Columns.File, s.Columns.Id)
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, er0
	}
	res, err := stmt.ExecContext(ctx, s.toArray(attachments), id)
	if err != nil {
		return -1, err
	}
	row, er2 := res.RowsAffected()

	if row < 0 {
		return -1, er2
	}
	return row, er2
}
