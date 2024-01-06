package settings

import "context"

type Settings struct {
	Language   string `yaml:"language" mapstructure:"language" json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	DateFormat string `yaml:"date_format" mapstructure:"date_format" json:"dateFormat,omitempty" gorm:"column:dateformat" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
}

type Service interface {
	Save(ctx context.Context, id string, settings Settings) (int64, error)
}
