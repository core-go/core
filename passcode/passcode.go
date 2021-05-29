package passcode

import "time"

type Passcode struct {
	Id        string    `json:"id,omitempty" gorm:"column:id;primary_key" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"-"`
	Passcode  string    `json:"passcode,omitempty" gorm:"column:passcode" bson:"passcode,omitempty" dynamodbav:"passcode,omitempty" firestore:"passcode,omitempty"`
	ExpiredAt time.Time `json:"expiredAt,omitempty" gorm:"column:expiredat" bson:"expiredAt,omitempty" dynamodbav:"expiredAt,omitempty" firestore:"expiredAt,omitempty"`
}
