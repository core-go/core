package passcode

import "time"

type Passcode struct {
	Id        string    `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id;primary_key" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"-"`
	Passcode  string    `yaml:"passcode" mapstructure:"passcode" json:"passcode,omitempty" gorm:"column:passcode" bson:"passcode,omitempty" dynamodbav:"passcode,omitempty" firestore:"passcode,omitempty"`
	ExpiredAt time.Time `yaml:"expired_at" mapstructure:"expired_at" json:"expiredAt,omitempty" gorm:"column:expiredat" bson:"expiredAt,omitempty" dynamodbav:"expiredAt,omitempty" firestore:"expiredAt,omitempty"`
}
