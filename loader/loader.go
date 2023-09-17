package loader

import (
	"context"
	"database/sql"
)

type Load func(context.Context, string) ([]string, error)

type Repository interface {
	Load(ctx context.Context, id string) ([]string, error)
}
type Port interface {
	Load(ctx context.Context, id string) ([]string, error)
}
type Query interface {
	Query(ctx context.Context, id string) ([]string, error)
}
type Loader struct {
	db    *sql.DB
	query string
}
func NewQuery(db *sql.DB, query string) Query {
	return &Loader{db: db, query: query}
}
func NewLoader(db *sql.DB, query string) *Loader {
	return &Loader{db: db, query: query}
}
func (a *Loader) Query(ctx context.Context, id string) ([]string, error) {
	return a.Load(ctx, id)
}
func (a *Loader) Load(ctx context.Context, id string) ([]string, error) {
	var ids []string
	rows, err := a.db.QueryContext(ctx, a.query, id)
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
