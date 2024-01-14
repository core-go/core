package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	h "github.com/core-go/core/history"
)

type HistoryAdapter struct {
	DB         *sql.DB
	Generate   func(ctx context.Context) (string, error)
	BuildParam func(int) string
	Type       string
	Tx         string
	Table      string
	HistoryId  string
	Resource   string
	Id         string
	User       string
	Time       string
	Data       string
}

func NewHistoryAdapter(db *sql.DB, generate func(ctx context.Context) (string, error), buildParam func(int) string, stype string, resource string, table string, user string, time string, opts ...string) *HistoryAdapter {
	var historyId, id, data, tx string
	if len(opts) > 0 {
		historyId = opts[0]
	} else {
		historyId = "history_id"
	}
	if len(opts) > 1 {
		id = opts[1]
	} else {
		id = "id"
	}
	if len(opts) > 2 {
		data = opts[2]
	} else {
		data = "data"
	}
	if len(opts) > 3 {
		tx = opts[3]
	} else {
		tx = "tx"
	}
	return &HistoryAdapter{DB: db, Generate: generate, BuildParam: buildParam, Type: stype, Tx: tx, Table: table, HistoryId: historyId, Resource: resource, Id: id, User: user, Time: time, Data: data}
}
func (a *HistoryAdapter) Create(ctx context.Context, id string, userId string, data map[string]interface{}) (int64, error) {
	hid, err := a.Generate(ctx)
	if err != nil {
		return -1, nil
	}
	history := &h.History{Data: data}
	now := time.Now()
	query := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s) values (%s,%s,%s,%s,%s,%s)", a.Table,
		a.HistoryId, a.Type, a.Id, a.User, a.Time, a.Data,
		a.BuildParam(1), a.BuildParam(2), a.BuildParam(3), a.BuildParam(4), a.BuildParam(5), a.BuildParam(6))
	tx := GetExec(ctx, a.DB, a.Tx)
	res, err := tx.ExecContext(ctx, query, hid, a.Resource, id, userId, now, history.Data)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

type Executor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func GetExec(ctx context.Context, db *sql.DB, name string) Executor {
	txi := ctx.Value(name)
	if txi != nil {
		txx, ok := txi.(*sql.Tx)
		if ok {
			return txx
		}
	}
	return db
}
