package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	n "github.com/core-go/core/notification"
)

type NotificationAdapter struct {
	DB *sql.DB
	Generate func(ctx context.Context) (string, error)
	BuildParam func(int)string
	Tx string
	Table string
	Id  string
	Sender string
	Receiver string
	Url string
	Message string
	Time  string
	Status  string
	SuffixColumn string
	SuffixValue string
}
func NewNotificationAdapter(db *sql.DB, generate func(ctx context.Context) (string, error), buildParam func(int)string, table string, tx string, time string, opts...string) *NotificationAdapter {
	var id, sender, receiver, message, url, suffixColumn, suffixValue string
	if len(opts) > 0 {
		id = opts[0]
	} else {
		id = "id"
	}
	if len(opts) > 1 {
		sender = opts[1]
	} else {
		sender = "sender"
	}
	if len(opts) > 2 {
		receiver = opts[2]
	} else {
		receiver = "receiver"
	}
	if len(opts) > 3 {
		message = opts[3]
	} else {
		message = "message"
	}
	if len(opts) > 4 {
		url = opts[4]
	} else {
		url = "url"
	}
	if len(opts) > 5 {
		suffixColumn = opts[5]
	} else {
		suffixColumn = ""
	}
	if len(opts) > 6 {
		suffixValue = opts[6]
	} else {
		suffixValue = ""
	}
	return &NotificationAdapter{DB: db, Generate: generate, BuildParam: buildParam, Table: table, Tx: tx, SuffixColumn: suffixColumn, SuffixValue: suffixValue, Time: time, Id: id, Sender: sender, Receiver: receiver, Message: message, Url: url}
}
func (a *NotificationAdapter) Push(ctx context.Context, noti *n.Notification) (int64, error) {
	if noti == nil {
		return 0, nil
	}
	id, err := a.Generate(ctx)
	if err != nil {
		return -1, nil
	}
	now := time.Now()
	query := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s %s) values (%s,%s,%s,%s,%s,%s %s)", a.Table,
		a.Id, a.Sender, a.Receiver, a.Message, a.Url, a.Time, a.SuffixColumn,
		a.BuildParam(1), a.BuildParam(2), a.BuildParam(3), a.BuildParam(4), a.BuildParam(5),a.BuildParam(6),a.SuffixValue)
	tx := GetExec(ctx, a.DB, a.Tx)
	res, err := tx.ExecContext(ctx, query, id, noti.Sender, noti.Receiver, noti.Message, noti.Url, now)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
func (a *NotificationAdapter) PushNotifications(ctx context.Context, ns []*n.Notification) (int64, error) {
	l := len(ns)
	for i :=0; i < l; i++ {
		id, err := a.Generate(ctx)
		if err != nil {
			return -1, err
		}
		ns[i].Id = id
	}
	now := time.Now()
	query := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s %s) values (%s,%s,%s,%s,%s,%s %s)", a.Table,
		a.Id, a.Sender, a.Receiver, a.Message, a.Url, a.Time, a.SuffixColumn,
		a.BuildParam(1), a.BuildParam(2), a.BuildParam(3), a.BuildParam(4), a.BuildParam(5),a.BuildParam(6),a.SuffixValue)
	tx := GetExec(ctx, a.DB, a.Tx)
	var sum int64
	sum = 0
	for i :=0; i < l; i++ {
		res, err := tx.ExecContext(ctx, query, ns[i].Id, ns[i].Sender, ns[i].Receiver, ns[i].Message, ns[i].Url, now)
		if err != nil {
			return -1, err
		}
		count, er2 := res.RowsAffected()
		if er2 != nil {
			return sum, er2
		}
		sum = sum + count
	}
	return sum, nil
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