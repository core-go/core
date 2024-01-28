package groups

import (
	"context"
	"database/sql"
)

type GetGroups func(context.Context, string) ([]string, int32, error)

type GroupRepository interface {
	GetGroups(ctx context.Context, id string) ([]string, int32, error)
}
type GroupPort interface {
	GetGroups(ctx context.Context, id string) ([]string, int32, error)
}

func NewGroupAdapter(db *sql.DB, query string) *GroupAdapter {
	return &GroupAdapter{DB: db, Query: query}
}
type GroupAdapter struct {
	DB    *sql.DB
	Query string
}

func (a *GroupAdapter) GetGroups(ctx context.Context, id string) ([]string, int32, error) {
	var ids []string
	var i int32
	rows, err := a.DB.QueryContext(ctx, a.Query, id)
	if err != nil {
		return ids, i, err
	}
	defer rows.Close()

	for rows.Next() {
		if er1 := rows.Scan(&ids, &i); er1 != nil {
			return nil, i, er1
		}
		return ids, i, rows.Err()
	}
	return ids, i, rows.Err()
}
