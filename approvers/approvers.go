package approvers

import (
	"context"
	"database/sql"
)

type GetApprovers func(context.Context, string) ([]string, error)

type ApproversRepository interface {
	GetApprovers(ctx context.Context, id string) ([]string, error)
}
type ApproversPort interface {
	GetApprovers(ctx context.Context, id string) ([]string, error)
}

func NewApproversAdapter(db *sql.DB, query string) *ApproversAdapter{
	return &ApproversAdapter{DB: db, Query: query}
}
type ApproversAdapter struct {
	DB    *sql.DB
	Query string
}

func (a *ApproversAdapter) GetApprovers(ctx context.Context, id string) ([]string, error) {
	var ids []string
	rows, err := a.DB.QueryContext(ctx, a.Query, id)
	defer rows.Close()

	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		ids = append(ids, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}
