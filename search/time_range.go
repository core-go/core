package search

import "time"

type TimeRange struct {
	StartTime *time.Time `mapstructure:"startTime" json:"startTime,omitempty" gorm:"column:starttime" bson:"startTime,omitempty" dynamodbav:"startTime,omitempty" firestore:"startTime,omitempty"`
	EndTime   *time.Time `mapstructure:"endTime" json:"endTime,omitempty" gorm:"column:endtime" bson:"endTime,omitempty" dynamodbav:"endTime,omitempty" firestore:"endTime,omitempty"`
}
