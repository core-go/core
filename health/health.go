package health

type Health struct {
	Status  string                 `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Data    map[string]interface{} `yaml:"data" mapstructure:"data" json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
	Details map[string]Health      `yaml:"details" mapstructure:"details" json:"details,omitempty" gorm:"column:details" bson:"details,omitempty" dynamodbav:"details,omitempty" firestore:"details,omitempty"`
}
