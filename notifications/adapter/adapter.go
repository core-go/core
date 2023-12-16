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
func UseNotification(db *sql.DB, buildParam func(int)string, getUsers func(ctx context.Context, ids []string) ([]u.User, error), readValue interface{}, table string, opts...string) func(ctx context.Context, receiver string, read *bool, limit int64, nextPageToken string) ([]n.Notification, string, error) {
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
func (a *NotificationAdapter) GetNotifications(ctx context.Context, receiver string, read *bool, limit int64, nextPageToken string) ([]n.Notification, string, error) {
	if limit <= 0 {
		limit = 20
	}
	var items []n.Notification
	var whereRead = ""
	if read != nil && *read == true {
		whereRead = fmt.Sprintf(" and %s = %s", a.Read, a.BuildParam(2))
	}
	var offset int64
	if len(nextPageToken) > 0 {
		k := 2
		if read != nil && *read == true {
			k = 3
		}
		positionQuery := fmt.Sprintf("select position from (select %s, row_number() over(order by %s desc) as position from %s where %s = %s %s) result where %s = %s",
			a.Id, a.Time, a.Table, a.Receiver, a.BuildParam(1), whereRead, a.Id, a.BuildParam(k))
		var row *sql.Row
		if read != nil && *read == true {
			row = a.DB.QueryRowContext(ctx, positionQuery, receiver, a.ReadValue, nextPageToken)
		} else {
			row = a.DB.QueryRowContext(ctx, positionQuery, receiver, nextPageToken)
		}
		if row.Err() != nil {
			return items, "", row.Err()
		}
		err := row.Scan(&offset)
		if offset < 0 {
			offset = 0
		}
		if err != nil {
			return items, "", err
		}
	} else {
		offset = 0
	}
	query := fmt.Sprintf("select %s, %s, %s, %s, %s, %s from %s where %s = %s %s order by %s desc limit %d offset %d",
		a.Id, a.Time, a.Read, a.Sender, a.Message, a.Url, a.Table, a.Receiver, a.BuildParam(1), whereRead, a.Time, limit, offset)
	var rows *sql.Rows
	var err error
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
			return items, "", err
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return items, "", nil
	}
	if a.GetUsers != nil {
		var userIds []string
		for _, hi := range items {
			userIds = append(userIds, hi.Sender)
		}
		ids := u.Unique(userIds)
		users, err := a.GetUsers(ctx, ids)
		if err != nil {
			return items, "", err
		}
		usersMap := u.ToMap(users)
		l := len(items)
		for i := 0; i < l; i++ {
			if u, ok := usersMap[items[i].Sender]; ok {
				ur := n.Notifier{Id: u.Id, Name: u.Name, Email: u.Email, Phone: u.Phone, Url: u.Url}
				items[i].Notifier = &ur
				items[i].Sender = ""
			}
		}
	}
	return items, items[len(items) - 1].Id, nil
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
