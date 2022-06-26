package search

import "time"

type TimeRange struct {
	Min       *time.Time `yaml:"min" mapstructure:"min" json:"min,omitempty" gorm:"column:startdate" bson:"min,omitempty" dynamodbav:"min,omitempty" firestore:"min,omitempty"`
	Max       *time.Time `yaml:"max" mapstructure:"max" json:"max,omitempty" gorm:"column:max" bson:"max,omitempty" dynamodbav:"max,omitempty" firestore:"max,omitempty"`
	Bottom    *time.Time `yaml:"bottom" mapstructure:"bottom" json:"bottom,omitempty" gorm:"column:bottom" bson:"bottom,omitempty" dynamodbav:"bottom,omitempty" firestore:"bottom,omitempty"`
	Top       *time.Time `yaml:"top" mapstructure:"top" json:"top,omitempty" gorm:"column:top" bson:"top,omitempty" dynamodbav:"top,omitempty" firestore:"top,omitempty"`
	Floor     *time.Time `yaml:"floor" mapstructure:"floor" json:"floor,omitempty" gorm:"column:floor" bson:"floor,omitempty" dynamodbav:"floor,omitempty" firestore:"floor,omitempty"`
	Ceiling   *time.Time `yaml:"ceiling" mapstructure:"ceiling" json:"ceiling,omitempty" gorm:"column:ceiling" bson:"ceiling,omitempty" dynamodbav:"ceiling,omitempty" firestore:"ceiling,omitempty"`
	Lower     *time.Time `yaml:"lower" mapstructure:"lower" json:"lower,omitempty" gorm:"column:lower" bson:"lower,omitempty" dynamodbav:"lower,omitempty" firestore:"lower,omitempty"`
	Upper     *time.Time `yaml:"upper" mapstructure:"upper" json:"upper,omitempty" gorm:"column:upper" bson:"upper,omitempty" dynamodbav:"upper,omitempty" firestore:"upper,omitempty"`
	StartTime *time.Time `yaml:"startTime" mapstructure:"startTime" json:"startTime,omitempty" gorm:"column:starttime" bson:"startTime,omitempty" dynamodbav:"startTime,omitempty" firestore:"startTime,omitempty"`
	EndTime   *time.Time `yaml:"endTime" mapstructure:"endTime" json:"endTime,omitempty" gorm:"column:endtime" bson:"endTime,omitempty" dynamodbav:"endTime,omitempty" firestore:"endTime,omitempty"`
}
