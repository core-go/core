package approval

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	Group string     `yaml:"group" mapstructure:"group" json:"group,omitempty" gorm:"column:group" bson:"group,omitempty" dynamodbav:"group,omitempty" firestore:"group,omitempty"`
	User  *string    `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Time  *time.Time `yaml:"time" mapstructure:"time" json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
}

func (u User) Value() (driver.Value, error) {
	return json.Marshal(u)
}
func (u *User) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &u)
}

func BuildApprovers(ctx context.Context, id string, getArray func(context.Context, string) ([]string, error)) (string, []User, error) {
	groups, err := getArray(ctx, id)
	if err != nil {
		return "", nil, err
	}
	var approvers []User
	l := len(groups)
	if l > 0 {
		for i := 0; i < l; i++ {
			approvers = append(approvers, User{Group: groups[i]})
		}
	}
	return "", approvers, nil
}

func Approve(ctx context.Context, userId string, approvers []User, getUserGroup func(context.Context, string) (*string, error)) (*string, int, error) {
	userGroup, err := getUserGroup(ctx, userId)
	if err != nil {
		return nil, 0, err
	}
	if userGroup == nil {
		return nil, -1, err
	}
	l := len(approvers)
	if l == 0 {
		return nil, 0, nil
	}
	for i := 0; i < l; i++ {
		if approvers[i].User != nil {
			continue
		}
		if approvers[i].Group != *userGroup {
			return nil, -2, nil
		}
		now := time.Now()
		approvers[i].User = &userId
		approvers[i].Time = &now
		if i < l - 1 {
			return &approvers[i+1].Group, i + 1, nil
		}
		return nil, i+1, nil
	}
	return nil, 0, nil
}

func Reject(ctx context.Context, userId string, approvers []User, getUserGroup func(context.Context, string) (*string, error)) (int, error) {
	userGroup, err := getUserGroup(ctx, userId)
	if err != nil {
		return 0, err
	}
	if userGroup == nil {
		return -1, err
	}
	l := len(approvers)
	if l == 0 {
		return 0, nil
	}
	for i := 0; i < l; i++ {
		if approvers[i].User != nil {
			continue
		}
		if approvers[i].Group != *userGroup {
			return -2, nil
		}
		now := time.Now()
		approvers[i].User = &userId
		approvers[i].Time = &now
		return i, nil
	}
	return 0, nil
}

func Handle(ctx context.Context, userId string, groups []string, approvers *[]User, getGroup func(context.Context, string) (*string, error)) (int64, error) {
	userGroup, err := getGroup(ctx, userId)
	if err != nil || userGroup == nil {
		return -3, err
	}
	for _, approver := range *approvers {
		if approver.Group == *userGroup {
			return -2, nil
		}
	}
	if len(groups) > 0 {
		found := false
		for _, group := range groups {
			if *userGroup == group {
				found = true
				break
			}
		}
		if !found {
			return -3, nil
		}
	}
	now := time.Now()
	*approvers = append(*approvers, User{User: &userId, Group: *userGroup, Time: &now})
	return 0, nil
}
