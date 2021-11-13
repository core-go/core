package search

type NumberRange struct {
	Min   *float64 `mapstructure:"min" json:"min,omitempty" gorm:"column:min" bson:"min,omitempty" dynamodbav:"min,omitempty" firestore:"min,omitempty"`
	Max   *float64 `mapstructure:"max" json:"max,omitempty" gorm:"column:max" bson:"max,omitempty" dynamodbav:"max,omitempty" firestore:"max,omitempty"`
	Lower *float64 `mapstructure:"lower" json:"lower,omitempty" gorm:"column:lower" bson:"lower,omitempty" dynamodbav:"lower,omitempty" firestore:"lower,omitempty"`
	Upper *float64 `mapstructure:"upper" json:"upper,omitempty" gorm:"column:upper" bson:"upper,omitempty" dynamodbav:"upper,omitempty" firestore:"upper,omitempty"`
}
