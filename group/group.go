package group

import (
	"context"
	"database/sql"
)

type GroupRepository interface {
	GetUserGroup(ctx context.Context, userId string) (*string, error)
}

type GroupPort interface {
	GetUserGroup(ctx context.Context, userId string) (*string, error)
}

type GroupAdapter struct {
	DB *sql.DB
	Query  string
}

func NewGroupAdapter(db *sql.DB, query string) *GroupAdapter {
	return &GroupAdapter{DB: db, Query: query}
}

func (a *GroupAdapter) GetUserGroup(ctx context.Context, userId string) (*string, error) {
	var groupId string
	rows, err := a.DB.QueryContext(ctx, a.Query, userId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&groupId)
		if err != nil {
			return nil, err
		}
	}
	return &groupId, nil
}
