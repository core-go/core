package groups

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type GetGroups func(context.Context, string) ([]string, int, error)

type GroupRepository interface {
	GetGroups(ctx context.Context, id string) ([]string, int, error)
}
type GroupPort interface {
	GetGroups(ctx context.Context, id string) ([]string, int, error)
}

func NewGroupAdapter(db *sql.DB, query string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) *GroupAdapter {
	return &GroupAdapter{DB: db, Query: query, Array: toArray}
}
type GroupAdapter struct {
	DB    *sql.DB
	Query string
	Array func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (a *GroupAdapter) GetGroups(ctx context.Context, id string) ([]string, int, error) {
	var ids []string
	var i int
	rows, err := a.DB.QueryContext(ctx, a.Query, id)
	if err != nil {
		return ids, i, err
	}
	defer rows.Close()

	for rows.Next() {
		if er1 := rows.Scan(a.Array(&ids), &i); er1 != nil {
			return nil, i, er1
		}
		return ids, i, rows.Err()
	}
	return ids, i, rows.Err()
}
