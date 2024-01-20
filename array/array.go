package approver

import (
	"context"
	"database/sql"
)

type GetArray func(context.Context, string) ([]string, error)

type ArrayRepository interface {
	GetArray(ctx context.Context, id string) ([]string, error)
}
type ArrayPort interface {
	GetArray(ctx context.Context, id string) ([]string, error)
}

func NewArrayAdapter(db *sql.DB, query string) *ArrayAdapter {
	return &ArrayAdapter{DB: db, Query: query}
}
type ArrayAdapter struct {
	DB    *sql.DB
	Query string
}

func (a *ArrayAdapter) GetArray(ctx context.Context, id string) ([]string, error) {
	var ids []string
	rows, err := a.DB.QueryContext(ctx, a.Query, id)
	if err != nil {
		return ids, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ids); err != nil {
			return nil, err
		}
	}
	return ids, rows.Err()
}
