package histories

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

type History struct {
	Id     string     `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Author string     `yaml:"author" mapstructure:"author" json:"author,omitempty" gorm:"column:author" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty"`
	Time   *time.Time `yaml:"time" mapstructure:"time" json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	Data   Data       `yaml:"data" mapstructure:"data" json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
	Note   *string    `yaml:"note" mapstructure:"note" json:"note,omitempty" gorm:"column:note" bson:"note,omitempty" dynamodbav:"note,omitempty" firestore:"note,omitempty"`
	User   *User      `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
}
type User struct {
	Id    string  `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name  *string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Email *string `yaml:"email" mapstructure:"email" json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone *string `yaml:"phone" mapstructure:"phone" json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	Url   *string `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
}

type HistoriesPort interface {
	GetHistories(ctx context.Context, resource string, id string, limit int64, nextPageToken int64) ([]History, string, error)
}
