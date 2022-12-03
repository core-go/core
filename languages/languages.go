package languages

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Languages struct {
	En string `yaml:"en" mapstructure:"en" json:"en,omitempty" gorm:"column:en" bson:"en,omitempty" dynamodbav:"en,omitempty" firestore:"en,omitempty"`
	Vi string `yaml:"vi" mapstructure:"vi" json:"vi,omitempty" gorm:"column:vi" bson:"vi,omitempty" dynamodbav:"vi,omitempty" firestore:"vi,omitempty"`
	Zh string `yaml:"zh" mapstructure:"zh" json:"zh,omitempty" gorm:"column:zh" bson:"zh,omitempty" dynamodbav:"zh,omitempty" firestore:"zh,omitempty"`
	Th string `yaml:"th" mapstructure:"th" json:"th,omitempty" gorm:"column:th" bson:"th,omitempty" dynamodbav:"th,omitempty" firestore:"th,omitempty"`
	Lo string `yaml:"lo" mapstructure:"lo" json:"lo,omitempty" gorm:"column:lo" bson:"lo,omitempty" dynamodbav:"lo,omitempty" firestore:"lo,omitempty"`
	Km string `yaml:"km" mapstructure:"km" json:"km,omitempty" gorm:"column:km" bson:"km,omitempty" dynamodbav:"km,omitempty" firestore:"km,omitempty"`
	Id string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Ms string `yaml:"ms" mapstructure:"ms" json:"ms,omitempty" gorm:"column:ms" bson:"ms,omitempty" dynamodbav:"ms,omitempty" firestore:"ms,omitempty"`
	Ja string `yaml:"ja" mapstructure:"ja" json:"ja,omitempty" gorm:"column:ja" bson:"ja,omitempty" dynamodbav:"ja,omitempty" firestore:"ja,omitempty"`
	Ko string `yaml:"ko" mapstructure:"ko" json:"ko,omitempty" gorm:"column:ko" bson:"ko,omitempty" dynamodbav:"ko,omitempty" firestore:"ko,omitempty"`
	De string `yaml:"de" mapstructure:"de" json:"de,omitempty" gorm:"column:de" bson:"de,omitempty" dynamodbav:"de,omitempty" firestore:"de,omitempty"`
	Fr string `yaml:"fr" mapstructure:"fr" json:"fr,omitempty" gorm:"column:fr" bson:"fr,omitempty" dynamodbav:"fr,omitempty" firestore:"fr,omitempty"`
	It string `yaml:"it" mapstructure:"it" json:"it,omitempty" gorm:"column:it" bson:"it,omitempty" dynamodbav:"it,omitempty" firestore:"it,omitempty"`
	Es string `yaml:"es" mapstructure:"es" json:"es,omitempty" gorm:"column:es" bson:"es,omitempty" dynamodbav:"es,omitempty" firestore:"es,omitempty"`
	Pt string `yaml:"pt" mapstructure:"pt" json:"pt,omitempty" gorm:"column:pt" bson:"pt,omitempty" dynamodbav:"pt,omitempty" firestore:"pt,omitempty"`
}

func (l Languages) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *Languages) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &l)
}
