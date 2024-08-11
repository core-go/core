package search

import "time"

type DateRange struct {
	Min    *time.Time `yaml:"min" mapstructure:"min" json:"min,omitempty" gorm:"column:startdate" bson:"min,omitempty" dynamodbav:"min,omitempty" firestore:"min,omitempty"`
	Max    *time.Time `yaml:"max" mapstructure:"max" json:"max,omitempty" gorm:"column:max" bson:"max,omitempty" dynamodbav:"max,omitempty" firestore:"max,omitempty"`
	Bottom *time.Time `yaml:"bottom" mapstructure:"bottom" json:"bottom,omitempty" gorm:"column:bottom" bson:"bottom,omitempty" dynamodbav:"bottom,omitempty" firestore:"bottom,omitempty"`
	Top    *time.Time `yaml:"top" mapstructure:"top" json:"top,omitempty" gorm:"column:top" bson:"top,omitempty" dynamodbav:"top,omitempty" firestore:"top,omitempty"`
}
