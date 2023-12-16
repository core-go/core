package notifications

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Data map[string]interface{}

func (h Data) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *Data) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &h)
}

type Notification struct {
	Id       string     `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Sender   string     `yaml:"sender" mapstructure:"sender" json:"sender,omitempty" gorm:"column:sender" bson:"sender,omitempty" dynamodbav:"sender,omitempty" firestore:"sender,omitempty"`
	Url      string     `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Message  string     `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
	Time     *time.Time `yaml:"time" mapstructure:"time" json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	Read     *bool      `yaml:"read" mapstructure:"read" json:"read,omitempty" gorm:"column:read" bson:"read,omitempty" dynamodbav:"read,omitempty" firestore:"read,omitempty"`
	Notifier *Notifier  `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
}
type Notifier struct {
	Id    string  `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name  *string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Email *string `yaml:"email" mapstructure:"email" json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone *string `yaml:"phone" mapstructure:"phone" json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	Url   *string `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
}

type NotificationsPort interface {
	GetNotifications(ctx context.Context, receiver string, read *bool, limit int64, nextPageToken string) ([]Notification, string, error)
	SetRead(ctx context.Context, id string, v bool) (int64, error)
}
