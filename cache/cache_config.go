package cache

import "time"

// CacheConfig ...
type CacheConfig struct {
	Size             int64         `yaml:"size" mapstructure:"size" json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"` // byte
	CleaningEnable   bool          `yaml:"cleaning_enable" mapstructure:"cleaning_enable" json:"cleaningEnable,omitempty" gorm:"column:cleaningenable" bson:"cleaningEnable,omitempty" dynamodbav:"cleaningEnable,omitempty" firestore:"cleaningEnable,omitempty"`
	CleaningInterval time.Duration `yaml:"cleaning_interval" mapstructure:"cleaning_interval" json:"cleaningInterval,omitempty" gorm:"column:cleaninginterval" bson:"cleaningInterval,omitempty" dynamodbav:"cleaningInterval,omitempty" firestore:"cleaningInterval,omitempty"` // nano-second
}
