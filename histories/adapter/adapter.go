package adapter

import (
	"context"
	"database/sql"
	"fmt"

	h "github.com/core-go/core/histories"
	u "github.com/core-go/core/user"
)

type HistoryAdapter struct {
	DB         *sql.DB
	BuildParam func(int) string
	Table      string
	HistoryId  string
	Resource   string
	Id         string
	User       string
	Time       string
	Data       string
	GetUsers   func(ctx context.Context, ids []string) ([]u.User, error)
}

func UseHistories(db *sql.DB, buildParam func(int) string, getUsers func(ctx context.Context, ids []string) ([]u.User, error), table string, resource string, user string, time string, opts ...string) func(ctx context.Context, resource string, id string, limit int64, nextPageToken string) ([]h.History, string, error) {
	adapter := NewHistoryAdapter(db, buildParam, getUsers, table, resource, user, time, opts...)
	return adapter.GetHistories
}
func NewHistoryAdapter(db *sql.DB, buildParam func(int) string, getUsers func(ctx context.Context, ids []string) ([]u.User, error), table string, resource string, user string, time string, opts ...string) *HistoryAdapter {
	var historyId, id, data string
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
	return &HistoryAdapter{DB: db, BuildParam: buildParam, Table: table, HistoryId: historyId, Resource: resource, Id: id, User: user, Time: time, Data: data, GetUsers: getUsers}
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
	query := fmt.Sprintf("select %s, %s, %s, %s from %s where %s = %s and %s = %s order by %s desc limit %d offset %d",
		a.HistoryId, a.User, a.Time, a.Data, a.Table, a.Id, a.BuildParam(1), a.Resource, a.BuildParam(2), a.Time, limit, offset)
	rows, err := a.DB.QueryContext(ctx, query, id, resource)
	defer rows.Close()
	for rows.Next() {
		var item h.History
		err = rows.Scan(&item.Id, &item.Author, &item.Time, &item.Data)
		if err != nil {
			return histories, "", err
		}
		histories = append(histories, item)
	}
	if len(histories) == 0 {
		return histories, "", nil
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
		usersMap := u.ToMap(users)
		l := len(histories)
		for i := 0; i < l; i++ {
			if u, ok := usersMap[histories[i].Author]; ok {
				ur := h.User{Id: u.Id, Name: u.Name, Email: u.Email, Phone: u.Phone, Url: u.Url}
				histories[i].User = &ur
				histories[i].Author = ""
			}
		}
	}
	return histories, histories[len(histories)-1].Id, nil
}
