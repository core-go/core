package approver

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
	if err != nil {
		return ids, err
	}

	for rows.Next() {
		var s string
		if er1 := rows.Scan(&s); er1 != nil {
			return nil, er1
		}
		ids = append(ids, s)
	}
	return ids, rows.Err()
}
