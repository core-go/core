package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	h "github.com/core-go/core/history/note"
)

type HistoryAdapter struct {
	DB *sql.DB
	Generate func(ctx context.Context) (string, error)
	BuildParam func(int)string
	Type string
	Tx string
	Table string
	HistoryId string
	Resource string
	Id  string
	User string
	Time string
	Data string
	Note string
}
func NewHistoryAdapter(db *sql.DB, generate func(ctx context.Context) (string, error), buildParam func(int)string, stype string, resource string, table string, user string, time string, opts...string) *HistoryAdapter {
	var historyId, id, note, data, tx string
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
		note = opts[2]
	} else {
		note = "note"
	}
	if len(opts) > 3 {
		data = opts[3]
	} else {
		data = "data"
	}
	if len(opts) > 4 {
		tx = opts[4]
	} else {
		tx = "tx"
	}
	return &HistoryAdapter{DB: db, Generate: generate, BuildParam: buildParam, Type: stype, Tx: tx, Table: table, HistoryId: historyId, Resource: resource, Id: id, User: user, Time: time, Data: data, Note: note}
}
func (a *HistoryAdapter) Create(ctx context.Context, id string, userId string, data map[string]interface{}, note string) (int64, error) {
	hid, err := a.Generate(ctx)
	if err != nil {
		return -1, nil
	}
	history := &h.History{
		Data: data,
		Note: note,
	}
	now := time.Now()
	query := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s,%s) values (%s,%s,%s,%s,%s,%s,%s)", a.Table,
		hid, a.Type, a.Id, a.User, a.Time, a.Data, a.Note,
		a.BuildParam(1), a.BuildParam(2), a.BuildParam(3), a.BuildParam(4), a.BuildParam(5), a.BuildParam(6), a.BuildParam(7))
	tx := GetExec(ctx, a.DB, a.Tx)
	res, err := tx.ExecContext(ctx, query, hid, a.Resource, id, userId, now, history.Data, note)
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
