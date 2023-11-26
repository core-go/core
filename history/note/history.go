package history

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
	HistoryId string     `yaml:"historyId" mapstructure:"historyId" json:"historyId,omitempty" gorm:"column:history_id" bson:"historyId,omitempty" dynamodbav:"historyId,omitempty" firestore:"historyId,omitempty"`
	Resource  string     `yaml:"resource" mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Id        string     `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Author    string     `yaml:"author" mapstructure:"author" json:"author,omitempty" gorm:"column:author" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty"`
	Time      *time.Time `yaml:"time" mapstructure:"time" json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	Data      Data       `yaml:"data" mapstructure:"data" json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
	Note      string     `yaml:"note" mapstructure:"note" json:"note,omitempty" gorm:"column:note" bson:"note,omitempty" dynamodbav:"note,omitempty" firestore:"note,omitempty"`
}

type HistoryPort interface {
	Create(ctx context.Context, id string, userId string, data map[string]interface{}, note string) (int64, error)
}
