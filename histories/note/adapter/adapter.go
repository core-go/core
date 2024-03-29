package adapter

import (
	"context"
	"database/sql"
	"fmt"

	h "github.com/core-go/core/histories/note"
	u "github.com/core-go/core/user"
)

type HistoryAdapter struct {
	DB *sql.DB
	BuildParam func(int)string
	Table string
	HistoryId string
	Resource string
	Id  string
	User string
	Time string
	Data string
	Note string
	GetUsers func(ctx context.Context, ids []string) ([]u.User, error)
}
func UseHistories(db *sql.DB, buildParam func(int)string, getUsers func(ctx context.Context, ids []string) ([]u.User, error), table string, resource string, user string, time string, opts...string) func(ctx context.Context, resource string, id string, limit int64, nextPageToken string) ([]h.History, string, error) {
	adapter := NewHistoryAdapter(db, buildParam, getUsers, table, resource, user, time, opts...)
	return adapter.GetHistories
}
func NewHistoryAdapter(db *sql.DB, buildParam func(int)string, getUsers func(ctx context.Context, ids []string) ([]u.User, error), table string, resource string, user string, time string, opts...string) *HistoryAdapter {
	var historyId, id, note, data string
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
	return &HistoryAdapter{DB: db, BuildParam: buildParam, Table: table, HistoryId: historyId, Resource: resource, Id: id, User: user, Time: time, Data: data, Note: note, GetUsers: getUsers}
}
func (a *HistoryAdapter) GetHistories(ctx context.Context, resource string, id string, limit int64, nextPageToken string) ([]h.History, string, error) {
	if limit <= 0 {
		limit = 20
	}
	var histories []h.History
	var offset int64
	if len(nextPageToken) > 0 {
		positionQuery := fmt.Sprintf("select position from (select %s, row_number() over(order by %s desc) as position from %s where %s = %s and %s = %s) result where %s = %s",
			a.HistoryId, a.Time, a.Table, a.Id, a.BuildParam(1), a.Resource, a.BuildParam(2), a.HistoryId, a.BuildParam(3))
		row := a.DB.QueryRowContext(ctx, positionQuery, id, resource, nextPageToken)
		if row.Err() != nil {
			return histories, "", row.Err()
		}
		err := row.Scan(&offset)
		if offset < 0 {
			offset = 0
		}
		if err != nil {
			return histories, "", err
		}
	} else {
		offset = 0
	}
	query := fmt.Sprintf("select %s, %s, %s, %s, %s from %s where %s = %s and %s = %s order by %s desc limit %d offset %d",
		a.HistoryId, a.User, a.Time, a.Note, a.Data, a.Table, a.Id, a.BuildParam(1), a.Resource, a.BuildParam(2), a.Time, limit, offset)
	rows, err := a.DB.QueryContext(ctx, query, id, resource)
	if err != nil {
		return histories, "", err
	}
	defer rows.Close()
	for rows.Next() {
		var item h.History
		err = rows.Scan(&item.Id, &item.Author, &item.Time, &item.Note, &item.Data)
		if err != nil {
			return histories, "", err
		}
		histories = append(histories, item)
	}
	if a.GetUsers != nil {
		var userIds []string
		for _, hi := range histories {
			userIds = append(userIds, hi.Author)
		}
		ids := u.Unique(userIds)
		users, err := a.GetUsers(ctx, ids)
		if err != nil {
			return histories, "", err
		}
		l := len(histories)
		for i := 0; i < l; i++ {
			p, _ := u.BinarySearch(histories[i].Author, users)
			if p >= 0 {
				us := users[p]
				ur := h.User{Id: us.Id, Name: us.Name, Email: us.Email, Phone: us.Phone, Url: us.Url}
				histories[i].User = &ur
				histories[i].Author = ""
			}
		}

	}
	return histories, histories[len(histories) - 1].Id, nil
}
