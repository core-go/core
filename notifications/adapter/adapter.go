package adapter

import (
	"context"
	"database/sql"
	"fmt"

	n "github.com/core-go/core/notifications"
	u "github.com/core-go/core/user"
)

type NotificationAdapter struct {
	DB *sql.DB
	BuildParam func(int)string
	Table string
	Id string
	Time string
	Receiver string
	Sender string
	Message string
	Url string
	Read string
	ReadValue interface{}
	GetUsers func(ctx context.Context, ids []string) ([]u.User, error)
}
func UseNotification(db *sql.DB, buildParam func(int)string, getUsers func(ctx context.Context, ids []string) ([]u.User, error), readValue interface{}, table string, opts...string) func(ctx context.Context, receiver string, read *bool, limit int64, offset int64) ([]n.Notification, int64, error) {
	adapter := NewNotificationAdapter(db, buildParam, getUsers, readValue, table, opts...)
	return adapter.GetNotifications
}
func NewNotificationAdapter(db *sql.DB, buildParam func(int)string, getUsers func(ctx context.Context, ids []string) ([]u.User, error), readValue interface{}, table string, opts...string) *NotificationAdapter {
	var receiver, sender, time, message, url, id string
	if len(opts) > 0 {
		receiver = opts[0]
	} else {
		receiver = "receiver"
	}
	if len(opts) > 1 {
		sender = opts[1]
	} else {
		sender = "sender"
	}
	if len(opts) > 2 {
		time = opts[2]
	} else {
		time = "time"
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
		id = opts[5]
	} else {
		id = "id"
	}
	return &NotificationAdapter{DB: db, BuildParam: buildParam, Table: table, ReadValue: readValue, Receiver: receiver, Sender: sender, Time: time, Message: message, Url: url, Id: id, GetUsers: getUsers}
}
func (a *NotificationAdapter) GetNotifications(ctx context.Context, receiver string, read *bool, limit int64, offset int64) ([]n.Notification, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	var items []n.Notification
	countQuery := fmt.Sprintf("select count(*) from %s where %s = %s", a.Table, a.Receiver, a.BuildParam(1))
	var row *sql.Row
	var whereRead = ""
	if read != nil && *read == true {
		whereRead = fmt.Sprintf(" and %s = %s", a.Read, a.BuildParam(2))
		countQuery = countQuery + whereRead
		row = a.DB.QueryRowContext(ctx, countQuery, receiver, a.ReadValue)
	} else {
		row = a.DB.QueryRowContext(ctx, countQuery, receiver)
	}
	if row.Err() != nil {
		return items, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return items, total, err
	}
	query := fmt.Sprintf("select %s, %s, %s, %s, %s, %s from %s where %s = %s %s order by %s desc limit %d offset %d",
		a.Id, a.Time, a.Read, a.Sender, a.Message, a.Url, a.Table, a.Receiver, a.BuildParam(1), whereRead, a.Time, limit, offset)
	var rows *sql.Rows
	if read != nil && *read == true {
		rows, err = a.DB.QueryContext(ctx, query, receiver, a.ReadValue)
	} else {
		rows, err = a.DB.QueryContext(ctx, query, receiver)
	}
	defer rows.Close()
	for rows.Next() {
		var item n.Notification
		err = rows.Scan(&item.Id, &item.Time, &item.Read, &item.Sender, &item.Message, &item.Url)
		if err != nil {
			return items, total, err
		}
		items = append(items, item)
	}
	if a.GetUsers != nil {
		var userIds []string
		for _, hi := range items {
			userIds = append(userIds, hi.Sender)
		}
		ids := u.Unique(userIds)
		users, err := a.GetUsers(ctx, ids)
		if err != nil {
			return items, total, err
		}
		l := len(items)
		for i := 0; i < l; i++ {
			p, _ := u.BinarySearch(items[i].Sender, users)
			if p >= 0 {
				us := users[p]
				ur := n.User{Id: us.Id, Name: us.Name, Email: us.Email, Phone: us.Phone, Url: us.Url}
				items[i].User = &ur
				items[i].Sender = ""
			}
		}
	}
	return items, total, nil
}
func (a *NotificationAdapter) SetRead(ctx context.Context, id string, v bool) (int64, error) {
	p := "null"
	k := 1
	if v {
		p = a.BuildParam(1)
		k = 2
	}
	query := fmt.Sprintf("update %s set %s = %s where %s = %s", a.Table, a.Read, p, a.Id, a.BuildParam(k))
	var res sql.Result
	var err error
	if v {
		res, err = a.DB.ExecContext(ctx, query, a.ReadValue, id)
	} else {
		res, err = a.DB.ExecContext(ctx, query, id)
	}
	if err != nil {
		return -1, nil
	}
	return res.RowsAffected()
}
