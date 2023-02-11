package user

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type User struct {
	Id    string  `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name  *string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Email *string `yaml:"email" mapstructure:"email" json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone *string `yaml:"phone" mapstructure:"phone" json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	Url   *string `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
}

type UserPort interface {
	Load(ctx context.Context, id string) (*User, error)
	Query(ctx context.Context, ids []string) ([]User, error)
}

type UserAdapter struct {
	db         *sql.DB
	Select     string
	BuildParam func(int) string
}

func NewUserAdapter(db *sql.DB, query string, opts ...func(i int) string) *UserAdapter {
	var buildParam func(i int) string
	if len(opts) > 0 && opts[0] != nil {
		buildParam = opts[0]
	} else {
		buildParam = getBuild(db)
	}
	return &UserAdapter{
		db:         db,
		Select:     query,
		BuildParam: buildParam,
	}
}

func (r *UserAdapter) Load(ctx context.Context, id string) (*User, error) {
	p := make([]string, 0)
	p = append(p, id)
	users, err := r.Query(ctx, p)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil
}

func (r *UserAdapter) Query(ctx context.Context, ids []string) ([]User, error) {
	var users []User
	if len(ids) == 0 {
		return users, nil
	}
	le := len(ids)
	p := make([]interface{}, 0)
	for _, str := range ids {
		p = append(p, str)
	}
	var arrValue []string
	for i := 1; i <= le; i++ {
		param := r.BuildParam(i)
		arrValue = append(arrValue, param)
	}
	query := r.Select + fmt.Sprintf(" in (%s)", strings.Join(arrValue, ","))
	rows, err := r.db.QueryContext(ctx, query, p...)
	defer rows.Close()

	for rows.Next() {
		var row User
		if err := rows.Scan(&row.Id, &row.Name, &row.Email, &row.Phone, &row.Url); err != nil {
			return users, err
		}
		users = append(users, row)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	SortById(users)
	return users, nil
}

func ToMap(rows []User) map[string]*User {
	rs := make(map[string]*User, 0)
	for _, row := range rows {
		rs[row.Id] = &row
	}
	return rs
}

func Unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
func SortById(users []User) {
	sort.Slice(users, func(i, j int) bool { return users[i].Id < users[j].Id })
}
func BinarySearch(id string, a []User) (result int, searchCount int) {
	mid := len(a) / 2
	x := strings.Compare(a[mid].Id, id)
	switch {
	case len(a) == 0:
		result = -1 // not found
	case x > 0:
		result, searchCount = BinarySearch(id, a[:mid])
	case x < 0:
		result, searchCount = BinarySearch(id, a[mid+1:])
		if result >= 0 { // if anything but the -1 "not found" result
			result += mid + 1
		}
	default: // a[mid] == id
		result = mid // found
	}
	searchCount++
	return
}
func getBuild(db *sql.DB) func(i int) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return buildDollarParam
	case "*godror.drv":
		return buildOracleParam
	case "*mssql.Driver":
		return buildMsSqlParam
	default:
		return buildParam
	}
}
func buildParam(i int) string {
	return "?"
}
func buildOracleParam(i int) string {
	return ":" + strconv.Itoa(i)
}
func buildMsSqlParam(i int) string {
	return "@p" + strconv.Itoa(i)
}
func buildDollarParam(i int) string {
	return "$" + strconv.Itoa(i)
}
