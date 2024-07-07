package flow

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Group struct {
	Group  string   `yaml:"group" mapstructure:"group" json:"group,omitempty" gorm:"column:user_group" bson:"group,omitempty" dynamodbav:"group,omitempty" firestore:"group,omitempty"`
	Fields []string `yaml:"fields" mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
}

func (g Group) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *Group) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, g)
}

type GetGroups func(context.Context, string) ([]Group, error)

type GroupRepository interface {
	GetGroups(ctx context.Context, id string) ([]Group, error)
}
type GroupPort interface {
	GetGroups(ctx context.Context, id string) ([]Group, error)
}

func NewGroupAdapter(db *sql.DB, buildParam func(int) string, table string, id string, groups string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) *GroupAdapter {
	return &GroupAdapter{DB: db, Table: table, Id: id, Groups: groups, BuildParam: buildParam, Array: toArray}
}

type GroupAdapter struct {
	DB         *sql.DB
	Table      string
	Id         string
	Groups     string
	BuildParam func(int) string
	Array      func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

/*
	func (a *GroupAdapter) CreateGroups(ctx context.Context, id string, groups []Group) (int64, error) {
		query := fmt.Sprintf("insert into %s (%s, %s) values (%s, %s)", a.Table, a.Id, a.Groups, a.BuildParam(1), a.BuildParam(2))
		res, err := a.DB.ExecContext(ctx, query, id, a.Array(groups))
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}

	func (a *GroupAdapter) UpdateGroups(ctx context.Context, id string, groups []Group) (int64, error) {
		query := fmt.Sprintf("update %s set %s = %s where %s = %s", a.Table, a.Groups, a.BuildParam(1), a.Id, a.BuildParam(2))
		res, err := a.DB.ExecContext(ctx, query, a.Array(groups), id)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}
*/
func (a *GroupAdapter) GetGroups(ctx context.Context, id string) ([]Group, error) {
	var groups []Group
	query := fmt.Sprintf("select %s from %s where %s = %s", a.Groups, a.Table, a.Id, a.BuildParam(1))
	rows, err := a.DB.QueryContext(ctx, query, id)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		if er1 := rows.Scan(a.Array(&groups)); er1 != nil {
			return nil, er1
		}
		return groups, rows.Err()
	}
	return groups, rows.Err()
}
