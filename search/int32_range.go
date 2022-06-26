package search

type Int32Range struct {
	Min     *int32 `yaml:"min" mapstructure:"min" json:"min,omitempty" gorm:"column:min" bson:"min,omitempty" dynamodbav:"min,omitempty" firestore:"min,omitempty"`
	Max     *int32 `yaml:"max" mapstructure:"max" json:"max,omitempty" gorm:"column:max" bson:"max,omitempty" dynamodbav:"max,omitempty" firestore:"max,omitempty"`
	Bottom  *int32 `yaml:"bottom" mapstructure:"bottom" json:"bottom,omitempty" gorm:"column:bottom" bson:"bottom,omitempty" dynamodbav:"bottom,omitempty" firestore:"bottom,omitempty"`
	Top     *int32 `yaml:"top" mapstructure:"top" json:"top,omitempty" gorm:"column:top" bson:"top,omitempty" dynamodbav:"top,omitempty" firestore:"top,omitempty"`
	Floor   *int32 `yaml:"floor" mapstructure:"floor" json:"floor,omitempty" gorm:"column:floor" bson:"floor,omitempty" dynamodbav:"floor,omitempty" firestore:"floor,omitempty"`
	Ceiling *int32 `yaml:"ceiling" mapstructure:"ceiling" json:"ceiling,omitempty" gorm:"column:ceiling" bson:"ceiling,omitempty" dynamodbav:"ceiling,omitempty" firestore:"ceiling,omitempty"`
	Lower   *int32 `yaml:"lower" mapstructure:"lower" json:"lower,omitempty" gorm:"column:lower" bson:"lower,omitempty" dynamodbav:"lower,omitempty" firestore:"lower,omitempty"`
	Upper   *int32 `yaml:"upper" mapstructure:"upper" json:"upper,omitempty" gorm:"column:upper" bson:"upper,omitempty" dynamodbav:"upper,omitempty" firestore:"upper,omitempty"`
}
