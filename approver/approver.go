package approver

import (
	"context"
	"database/sql"
)

type GetApprovers func(context.Context, string) ([]string, error)

type ApproversRepository interface {
	GetArray(ctx context.Context, id string) ([]string, error)
}
type ApproversPort interface {
	GetArray(ctx context.Context, id string) ([]string, error)
}

func NewApproversAdapter(db *sql.DB, query string) *ApproversAdapter{
	return &ApproversAdapter{DB: db, Query: query}
}
type ApproversAdapter struct {
	DB    *sql.DB
	Query string
}

func (a *ApproversAdapter) GetArray(ctx context.Context, id string) ([]string, error) {
	var ids []string
	rows, err := a.DB.QueryContext(ctx, a.Query, id)
	defer rows.Close()
	if err != nil {
		return ids, err
	}

	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		ids = append(ids, s)
	}
	return ids, rows.Err()
}
