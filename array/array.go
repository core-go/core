package array

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type GetArray func(context.Context, string) ([]string, error)

type ArrayRepository interface {
	GetArray(ctx context.Context, id string) ([]string, error)
}
type ArrayPort interface {
	GetArray(ctx context.Context, id string) ([]string, error)
}

func NewArrayAdapter(db *sql.DB, query string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) *ArrayAdapter {
	return &ArrayAdapter{DB: db, Query: query, Array: toArray}
}
type ArrayAdapter struct {
	DB    *sql.DB
	Query string
	Array     func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (a *ArrayAdapter) GetArray(ctx context.Context, id string) ([]string, error) {
	var ids []string
	rows, err := a.DB.QueryContext(ctx, a.Query, id)
	if err != nil {
		return ids, err
	}
	defer rows.Close()

	for rows.Next() {
		if er1 := rows.Scan(a.Array(&ids)); er1 != nil {
			return nil, er1
		}
		return ids, rows.Err()
	}
	return ids, rows.Err()
}
