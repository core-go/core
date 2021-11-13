package search

import "time"

type DateRange struct {
	Min *time.Time `mapstructure:"min" json:"min,omitempty" gorm:"column:startdate" bson:"min,omitempty" dynamodbav:"min,omitempty" firestore:"min,omitempty"`
	Max *time.Time `mapstructure:"max" json:"max,omitempty" gorm:"column:max" bson:"max,omitempty" dynamodbav:"max,omitempty" firestore:"max,omitempty"`
}
